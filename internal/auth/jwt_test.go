package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func createTestToken(subject, secret string, expiry time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   subject,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func TestValidateJWT(t *testing.T) {
	validUUID := uuid.New()
	secret := "test-secret"

	tests := []struct {
		name        string
		setupToken  func() string
		secret      string
		expectedID  uuid.UUID
		expectError bool
	}{
		{
			name: "valid token returns correct UUID",
			setupToken: func() string {
				token, _ := createTestToken(validUUID.String(), secret, time.Hour)
				return token
			},
			secret:      secret,
			expectedID:  validUUID,
			expectError: false,
		},
		{
			name: "expired token returns error",
			setupToken: func() string {
				token, _ := createTestToken(validUUID.String(), secret, -time.Hour)
				return token
			},
			secret:      secret,
			expectedID:  uuid.Nil,
			expectError: true,
		},
		{
			name: "wrong secret returns error",
			setupToken: func() string {
				token, _ := createTestToken(validUUID.String(), secret, time.Hour)
				return token
			},
			secret:      "wrong-secret",
			expectedID:  uuid.Nil,
			expectError: true,
		},
		{
			name: "malformed token returns error",
			setupToken: func() string {
				return "this.is.notavalidtoken"
			},
			secret:      secret,
			expectedID:  uuid.Nil,
			expectError: true,
		},
		{
			name: "empty token returns error",
			setupToken: func() string {
				return ""
			},
			secret:      secret,
			expectedID:  uuid.Nil,
			expectError: true,
		},
		{
			name: "subject is not a valid UUID returns error",
			setupToken: func() string {
				token, _ := createTestToken("not-a-uuid", secret, time.Hour)
				return token
			},
			secret:      secret,
			expectedID:  uuid.Nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString := tt.setupToken()

			gotID, err := ValidateJWT(tokenString, tt.secret)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				if gotID != uuid.Nil {
					t.Errorf("expected uuid.Nil on error, got %v", gotID)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if gotID != tt.expectedID {
				t.Errorf("expected UUID %v, got %v", tt.expectedID, gotID)
			}
		})
	}
}