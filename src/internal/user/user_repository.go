package user

import (
	"context"
	"course-backend/src/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID, permanent bool) error

	ListUsers(ctx context.Context, limit, offset int) ([]models.User, error)
	GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// CreateUser creates a new user in the database
func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// UpdateUser updates an existing user in the database
func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// GetUserByID creates a new user in the database
func (r *userRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// DeleteUser deletes a user from the database
func (r *userRepository) DeleteUser(ctx context.Context, userID uuid.UUID, permanent bool) error {
	if permanent {
		// Permanent delete
		return r.db.WithContext(ctx).Unscoped().Delete(&models.User{}, "id = ?", userID).Error
	} else {
		// Soft delete
		return r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", userID).Error
	}
}

// ListUsers retrieves all users from the database with pagination
func (r *userRepository) ListUsers(ctx context.Context, limit, offset int) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

// GetUsersByIDs retrieves users by a slice of IDs
func (r *userRepository) GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).Where("id IN ?", userIDs).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
