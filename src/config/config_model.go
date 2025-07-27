package config

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
}

type MailerConfig struct {
	Sender string
	//TODO: ADD THE REMAINING THINGS HERE
}

type RedisConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database int
	SSLMode  string
}
