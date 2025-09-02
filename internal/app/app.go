package app

import (
	"fmt"
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
	logger "github.com/Arpitmovers/reviewservice/internal/logging"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
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

func (a *App) Initialize(cfg *config.Config) {

	logger.InitLogger()
	defer logger.Logger.Sync()

	client, err := s3.GetS3Client(cfg)
	if err != nil {
		logger.Logger.Error("unable to get S3 client", zap.Error(err))
	}
	a.s3Client = client

	a.amqpConnect = a.setupAmqp(cfg)
	a.amqpPubliser = handlers.GetPublisher(a.amqpConnect)
	a.amqpConsumer = handlers.GetSubscriber(a.amqpConnect)
	var redisError error
	a.redisClient, redisError = redis.GetRedisClient(cfg)

	if redisError != nil {
		logger.Logger.Error("error in GetRedisClient ", zap.Error(err))
	}

	var dbErr error
	a.dbConnect, dbErr = db.NewDBConnect(cfg)
	if dbErr != nil {
		logger.Logger.Error("error in NewDBConnect ", zap.Error(err))
	}

	a.Router = mux.NewRouter()
	a.setRouters(cfg)

	// setup services
	reviewRepo := models.NewReviewRepository(a.dbConnect)
	reviewService := services.NewReviewService(reviewRepo)
	reviewConsumer := handlers.NewReviewConsumer(reviewService)

	rmqError := a.amqpConsumer.Consume(reviewConsumer.ConsumeReview)
	if rmqError != nil {
		logger.Logger.Error("failed to start consumer", zap.Error(rmqError))
	}
}

func (a *App) setupAmqp(cfg *config.Config) *mq.AmqpConnection {
	amqpURL := fmt.Sprintf(
		"amqp://%s:%s@localhost:%s/%s",
		cfg.AmqpUserName, cfg.AmqpPwd, cfg.AmqpPort, cfg.AmqpVhost)

	conn, err := mq.NewConnection(amqpURL)
	if err != nil {
		logger.Logger.Error("failed to connect to AMQP", zap.Error(err))
		return nil
	}

	logger.Logger.Info("Successfully connected to AMQP", zap.String("url", amqpURL))
	return conn
}

func (a *App) setRouters(cfg *config.Config) {
	reviewHandler := &handlers.ReviewHandler{
		S3:        a.s3Client,
		Amqp:      a.amqpConnect,
		Redis:     a.redisClient,
		Publisher: a.amqpPubliser,
		Consumer:  a.amqpConsumer,
	}

	a.Router.Handle(
		"/v1/reviews/injest",
		auth.JWTAuthMiddleware(cfg, http.HandlerFunc(reviewHandler.TriggerReviewInjest)),
	).Methods(http.MethodPost)

	a.Router.HandleFunc("/login", auth.LoginHandler(cfg)).Methods("POST")
}

func (a *App) Run(host string) {
	logger.Logger.Info("starting server", zap.String("host", host))
	http.ListenAndServe(host, a.Router)
}
