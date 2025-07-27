package dto

import (
	"course-backend/src/internal/models"
	"github.com/google/uuid"
)

// UserResponse represents the response structure for user-related operations.
type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
}

// ToUserResponse converts a models.User to a UserResponse DTO.
func ToUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
	}
}

// ToUserResponseList converts a slice of models.User to a slice of UserResponse DTOs.
func ToUserResponseList(users []*models.User) []UserResponse {
	res := make([]UserResponse, len(users))
	for i, u := range users {
		res[i] = ToUserResponse(u)
	}
	return res
}

type CreateUserDTO struct {
	FullName string `json:"full_name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type GetUserByEmailDTO struct {
	Email string `json:"email" validate:"required,email"`
}
