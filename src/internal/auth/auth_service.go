package auth

import (
	"authentication/src/internal/dto"
	"authentication/src/internal/errs"
	"authentication/src/internal/models"
	"authentication/src/internal/user"
	"authentication/src/utils"
	"context"
	"github.com/gofiber/fiber/v2/middleware/session"
	"time"
)

// AuthService defines authentication-related operations for users.
type AuthService interface {
	// Login authenticates a user with the provided credentials.
	Login(ctx context.Context, req *dto.LoginRequest, sess *session.Session) (*models.User, error)
	// Register creates a new user with the provided details.
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	// Logout logs out the user.
	Logout(ctx context.Context, req *dto.LogoutRequest, sess *session.Session) error
	// SendVerificationEmail sends a verification email to the user.
	SendVerificationEmail(ctx context.Context, req *dto.SendEmailVerificationRequest) error
	// VerifyEmail verifies the user's email using the provided token.
	VerifyEmail(ctx context.Context, req *dto.VerifyEmailRequest) error
	// ForgotPassword initiates the forgot password process for the user.
	ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) error
	// ResetPassword resets the user's password using the provided reset token.
	ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error

	// Additional methods can be added as needed

}

// authService implements AuthService for authentication logic.
type authService struct {
	UserService  user.UserService
	TokenService TokenService
	Mailer       utils.Mailer
}

// NewAuthService creates a new AuthService instance.
func NewAuthService(us user.UserService, ts TokenService) AuthService {
	return &authService{
		UserService:  us,
		TokenService: ts,
		Mailer:       utils.NewMailer(), // Assuming you have a Mailer implementation
	}
}

// Login authenticates a user with the provided credentials.
func (s *authService) Login(ctx context.Context, req *dto.LoginRequest, sess *session.Session) (*models.User, error) {

	getUserByEmailDTO := &dto.GetUserByEmailDTO{
		Email: req.Email,
	}

	loggedInUser, err := s.UserService.GetUserByEmail(ctx, getUserByEmailDTO)
	if err != nil {
		return nil, err
	}

	if loggedInUser == nil {
		return nil, errs.ErrUserNotFound
	}

	if isPasswordValid := utils.ComparePassword(req.Password, loggedInUser.PasswordHash); !isPasswordValid {
		return nil, errs.ErrInvalidCredentials
	}

	if !loggedInUser.Verified {
		return nil, errs.ErrEmailNotVerified
	}

	sess.Set("userID", loggedInUser.ID)
	err = sess.Save()
	if err != nil {
		return nil, err
	}

	return loggedInUser, nil
}

// Register creates a new user with the provided details
func (s *authService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {

	createUserDTO := &dto.CreateUserDTO{
		Email:    req.Email,
		FullName: req.FullName,
		Password: req.Password,
	}

	createdUser, err := s.UserService.CreateUser(ctx, createUserDTO)

	if err != nil {
		return nil, err
	}

	res := &dto.RegisterResponse{
		UserID: createdUser.ID,
	}

	return res, nil
}

// Logout logs out the user
func (s *authService) Logout(ctx context.Context, req *dto.LogoutRequest, sess *session.Session) error {
	// Destroy the session
	return sess.Destroy()
}

// SendVerificationEmail sends a verification email to the user
func (s *authService) SendVerificationEmail(ctx context.Context, req *dto.SendEmailVerificationRequest) error {
	purpose := "email_verification"
	expiry := time.Duration(time.Minute * 30)
	token, err := s.TokenService.GenerateToken(ctx, req.ID, purpose, expiry)
	if err != nil {
		return err
	}

	err = s.Mailer.SendVerificationMail(req.Email, token)
	if err != nil {
		return err
	}
	return nil
}

// VerifyEmail verifies the user's email using the provided token
func (s *authService) VerifyEmail(ctx context.Context, req *dto.VerifyEmailRequest) error {

	expectedPurpose := "email_verification"

	claims, err := s.TokenService.ValidateToken(ctx, req.Token, expectedPurpose)

	if err != nil {
		return err
	}

	unverifiedUser, err := s.UserService.GetUserByID(ctx, claims.UserID)

	if err != nil {
		return err
	}

	unverifiedUser.Verified = true

	_, err = s.UserService.UpdateUser(ctx, unverifiedUser)

	if err != nil {
		return err
	}

	return nil
}

// ForgotPassword initiates the forgot password process for the user
func (s *authService) ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) error {

	getUserByEmailDTO := &dto.GetUserByEmailDTO{
		Email: req.Email,
	}

	existingUser, err := s.UserService.GetUserByEmail(ctx, getUserByEmailDTO)
	if err != nil {
		return err
	}

	if existingUser == nil {
		return errs.ErrUserNotFound
	}

	purpose := "password_reset"

	expiry := time.Duration(time.Minute * 30)

	token, err := s.TokenService.GenerateToken(ctx, existingUser.ID, purpose, expiry)

	if err != nil {
		return err
	}

	err = s.Mailer.SendPasswordResetMail(existingUser.Email, token)
	if err != nil {
		return err
	}

	return nil
}

// ResetPassword resets the user's password using the provided reset token
func (s *authService) ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error {
	expectedPurpose := "password_reset"

	claims, err := s.TokenService.ValidateToken(ctx, req.ResetToken, expectedPurpose)

	if err != nil {
		return err
	}

	existingUser, err := s.UserService.GetUserByID(ctx, claims.UserID)

	if err != nil {
		return err
	}

	if existingUser == nil {
		return errs.ErrUserNotFound
	}

	existingUser.PasswordHash, err = utils.HashPassword(req.NewPassword)
	if err != nil {
		return errs.ErrInternalServerError // Error hashing password
	}

	_, err = s.UserService.UpdateUser(ctx, existingUser)
	if err != nil {
		return err
	}

	return nil
}
