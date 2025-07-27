package auth

import (
	"github.com/gofiber/fiber/v2/middleware/session"
	redisstore "github.com/gofiber/storage/redis"
	"time"
)

var store *session.Store

func InitSessionStore() {
	redis := redisstore.New(redisstore.Config{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		Database: 0,
		Reset:    false,
	})

	store = session.New(session.Config{
		Storage:        redis,
		Expiration:     24 * time.Hour,
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: "Lax",
	})
}
