package config

import "os"

type Config struct {
	DatabaseURL string
	ServerPort  string
}

func LoadConfig() Config {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://user:secretpassword@localhost/auth?sslmode=disable"
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	return Config{
		DatabaseURL: databaseURL,
		ServerPort:  serverPort,
	}
}
