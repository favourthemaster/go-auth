package user

import (
	"context"
	"course-backend/src/internal/dto"
	"course-backend/src/internal/errs"
	"course-backend/src/internal/models"
	"course-backend/src/utils"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	// GetUserByEmail retrieves a user by email
	GetUserByEmail(ctx context.Context, emailDTO *dto.GetUserByEmailDTO) (*models.User, error)
	// CreateUser creates a new user
	CreateUser(ctx context.Context, userDTO *dto.CreateUserDTO) (*models.User, error)
	// UpdateUser updates an existing user
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	// DeleteUser deletes a user by ID
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	// ListUsers lists users with pagination
	//ListUsers(ctx context.Context, limit, offset int) ([]dto.UserResponse, error)
}

type userService struct {
	ur UserRepository
}

// NewUserService creates a new UserService instance
func NewUserService(ur UserRepository) UserService {
	return &userService{
		ur: ur,
	}
}

func (u userService) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := u.ur.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound // User not found
		}
		return nil, err // Other error
	}
	return user, nil
}

func (u userService) GetUserByEmail(ctx context.Context, emailDTO *dto.GetUserByEmailDTO) (*models.User, error) {
	user, err := u.ur.GetUserByEmail(ctx, emailDTO.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // User not found
		}
		return nil, err // Other error
	}

	return user, nil
}

func (u userService) CreateUser(ctx context.Context, userDTO *dto.CreateUserDTO) (*models.User, error) {

	existingUser, err := u.ur.GetUserByEmail(ctx, userDTO.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User not found, proceed with creation
		} else {
			// Some other error occurred
			return nil, err
		}
	}

	if existingUser != nil {
		if !existingUser.Verified {
			err := u.ur.DeleteUser(ctx, existingUser.ID, true)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errs.ErrUserAlreadyExists
		}

	}

	hashedPassword, err := utils.HashPassword(userDTO.Password)
	if err != nil {
		return nil, errs.ErrInternalServerError // Error hashing password
	}

	newUser := &models.User{
		Email:        userDTO.Email,
		PasswordHash: hashedPassword,
		FullName:     userDTO.FullName,
	}

	err = u.ur.CreateUser(ctx, newUser)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (u userService) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {

	_, err := u.ur.GetUserByID(ctx, user.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound // User not found
		}
		return nil, err // Other error
	}

	err = u.ur.UpdateUser(ctx, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u userService) DeleteUser(ctx context.Context, userID uuid.UUID) error {

	_, err := u.ur.GetUserByID(ctx, userID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.ErrUserNotFound // User not found
		}
	}

	err = u.ur.DeleteUser(ctx, userID, false) // Soft delete

	if err != nil {
		return err // Other error
	}

	return nil
}

func (u userService) ListUsers(ctx context.Context, limit, offset int) ([]dto.UserResponse, error) {

	users, err := u.ur.ListUsers(ctx, limit, offset)

	if err != nil {
		return nil, err // Error retrieving users
	}

	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponse{
			ID:       user.ID,
			FullName: user.FullName,
			Email:    user.Email,
		})
	}
	return userResponses, nil
}
