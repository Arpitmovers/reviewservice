package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServicePort  string
	DBUrl        string
	LogLevel     string
	AmqpPort     string
	AmqpVhost    string
	AmqpUserName string
	AmqpPwd      string
	RedisHost    string
	RedisPort    string
	AwsRegion    string
	AwsAccessKey string
	AwsSecretKey string
	DbHost       string
	DbPort       string
	DbName       string
	DbPwd        string
	DbUser       string
	ReviewTable  string
	JwtKey       string
	ApiUser      string
	ApiPwd       string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, falling back to environment variables")
	}

	cfg := &Config{
		ServicePort:  getEnv("PORT", "8080"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		AmqpPort:     getEnv("AMQP_PORT", "5672"),
		AmqpVhost:    getEnv("AMQP_VHOST", "info"),
		AmqpUserName: getEnv("AMQP_USERNAME", "info"),
		AmqpPwd:      getEnv("AMQP_PWD", "info"),
		RedisPort:    getEnv("AMQP_PWD", "6379"),
		RedisHost:    getEnv("REDIS_HOST", "6379"),
		AwsAccessKey: getEnv("AWS_ACCESS_KEY_ID", "6379"),
		AwsSecretKey: getEnv("AWS_SECRET_ACCESS_KEY", "6379"),

		DbHost:      getEnv("MARIA_HOST", "localhost"),
		DbPort:      getEnv("MARIA_PORT", "3306"),
		DbName:      getEnv("REVIEW_DB", "zuzu_db"),
		DbPwd:       getEnv("REVIEW_DBPWD", " "),
		DbUser:      getEnv("REVIEW_USER", " "),
		ReviewTable: getEnv("REVIEW_TABLE", " "),
		JwtKey:      getEnv("JWT_KEY", " "),
		ApiUser:     getEnv("API_USER", " "),
		ApiPwd:      getEnv("API_PWD", " "),
	}
	return cfg

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
