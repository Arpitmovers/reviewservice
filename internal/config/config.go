package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

}

func main() {
	key := os.Getenv("AWS_ACCESS_KEY_ID")
	log.Println("Using key:", key)
}
