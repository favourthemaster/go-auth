package auth

import (
	"course-backend/src/internal/dto"
	"course-backend/src/internal/errs"
	"course-backend/src/utils"
	"errors"
	"github.com/gofiber/fiber/v2"
	"log"
)

type AuthHandler struct {
	AuthService
}

// NewAuthHandler creates a new authHandler with the provided AuthService.
func NewAuthHandler(as AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: as,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {

	ctx := c.Context()
	var req dto.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
			err, "Failed to parse request body"))
	}

	// Validate the request
	if validationErr := utils.ValidateStruct(req); validationErr != nil {
		log.Printf("Validation error: %v", validationErr)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
			validationErr, "Validation failed"))
	}

	res, err := h.AuthService.Register(ctx, &req)
	if err != nil {
		log.Printf("Error during registration: %v", err)

		if errors.Is(err, errs.ErrUserAlreadyExists) {
			return c.Status(fiber.StatusConflict).JSON(utils.ErrorResponse(
				err, "A user with this email address already exists"))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			err, "Registration failed"))
	}

	emailReq := &dto.SendEmailVerificationRequest{
		ID:    res.UserID,
		Email: req.Email,
	}
	err = h.AuthService.SendVerificationEmail(ctx, emailReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			err, "Failed to send verification email"))
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(nil, "User registered successfully, please check your email for verification instructions"))
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	ctx := c.Context()
	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
			err, "Failed to parse request body"))
	}

	// Validate the request
	if validationErr := utils.ValidateStruct(req); validationErr != nil {
		log.Printf("Validation error: %v", validationErr)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
			validationErr, "Validation failed"))
	}

	sess, err := store.Get(c)
	if err != nil {
		log.Printf("Error retrieving session: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			err, "Failed to retrieve session"))
	}

	loggedInUser, err := h.AuthService.Login(ctx, &req, sess)
	if err != nil {
		log.Printf("Error during login: %v", err)

		if errors.Is(err, errs.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(
				err, "This user does not exist"))
		}

		if errors.Is(err, errs.ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse(
				err, "Invalid email or password"))
		}

		if errors.Is(err, errs.ErrEmailNotVerified) {
			return c.Status(fiber.StatusForbidden).JSON(utils.ErrorResponse(
				err, "Email not verified, please check your inbox for the verification email or sign up again"))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			err, "Login failed"))
	}

	loginResponse := dto.LoginResponse{
		User: dto.UserResponse{
			ID:       loggedInUser.ID,
			FullName: loggedInUser.FullName,
			Email:    loggedInUser.Email,
		},
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(loginResponse, "Login successful"))
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	ctx := c.Context()
	var req dto.LogoutRequest

	sess, err := store.Get(c)
	if err != nil {
		log.Printf("Error retrieving session: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			err, "Failed to retrieve session"))
	}

	err = h.AuthService.Logout(ctx, &req, sess)
	if err != nil {
		log.Printf("Error deleting session: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(err, "The Session could not be removed"))
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(nil, "Logout successful"))
}

// VerifyEmail verifies the user's email address
func (h *AuthHandler) VerifyEmail(c *fiber.Ctx) error {
	ctx := c.Context()
	var req dto.VerifyEmailRequest

	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
			err, "Failed to parse request body"))
	}

	// Validate the request
	if validationErr := utils.ValidateStruct(req); validationErr != nil {
		log.Printf("Validation error: %v", validationErr)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
			validationErr, "Validation failed"))
	}

	err := h.AuthService.VerifyEmail(ctx, &req)
	if err != nil {
		log.Printf("Error during email verification: %v", err)

		if errors.Is(err, errs.ErrInvalidToken) {
			return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
				err, "Invalid verification link"))
		}

		if errors.Is(err, errs.ErrTokenExpired) {
			return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
				err, "Verification link has expired, please request a new one"))
		}

		if errors.Is(err, errs.ErrInvalidTokenPurpose) {
			return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
				err, "Invalid token purpose, please request a new verification email"))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			err, "Email verification failed"))
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(nil, "Email verified successfully, you can now log in"))
}

// SendVerificationEmail sends a verification email to the user
func (h *AuthHandler) SendVerificationEmail(c *fiber.Ctx) error {
	ctx := c.Context()
	var req dto.SendEmailVerificationRequest

	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
			err, "Failed to parse request body"))
	}

	// Validate the request
	if validationErr := utils.ValidateStruct(req); validationErr != nil {
		log.Printf("Validation error: %v", validationErr)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
			validationErr, "Validation failed"))
	}

	err := h.AuthService.SendVerificationEmail(ctx, &req)
	if err != nil {
		log.Printf("Error sending verification email: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			err, "Failed to send verification email"))
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(nil, "Verification email sent successfully"))
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	ctx := c.Context()
	var req dto.ForgotPasswordRequest

	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
			err, "Failed to parse request body"))
	}

	// Validate the request
	if validationErr := utils.ValidateStruct(req); validationErr != nil {
		log.Printf("Validation error: %v", validationErr)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
			validationErr, "Validation failed"))
	}

	err := h.AuthService.ForgotPassword(ctx, &req)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			log.Printf("User not found: %v", err)
			return c.Status(fiber.StatusNotFound).JSON(utils.ErrorResponse(
				err, "This email address is not registered"))
		}
		log.Printf("Error during forgot password: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			err, "Forgot password failed"))
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(nil, "A password reset link has been sent to your email"))
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	ctx := c.Context()
	var req dto.ResetPasswordRequest

	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
			err, "Failed to parse request body"))
	}

	// Validate the request
	if validationErr := utils.ValidateStruct(req); validationErr != nil {
		log.Printf("Validation error: %v", validationErr)
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
			validationErr, "Validation failed"))
	}

	//Get and destroy any previously existing sessions
	sess, err := store.Get(c)
	if err != nil {
		log.Printf("Error retrieving session: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			err, "Failed to retrieve session"))
	}
	err = sess.Destroy()
	if err != nil {
		log.Printf("Error destroying session: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			err, "Failed to destroy session"))
	}

	err = h.AuthService.ResetPassword(ctx, &req)
	if err != nil {
		log.Printf("Error during password reset: %v", err)

		if errors.Is(err, errs.ErrInvalidToken) {
			return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
				err, "Invalid reset link"))
		}

		if errors.Is(err, errs.ErrTokenExpired) {
			return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(
				err, "Reset link has expired, please request a new one"))
		}

		return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse(
			err, "Password reset failed"))
	}

	return c.Status(fiber.StatusOK).JSON(utils.SuccessResponse(nil, "Password reset successfully, you can now log in with your new password"))
}
