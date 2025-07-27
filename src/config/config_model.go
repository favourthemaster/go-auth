// Package config provides configuration management for the application.
package config

// DBConfig holds database configuration values.
type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
}

// MailerConfig holds mailer configuration values.
type MailerConfig struct {
	Sender string
	//TODO: ADD THE REMAINING THINGS HERE
}

// RedisConfig holds Redis configuration values.
type RedisConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database int
	SSLMode  string
}
