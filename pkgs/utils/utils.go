package utils

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

func GetEnv(key string) string {
	env := os.Getenv(key)
	if env == "" {
		panic("Environment variable not set: " + key)
	}
	return env
}
