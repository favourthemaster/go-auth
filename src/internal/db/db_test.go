package db_test

import (
	"authentication/src/internal/db"
	"testing"
)

func TestConnect(t *testing.T) {
	err := db.Connect()
	if err != nil {
		t.Fatalf("Database connection failed: %v", err)
	}
}
