package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		header    http.Header
		expectErr bool
		expected  string
	}{
		{
			name:      "valid header",
			header:    http.Header{"Authorization": []string{"Bearer abc.def.ghi"}},
			expectErr: false,
			expected:  "abc.def.ghi",
		},
		{
			name:      "missing header",
			header:    http.Header{},
			expectErr: true,
		},
		{
			name:      "invalid format",
			header:    http.Header{"Authorization": []string{"Token abc.def.ghi"}},
			expectErr: true,
		},
		{
			name:      "empty token",
			header:    http.Header{"Authorization": []string{"Bearer "}},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.header)
			if (err != nil) != tt.expectErr {
				t.Errorf("unexpected error result: got %v, wantErr %v", err, tt.expectErr)
			}
			if token != tt.expected {
				t.Errorf("unexpected token: got %v, want %v", token, tt.expected)
			}
		})
	}
}
