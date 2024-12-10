package utils

import (
	"database/sql"
	"os"
	"time"

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

func ConvertToNullTime(timestamp string) sql.NullTime {
	parsedTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		}
	}
	return sql.NullTime{
		Time:  parsedTime,
		Valid: true,
	}
}
