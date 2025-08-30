package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"

	handlers "github.com/Arpitmovers/reviewservice/internal/handlers/dto"
	s3 "github.com/Arpitmovers/reviewservice/internal/repository/aws"
	"github.com/Arpitmovers/reviewservice/internal/repository/mq"
	"github.com/Arpitmovers/reviewservice/internal/repository/redis"
	"github.com/go-playground/validator"
	// "gorm.io/gorm"
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

func GetPublisher(amqpConn *mq.AmqpConnection) *mq.Publisher {
	publisherOnce.Do(func() {
		var err error
		publisherInstance, err = mq.NewPublisher(amqpConn, "reviews", "direct")
		if err != nil {
			fmt.Printf("failed to init publisher: %v\n", err)
		}
	})
	return publisherInstance
}

func GetSubscriber(amqpConn *mq.AmqpConnection) *mq.Consumer {
	subscriberOnce.Do(func() {
		var err error
		subscriberInstance, err = mq.NewConsumer(amqpConn, "reviewQueue", "reviews", "review.created")
		if err != nil {
			fmt.Printf("failed to init subscriber: %v\n", err)
		}
	})

	return subscriberInstance
}

func (h *ReviewHandler) TriggerReviewInjest(w http.ResponseWriter, r *http.Request) {

	var requestBody handlers.ProcessReviewRequest

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := validateRequest(requestBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(handlers.APIResponse{
			ErrorMsg: err.Error(),
			Success:  false,
		})
		return
	}

	fmt.Println("requestBody", requestBody)

	files, err := h.S3.ListFiles("reviews")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(handlers.APIResponse{
			ErrorMsg: "Failed to list files in S3 bucket",
			Success:  false,
		})
		fmt.Printf("Error listing files: %v", err)
		return
	}

	json.NewEncoder(w).Encode(handlers.APIResponse{
		ErrorMsg: "",
		Success:  true,
	})

	totalFiles := len(files)

	if totalFiles == 0 {
		json.NewEncoder(w).Encode(handlers.APIResponse{
			ErrorMsg: "no files found",
			Success:  false,
		})
		return
	}

	fmt.Println("totalFiles count ", totalFiles)
	numCPU := runtime.NumCPU()

	h.processFiles(files, numCPU)

}

func (h *ReviewHandler) processFiles(files []string, workerCount int) {
	if len(files) == 0 {
		return
	}
	if workerCount > len(files) {
		workerCount = len(files)
	}

	jobs := make(chan string, len(files))
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for fileName := range jobs {
				fmt.Printf("Worker %d picked file: %s\n", workerID, fileName)
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
	fmt.Printf("Processing file: %s", fileName)

	stream, err := h.S3.GetFileStream(fileName)
	if err != nil {
		fmt.Printf("Failed to get stream for %s: %v", fileName, err)
		return
	}
	defer stream.Close()

	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {

		line := scanner.Bytes()
		var review handlers.Review

		if err := json.Unmarshal(line, &review); err != nil {
			fmt.Printf("Invalid JSON in %s: %v", fileName, err)
			continue
		}

		if _, ok := h.validateReview(review); ok {
			// 2. Create Publisher (say exchange = "reviews", type = "direct")

			err = h.Publisher.Publish("review.created", line)
			if err != nil {
				fmt.Printf("failed to publish message: %v", err)
			} else {
				fmt.Println("message published successfully")
			}

		}
	}

	fmt.Print("all data published for", fileName)

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file %s: %v", fileName, err)
	}

}

func (h *ReviewHandler) validateReview(review handlers.Review) (error, bool) {
	if review.HotelID == 0 {
		fmt.Println(" missing HotelID in ", review)

		return fmt.Errorf("hotelId is required"), false
	}
	if review.Comment.HotelReviewID == 0 {
		fmt.Println(" missing HotelReviewID in ", review)
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

// 1 . get total   file count from s3
// spin  no of go routinesbased on no of cpus(B)
// B  go routines will parse fileand do file validation and publis
// noof iternations  needed = totalFileCnt/  cpus

// in each go rountine start reading file and validation  once validation passed add it  to queue
//, andpublishing to amqp,once  B os

//idempotency key: {fileName:"done/failed/inprogress"}  // ttl  of 30 days --> 1st check
// to resume in  case server crash: fileName + "count" : sucess / failed  , on next api call itshould resume fromsame  position
/*
	db wrute fail <> --> push to dead letterq
	s3 read failed --> update   {fileName:failed}
	inital connect  failed --> dea
*/

// api payload  validation

// totalIterations = totalFileCnt /  parallelismFactor
// each go routine,
// read file , parse record
// publish event , after getting ack, update status as processed
