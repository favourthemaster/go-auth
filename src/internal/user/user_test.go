package user_test

import (
	"authentication/src/internal/user"
	"authentication/src/models"
	"testing"
)

func TestCreateUser(t *testing.T) {
	service := user.NewService(nil) // Pass nil or mock repository
	newUser := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	created, err := service.CreateUser(newUser)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	if created.Username != newUser.Username {
		t.Error("Usernames do not match")
	}
}
