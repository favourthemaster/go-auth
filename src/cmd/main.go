package main

import (
	"authentication/src/config"
	"authentication/src/internal/auth"
	"authentication/src/internal/db"
	"authentication/src/internal/user"
	"encoding/gob"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
	"os"
)

func main() {
	// Initialize config and load environment variables
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	app := fiber.New()

	err := db.Connect()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		os.Exit(1)
	}

	auth.InitSessionStore()

	database := db.GetDB()
	if database == nil {
		log.Fatal("Failed to initialize database")
	}

	db.InitRedisFromConfig()

	gob.Register(uuid.UUID{})
	userRepo := user.NewUserRepository(database)
	userService := user.NewUserService(userRepo)

	tokenService := auth.NewTokenService("your-secret-key") // Replace with your actual secret key
	if tokenService == nil {
		log.Fatal("Failed to initialize token service")
	}
	authService := auth.NewAuthService(userService, tokenService)
	if authService == nil {
		log.Fatal("Failed to initialize auth service")
	}
	authHandler := auth.NewAuthHandler(authService)
	authGroup := app.Group("/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
	authGroup.Post("/logout", auth.RequireAuth(), authHandler.Logout)
	authGroup.Post("/forgot-password", authHandler.ForgotPassword)
	authGroup.Post("/resend-verification-email", authHandler.SendVerificationEmail)
	authGroup.Post("/verify-email/", authHandler.VerifyEmail)
	authGroup.Post("/reset-password", authHandler.ResetPassword)

	err = app.Listen(":3000")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
