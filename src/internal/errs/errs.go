package errs

import "errors"

var (
	//User errors
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailNotVerified  = errors.New("email not verified")

	ErrInvalidCredentials = errors.New("invalid credentials")

	// Token errors
	ErrTokenNotFound       = errors.New("token not found")
	ErrInvalidTokenPurpose = errors.New("invalid token purpose")
	ErrInvalidToken        = errors.New("invalid token")
	ErrTokenExpired        = errors.New("token expired")

	ErrRedisTokenDeletion = errors.New("error deleting token from Redis")

	ErrInvalidBlockData = errors.New("invalid block data")

	ErrInvalidRequestBody  = errors.New("invalid request body")
	ErrInternalServerError = errors.New("internal server error")
)
