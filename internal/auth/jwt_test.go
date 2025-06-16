package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	secret := "mysecret"
	userID := uuid.New()
	duration := time.Second * 2

	token, err := MakeJWT(userID, secret, duration)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	parsedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}

	if parsedID != userID {
		t.Errorf("Expected userID %v, got %v", userID, parsedID)
	}
}

func TestExpiredJWT(t *testing.T) {
	secret := "mysecret"
	userID := uuid.New()

	token, err := MakeJWT(userID, secret, -time.Second) // already expired
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Error("Expected error for expired token, got nil")
	}
}

func TestInvalidSecret(t *testing.T) {
	secret := "correctsecret"
	invalidSecret := "wrongsecret"
	userID := uuid.New()

	token, err := MakeJWT(userID, secret, time.Minute)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(token, invalidSecret)
	if err == nil {
		t.Error("Expected error for invalid secret, got nil")
	}
}
