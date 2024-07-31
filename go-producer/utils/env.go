package utils

import (
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string, defaultValue string) string {
	err := godotenv.Load()
	if err != nil {
		return defaultValue
	}

	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}
	return val
}
