package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Arpitmovers/reviewservice/internal/config"
	"github.com/Arpitmovers/reviewservice/internal/handlers"
	s3 "github.com/Arpitmovers/reviewservice/internal/repository/aws"
	"github.com/Arpitmovers/reviewservice/internal/repository/db"
	"github.com/Arpitmovers/reviewservice/internal/repository/models"
	"github.com/Arpitmovers/reviewservice/internal/repository/mq"
	"github.com/Arpitmovers/reviewservice/internal/repository/redis"
	services "github.com/Arpitmovers/reviewservice/internal/service"

	"github.com/Arpitmovers/reviewservice/internal/auth"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type App struct {
	Router       *mux.Router
	amqpConnect  *mq.AmqpConnection
	s3Client     *s3.S3Storage
	redisClient  *redis.RedisCache
	dbConnect    *gorm.DB
	amqpPubliser *mq.Publisher
	amqpConsumer *mq.Consumer
}

func (a *App) Initialize(config *config.Config) {
	a.s3Client = s3.GetS3Client(config)
	a.amqpConnect = setupAmqp(config)
	a.amqpPubliser = handlers.GetPublisher(a.amqpConnect)
	a.amqpConsumer = handlers.GetSubscriber(a.amqpConnect)

	a.redisClient = redis.GetRedisClient()
	a.dbConnect = db.NewDBConnect(config)
	a.Router = mux.NewRouter()
	a.setRouters(config)

	reviewRepo := models.NewReviewRepository(a.dbConnect)
	reviewService := services.NewReviewService(reviewRepo)
	reviewConsumer := handlers.NewReviewConsumer(reviewService)

	err := a.amqpConsumer.Consume(reviewConsumer.ConsumeReview())

	if err != nil {
		log.Fatalf("failed to start consumer: %v", err)
	}
	//  Block main goroutine so your app keeps running
	//select {}
	//	defer a.amqpConnect.Close()
}

func setupAmqp(cfg *config.Config) *mq.AmqpConnection {
	amqpURL := fmt.Sprintf(
		"amqp://%s:%s@localhost:%s/%s",
		cfg.AmqpUserName, cfg.AmqpPwd, cfg.AmqpPort, cfg.AmqpVhost)

	conn, err := mq.NewConnection(amqpURL)
	if err != nil {
		log.Fatalf("Failed to connect to AMQP: %v", err)

	}
	log.Println("Successfully connected to AMQP")
	return conn
}

func (a *App) setRouters(cfg *config.Config) {
	reviewHandler := &handlers.ReviewHandler{S3: a.s3Client, Amqp: a.amqpConnect, Redis: a.redisClient,
		Publisher: a.amqpPubliser, Consumer: a.amqpConsumer}

	a.Router.Handle(
		"/reviews/injest",
		auth.JWTAuthMiddleware(cfg, http.HandlerFunc(reviewHandler.TriggerReviewInjest)),
	).Methods(http.MethodPost)

	a.Router.HandleFunc("/login", auth.LoginHandler(cfg)).Methods("POST")
	// return reviewHandler
}

func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
