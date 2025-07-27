package db

import (
	"course-backend/src/config"
	"course-backend/src/internal/models"
	models2 "course-backend/src/internal/models/course"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

var (
	// DB is the global database connection
	DB *gorm.DB
	// once is used to ensure that the database connection is initialized only once
	once sync.Once
)

// Connect initializes the database connection
func Connect() error {
	var err error
	once.Do(func() {
		dbConfig := config.GetDBConfig()
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", dbConfig.Host, dbConfig.Username, dbConfig.Password, dbConfig.Database, dbConfig.Port, dbConfig.SSLMode)
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	})
	if err != nil {
		return err
	}

	//Migrate all models
	err = Migrate(
		models.User{},
		models2.Course{},
		models2.Lesson{},
		models2.Chapter{},
		models2.Block{},
		models2.PaymentInfo{},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetDB returns the global database connection
func GetDB() *gorm.DB {
	return DB
}

// Close closes the database connection
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Migrate runs the database migrations
func Migrate(models ...interface{}) error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	for _, model := range models {
		if err := DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate model %T: %w", model, err)
		}
	}
	return nil
}

// Ping checks the database connection
func Ping() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %s", err.Error())
	}
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %s", err.Error())
	}
	return nil
}
