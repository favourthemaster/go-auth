package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaims struct {
	UserID  uuid.UUID `json:"userId"`
	Purpose string    `json:"purpose"` // <- Custom claim: "email_verification", "reset_password", etc.
	jwt.RegisteredClaims
}
