# Chirpy API

Chirpy is a microblogging API that allows users to create, view, and manage "chirps" (short messages). It includes user authentication, refresh token management, and administrative features.

## Table of Contents

  - [Chirpy API](https://www.google.com/search?q=%23chirpy-api)
      - [Table of Contents](https://www.google.com/search?q=%23table-of-contents)
      - [Features](https://www.google.com/search?q=%23features)
      - [Getting Started](https://www.google.com/search?q=%23getting-started)
          - [Prerequisites](https://www.google.com/search?q=%23prerequisites)
          - [Installation](https://www.google.com/search?q=%23installation)
          - [Configuration](https://www.google.com/search?q=%23configuration)
          - [Running the API](https://www.google.com/search?q=%23running-the-api)
      - [API Endpoints](https://www.google.com/search?q=%23api-endpoints)
          - [Health Check](https://www.google.com/search?q=%23health-check)
          - [Metrics](https://www.google.com/search?q=%23metrics)
          - [Admin (Development Only)](https://www.google.com/search?q=%23admin-development-only)
          - [Users](https://www.google.com/search?q=%23users)
          - [Authentication](https://www.google.com/search?q=%23authentication)
          - [Chirps](https://www.google.com/search?q=%23chirps)
          - [Webhooks](https://www.google.com/search?q=%23webhooks)
      - [Authentication](https://www.google.com/search?q=%23authentication-1)
      - [Error Handling](https://www.google.com/search?q=%23error-handling)
      - [Development](https://www.google.com/search?q=%23development)
      - [Contributing](https://www.google.com/search?q=%23contributing)
      - [License](https://www.google.com/search?q=%23license)

## Features

  * **User Management:** Register, update, and retrieve user information.
  * **Authentication:** Secure user login with JWT access tokens and refresh tokens.
  * **Chirp Management:** Create, view, and delete chirps.
  * **Profanity Filtering:** Automatic filtering of certain words in chirps.
  * **Admin Tools:** Metrics tracking and a development-only reset endpoint.
  * **Webhook Support:** Endpoint for external services (e.g., Polka) to trigger user upgrades.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

  * Go (version 1.22 or higher recommended)
  * PostgreSQL database
  * `make` (optional, for convenience scripts)

### Installation

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/Rishan-Jadva/chirpy.git
    cd chirpy
    ```

2.  **Install Go dependencies:**

    ```bash
    go mod tidy
    ```

### Configuration

The API relies on environment variables for sensitive information and configuration. Create a `.env` file in the root directory of the project with the following variables:

```dotenv
DB_URL="postgresql://user:password@host:port/database_name?sslmode=disable"
JWT_SECRET="your_jwt_secret_key"
POLKA_KEY="your_polka_webhook_secret_key"
PLATFORM="dev" # Set to "dev" for development features like /admin/reset, otherwise leave it empty or set to "prod"
```

  * **`DB_URL`**: Your PostgreSQL connection string.
  * **`JWT_SECRET`**: A strong, secret key used for signing and verifying JWTs.
  * **`POLKA_KEY`**: A secret key used to authenticate requests from the Polka webhook.
  * **`PLATFORM`**: Set to `"dev"` to enable the `/admin/reset` endpoint, which is useful for development but should *not* be enabled in production.

### Running the API

1.  **Ensure your `.env` file is configured.**

2.  **Run the application:**

    ```bash
    go run .
    ```

    The server will start on `http://localhost:8080`.

## API Endpoints

The API serves content on `http://localhost:8080`.

### Health Check

  * **GET `/api/healthz`**
      * **Description:** Checks the health of the API.
      * **Response:** `200 OK` with a plain text body of `OK`.

### Metrics

  * **GET `/admin/metrics`**
      * **Description:** Displays the number of times the file server has been hit.
      * **Response:** HTML page showing the hit count.

### Admin (Development Only)

  * **POST `/admin/reset`**
      * **Description:** Resets the file server hit counter and deletes all users from the database. **Only available when `PLATFORM` is set to `"dev"` in the `.env` file.**
      * **Response:** `200 OK` with a plain text body `Hits counter reset`.

### Users

  * **POST `/api/users`**
      * **Description:** Registers a new user.
      * **Request Body (JSON):**
        ```json
        {
          "email": "user@example.com",
          "password": "mySecurePassword"
        }
        ```
      * **Response (JSON):** `201 Created`
        ```json
        {
          "id": "uuid",
          "created_at": "timestamp",
          "updated_at": "timestamp",
          "email": "user@example.com",
          "is_chirpy_red": false
        }
        ```
  * **PUT `/api/users`**
      * **Description:** Updates an existing user's email and password. Requires a valid JWT in the `Authorization` header.
      * **Request Headers:**
          * `Authorization: Bearer <JWT_ACCESS_TOKEN>`
      * **Request Body (JSON):**
        ```json
        {
          "email": "new_email@example.com",
          "password": "newSecurePassword"
        }
        ```
      * **Response (JSON):** `200 OK`
        ```json
        {
          "id": "uuid",
          "created_at": "timestamp",
          "updated_at": "timestamp",
          "email": "new_email@example.com",
          "is_chirpy_red": false
        }
        ```

### Authentication

  * **POST `/api/login`**
      * **Description:** Authenticates a user and returns JWT access and refresh tokens.
      * **Request Body (JSON):**
        ```json
        {
          "email": "user@example.com",
          "password": "mySecurePassword"
        }
        ```
      * **Response (JSON):** `200 OK`
        ```json
        {
          "id": "uuid",
          "created_at": "timestamp",
          "updated_at": "timestamp",
          "email": "user@example.com",
          "is_chirpy_red": false,
          "token": "jwt_access_token",
          "refresh_token": "jwt_refresh_token"
        }
        ```
  * **POST `/api/refresh`**
      * **Description:** Refreshes an expired access token using a valid refresh token.
      * **Request Headers:**
          * `Authorization: Bearer <JWT_REFRESH_TOKEN>`
      * **Response (JSON):** `200 OK`
        ```json
        {
          "token": "new_jwt_access_token"
        }
        ```
  * **POST `/api/revoke`**
      * **Description:** Revokes a refresh token, invalidating it. The corresponding access token will no longer be refreshable.
      * **Request Headers:**
          * `Authorization: Bearer <JWT_REFRESH_TOKEN>`
      * **Response:** `204 No Content`

### Chirps

  * **POST `/api/chirps`**
      * **Description:** Creates a new chirp. Requires a valid JWT in the `Authorization` header. Chirp body has a maximum length of 140 characters and undergoes profanity filtering.
      * **Request Headers:**
          * `Authorization: Bearer <JWT_ACCESS_TOKEN>`
      * **Request Body (JSON):**
        ```json
        {
          "body": "This is my first chirp!"
        }
        ```
      * **Response (JSON):** `201 Created`
        ```json
        {
          "id": "uuid",
          "created_at": "timestamp",
          "updated_at": "timestamp",
          "body": "This is my first ****!",
          "user_id": "user_uuid"
        }
        ```
  * **GET `/api/chirps`**
      * **Description:** Retrieves a list of chirps. Can be filtered by `author_id` and sorted by `created_at`.
      * **Query Parameters:**
          * `author_id` (optional): Filter chirps by the specified user ID.
          * `sort` (optional): Sort order for chirps.
              * `asc` (default): Sort by `created_at` in ascending order.
              * `desc`: Sort by `created_at` in descending order.
      * **Response (JSON):** `200 OK`
        ```json
        [
          {
            "id": "uuid",
            "created_at": "timestamp",
            "updated_at": "timestamp",
            "body": "Hello Chirpy!",
            "user_id": "user_uuid"
          }
        ]
        ```
  * **GET `/api/chirps/{chirpID}`**
      * **Description:** Retrieves a single chirp by its ID.
      * **Path Parameters:**
          * `chirpID` (UUID): The ID of the chirp to retrieve.
      * **Response (JSON):** `200 OK`
        ```json
        {
          "id": "uuid",
          "created_at": "timestamp",
          "updated_at": "timestamp",
          "body": "A specific chirp.",
          "user_id": "user_uuid"
        }
        ```
  * **DELETE `/api/chirps/{chirpID}`**
      * **Description:** Deletes a chirp by its ID. Only the author of the chirp can delete it. Requires a valid JWT in the `Authorization` header.
      * **Request Headers:**
          * `Authorization: Bearer <JWT_ACCESS_TOKEN>`
      * **Path Parameters:**
          * `chirpID` (UUID): The ID of the chirp to delete.
      * **Response:** `204 No Content`

### Webhooks

  * **POST `/api/polka/webhooks`**
      * **Description:** Endpoint for Polka webhooks to signal a user upgrade. Requires a valid `X-Api-Key` header matching `POLKA_KEY`. If the event is `user.upgraded`, it attempts to upgrade the specified user to Chirpy Red.
      * **Request Headers:**
          * `X-Api-Key: <POLKA_KEY>`
      * **Request Body (JSON):**
        ```json
        {
          "event": "user.upgraded",
          "data": {
            "user_id": "uuid_of_user_to_upgrade"
          }
        }
        ```
      * **Response:** `204 No Content` (if successful or event is not `user.upgraded`)
      * **Errors:** `401 Unauthorized` if `X-Api-Key` is missing or incorrect. `404 Not Found` if the user to upgrade is not found.

## Authentication

This API uses JSON Web Tokens (JWTs) for authentication.

  * **Access Tokens:** Short-lived tokens used to authenticate requests to protected endpoints. These should be sent in the `Authorization` header as `Bearer <token>`.
  * **Refresh Tokens:** Long-lived tokens used to obtain new access tokens when the current one expires. These are also sent in the `Authorization` header as `Bearer <token>` to the `/api/refresh` endpoint.

## Error Handling

The API returns JSON error responses with a `message` and an optional `error` field for detailed debugging information (though the `error` field might be omitted in production environments for security reasons).

Example error response:

```json
{
  "error": "Couldn't decode parameters",
  "message": "Internal Server Error"
}
```

Common HTTP status codes used for errors:

  * `400 Bad Request`: Invalid request payload or parameters.
  * `401 Unauthorized`: Missing or invalid authentication credentials (e.g., JWT).
  * `403 Forbidden`: Authenticated, but not authorized to perform the action.
  * `404 Not Found`: Resource not found.
  * `500 Internal Server Error`: An unexpected server-side error occurred.

## Development

To aid in development, the `PLATFORM` environment variable can be set to `"dev"`. This enables the `/admin/reset` endpoint, which is useful for clearing the database during testing.

## License

This project is licensed under the MIT License - see the `LICENSE` file for details