package functions

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(s string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	v := os.Getenv(s)

	return v
}
