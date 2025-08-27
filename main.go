package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Arpitmovers/reviewservice/internal/rabbitmq/mq"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	// r.GET("/user/:name", func(c *gin.Context) {
	// 	user := c.Params.ByName("name")
	// 	value, ok := db[user]
	// 	if ok {
	// 		c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
	// 	} else {
	// 		c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
	// 	}
	// })

	return r
}

func main() {

	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		log.Fatal("AMQP_URL environment variable is not set")
	}

	// Establish a connection using NewConnection from the mq package
	conn, err := mq.NewConnection(amqpURL)
	if err != nil {
		log.Fatalf("Failed to connect to AMQP: %v", err)
	}
	defer conn.Close()

	log.Println("Successfully connected to AMQP")
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")

}
