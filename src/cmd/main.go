package main

import (
	"course-backend/src/internal/auth"
	"course-backend/src/internal/course"
	"course-backend/src/internal/course/block"
	"course-backend/src/internal/course/chapter"
	"course-backend/src/internal/course/lesson"
	"course-backend/src/internal/db"
	"course-backend/src/internal/user"
	"encoding/gob"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
	"os"
)

func main() {

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

	// Course Module Wiring
	courseRepo := course.NewCourseRepository(database)
	courseService := course.NewCourseService(courseRepo)

	// Block Module Wiring
	blockRepo := block.NewBlockRepository(database)
	blockService := block.NewBlockService(blockRepo)

	// Lesson Module Wiring
	lessonRepo := lesson.NewLessonRepository(database)
	lessonService := lesson.NewLessonService(lessonRepo, blockService)

	// Chapter Module Wiring
	chapterRepo := chapter.NewChapterRepository(database)
	chapterService := chapter.NewChapterService(chapterRepo, lessonService)

	courseHandler := course.NewCourseHandler(courseService, chapterService, lessonService, blockService)
	courseGroup := app.Group("/courses")
	//Get Courses with pagination
	courseGroup.Get("/all", courseHandler.GetAllCourses)
	courseGroup.Get("/course/:id", courseHandler.GetCourseByID)
	//Get Full Course by ID
	courseGroup.Get("/full/:id", courseHandler.GetFullCourseByID)
	courseGroup.Post("/", auth.RequireAuth(), courseHandler.CreateCourse)
	courseGroup.Put("/", auth.RequireAuth(), courseHandler.UpdateCourse)
	courseGroup.Delete("/", auth.RequireAuth(), courseHandler.DeleteCourse)
	courseGroup.Post("/chapters", auth.RequireAuth(), courseHandler.CreateChapters)
	courseGroup.Put("/chapters", auth.RequireAuth(), courseHandler.UpdateChapters)
	courseGroup.Delete("/chapters/:id", auth.RequireAuth(), courseHandler.DeleteChapter)
	courseGroup.Post("/lessons", auth.RequireAuth(), courseHandler.CreateLessons)
	courseGroup.Put("/lessons", auth.RequireAuth(), courseHandler.UpdateLessons)
	courseGroup.Delete("/lessons/:id", auth.RequireAuth(), courseHandler.DeleteLesson)
	courseGroup.Post("/blocks", auth.RequireAuth(), courseHandler.CreateBlocks)
	courseGroup.Put("/blocks", auth.RequireAuth(), courseHandler.UpdateBlocks)
	courseGroup.Delete("/blocks/:id", auth.RequireAuth(), courseHandler.DeleteBlock)

	err = app.Listen(":3000")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
