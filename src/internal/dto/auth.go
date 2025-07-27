package dto

import "github.com/google/uuid"

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the response body for user login
type LoginResponse struct {
	User UserResponse `json:"user"`
}

//----------------------------Register------------------------------------

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	FullName string `json:"name" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required,min=8,max=64"`
	Email    string `json:"email" validate:"required,email"`
}

// RegisterResponse represents the response body for user registration
type RegisterResponse struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

// ----------------------------Logout------------------------------------

// LogoutRequest represents the request body for user logout
type LogoutRequest struct {
}

// LogoutResponse represents the response body for user logout
type LogoutResponse struct {
}

// -----------------------------Forgot-Password------------------------------
// ForgotPasswordRequest represents the request body for forgot password
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ForgotPasswordResponse represents the response body for forgot password
type ForgotPasswordResponse struct {
}

// -----------------------------Reset-Password------------------------------
// ResetPasswordRequest represents the request body for reset password
type ResetPasswordRequest struct {
	UserID      uuid.UUID `json:"user_id" validate:"required"`
	ResetToken  string    `json:"reset_token" validate:"required"`
	NewPassword string    `json:"new_password" validate:"required,min=8,max=64"`
}

// ResetPasswordResponse represents the response body for reset password
type ResetPasswordResponse struct {
}

// -----------------------------Email-Verification-----------------------------

// SendEmailVerificationRequest represents the request body for email verification
type SendEmailVerificationRequest struct {
	ID    uuid.UUID `json:"id" validate:"required"`
	Email string    `json:"email" validate:"required,email"`
}

// VerifyEmailRequest represents the request body for verifying email
type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}
