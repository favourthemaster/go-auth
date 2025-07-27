package auth

import (
	"authentication/src/internal/db"
	"authentication/src/internal/errs"
	"authentication/src/internal/models"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

type TokenService interface {
	GenerateToken(ctx context.Context, userID uuid.UUID, purpose string, expiry time.Duration) (string, error)
	ValidateToken(ctx context.Context, token, expectedPurpose string) (*models.CustomClaims, error)
}

type tokenService struct {
	secretKey  string
	redisStore *redis.Client
}

func NewTokenService(secretKey string) TokenService {
	return &tokenService{
		secretKey:  secretKey,
		redisStore: db.GetRedisClient(),
	}
}

func (t *tokenService) GenerateToken(ctx context.Context, userID uuid.UUID, purpose string, expiry time.Duration) (string, error) {
	claims := models.CustomClaims{
		UserID:  userID,
		Purpose: purpose,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   purpose,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(t.secretKey))
	if err != nil {
		return "", err
	}
	//Delete any existing token for the user and purpose
	existingToken := t.redisStore.Get(ctx, fmt.Sprintf("%s:%s", purpose, userID))
	if existingToken != nil {
		err = t.redisStore.Del(ctx, fmt.Sprintf("%s:%s", purpose, userID)).Err()
		if err != nil && !errors.Is(err, redis.Nil) {
			return "", err // Handle error if deletion fails
		}
	}
	// Store the token in Redis with an expiry
	t.redisStore.Set(ctx, fmt.Sprintf("%s:%s", purpose, userID), signedToken, expiry)

	return signedToken, nil
}

func (t *tokenService) ValidateToken(ctx context.Context, token, expectedPurpose string) (*models.CustomClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errs.ErrTokenExpired
		}

		return nil, err
	}

	claims, ok := parsedToken.Claims.(*models.CustomClaims)
	if !ok {
		return nil, err
	}

	expectedToken := t.redisStore.Get(ctx, fmt.Sprintf("%s:%s", expectedPurpose, claims.UserID))

	if expectedToken == nil {
		return nil, errs.ErrTokenNotFound
	}

	if expectedToken.Val() != token {
		return nil, errs.ErrInvalidToken
	}

	if claims.Purpose != expectedPurpose {
		return nil, errs.ErrInvalidTokenPurpose
	}

	// Delete the token from Redis after validation
	err = t.redisStore.Del(ctx, fmt.Sprintf("%s:%s", expectedPurpose, claims.UserID)).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, errs.ErrRedisTokenDeletion
	}

	return claims, nil
}
