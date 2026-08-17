package util

import (
	"testing"
)

func TestGenerateAndParseToken(t *testing.T) {
	secret := "test-secret-for-unit-test"
	token, err := GenerateToken(secret, 1, "user", "USER", 3600)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("token should not be empty")
	}
	claims, err := ParseToken(secret, token)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}
	if claims.UserID != 1 || claims.Username != "user" || claims.Role != "USER" {
		t.Errorf("claims mismatch: %+v", claims)
	}
}

func TestParseTokenInvalid(t *testing.T) {
	if _, err := ParseToken("secret", "invalid.token.here"); err == nil {
		t.Error("expected error for invalid token")
	}
}
