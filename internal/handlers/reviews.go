package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"

	handlers "github.com/Arpitmovers/reviewservice/internal/handlers/dto"
	logger "github.com/Arpitmovers/reviewservice/internal/logging"
	s3 "github.com/Arpitmovers/reviewservice/internal/repository/aws"
	"github.com/Arpitmovers/reviewservice/internal/repository/mq"
	"github.com/Arpitmovers/reviewservice/internal/repository/redis"
	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

var validate = validator.New()

func validateRequest(req handlers.ProcessReviewRequest) error {
	return validate.Struct(req)
}

var (
	publisherInstance  *mq.Publisher
	publisherOnce      sync.Once
	subscriberInstance *mq.Consumer
	subscriberOnce     sync.Once
)

type ReviewHandler struct {
	Amqp      *mq.AmqpConnection
	S3        *s3.S3Storage
	Redis     *redis.RedisCache
	Publisher *mq.Publisher
	Consumer  *mq.Consumer
	// DB        *gorm.DB
}

// ---- publisher/subscriber ----

func GetPublisher(amqpConn *mq.AmqpConnection) *mq.Publisher {
	publisherOnce.Do(func() {
		var err error
		publisherInstance, err = mq.NewPublisher(amqpConn, "reviews", "direct")
		if err != nil {
			logger.Logger.Error("failed to init publisher", zap.Error(err))
		} else {
			logger.Logger.Info("publisher initialized", zap.String("exchange", "reviews"))
		}
	})
	return publisherInstance
}

func GetSubscriber(amqpConn *mq.AmqpConnection) *mq.Consumer {
	subscriberOnce.Do(func() {
		var err error
		subscriberInstance, err = mq.NewConsumer(amqpConn, "reviewQueue", "reviews", "review.created")
		if err != nil {
			logger.Logger.Error("failed to init subscriber", zap.Error(err))
		} else {
			logger.Logger.Info("subscriber initialized")
		}
	})
	return subscriberInstance
}

// ---- handler ----

func (h *ReviewHandler) TriggerReviewInjest(w http.ResponseWriter, r *http.Request) {
	var requestBody handlers.ProcessReviewRequest

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		logger.Logger.Error("invalid request payload", zap.Error(err))
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := validateRequest(requestBody); err != nil {
		logger.Logger.Warn("request validation failed", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handlers.APIResponse{
			ErrorMsg: err.Error(),
			Success:  false,
		})
		return
	}

	logger.Logger.Info("received review ingest request", zap.Any("requestBody", requestBody))

	files, err := h.S3.ListFiles("reviews")
	if err != nil {
		logger.Logger.Error("failed to list files in S3", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(handlers.APIResponse{
			ErrorMsg: "Failed to list files in S3 bucket",
			Success:  false,
		})
		return
	}

	totalFiles := len(files)
	if totalFiles == 0 {
		logger.Logger.Warn("no files found in S3 bucket", zap.String("prefix", "reviews"))
		json.NewEncoder(w).Encode(handlers.APIResponse{
			ErrorMsg: "no files found",
			Success:  false,
		})
		return
	}

	logger.Logger.Info("files found for processing", zap.Int("count", totalFiles))

	h.processFiles(files)

	json.NewEncoder(w).Encode(handlers.APIResponse{
		ErrorMsg: "",
		Success:  true,
	})
}

func (h *ReviewHandler) processFiles(files []string) {
	if len(files) == 0 {
		return
	}
	workerCount := runtime.NumCPU()

	if workerCount > len(files) {
		workerCount = len(files)
	}

	logger.Logger.Info("workerCount is ", zap.Int("workerCount", workerCount))
	jobs := make(chan string, len(files))
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for fileName := range jobs {
				logger.Logger.Debug("worker picked file", zap.Int("workerID", workerID), zap.String("file", fileName))
				h.processFile(fileName)
			}
		}(i + 1)
	}

	for _, fileName := range files {
		jobs <- fileName
	}
	close(jobs)

	wg.Wait()
}

func (h *ReviewHandler) processFile(fileName string) {
	logger.Logger.Info("processing file", zap.String("file", fileName))

	stream, err := h.S3.GetFileStream(fileName)
	if err != nil {
		logger.Logger.Error("failed to get file stream", zap.String("file", fileName), zap.Error(err))
		return
	}
	defer stream.Close()

	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		line := scanner.Bytes()
		var review handlers.Review

		if err := json.Unmarshal(line, &review); err != nil {
			logger.Logger.Warn("invalid JSON line, ignoring line", zap.ByteString("line", line), zap.Error(err))
			continue
		}

		if _, ok := h.validateReview(review); ok {

			if err := h.Publisher.PublishSafe("review.created", line); err != nil {
				logger.Logger.Error("failed to publish message", zap.Error(err), zap.String("file", fileName))
			} else {
				logger.Logger.Debug("message published successfully", zap.String("file", fileName))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Logger.Error("error reading file", zap.String("file", fileName), zap.Error(err))
	} else {
		logger.Logger.Info("all data published", zap.String("file", fileName))
	}
}

func (h *ReviewHandler) validateReview(review handlers.Review) (error, bool) {
	if review.HotelID == 0 {
		logger.Logger.Warn("missing HotelID", zap.Any("review", review))
		return fmt.Errorf("hotelId is required"), false
	}
	if review.Comment.HotelReviewID == 0 {
		logger.Logger.Warn("missing HotelReviewID", zap.Any("review", review))
		return fmt.Errorf("missing HotelReviewID"), false
	}

	return nil, true
}

/*
failure scnearios.

1.  s3 list  failed --> api return  500  or ->exponentional backoff
2. s3 reading file failed (pipereading) -- >  exponentialbackoff ( need to read from same position) (savethe  position  inredis

fileName:{ status:" ", count :int})
3. db failed --> park to dead letter  queue.

*/

//idempotency key: {fileName:"done/failed/inprogress"}  // ttl  of 30 days --> 1st check
// to resume in  case server crash: fileName + "count" : sucess / failed  , on next api call itshould resume fromsame  position
/*
	db wrute fail <> --> push to dead letterq
	s3 read failed --> update   {fileName:failed}
	inital connect  failed --> dea
*/
