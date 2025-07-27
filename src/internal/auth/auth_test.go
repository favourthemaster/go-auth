package auth_test

import (
	"authentication/src/internal/auth"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "securePassword123"
	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}
	if hash == password {
		t.Error("Hashed password should not match original password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "securePassword123"
	hash, _ := auth.HashPassword(password)
	if !auth.CheckPasswordHash(password, hash) {
		t.Error("Password hash check failed")
	}
}
