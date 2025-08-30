package s3

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/Arpitmovers/reviewservice/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	s3Config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	s3Client *s3.Client
	once     sync.Once
)

type S3Storage struct {
	client           *s3.Client
	reviewFileBucket string
}

type Storage interface {
	ListFiles(path string) ([]string, error)
	GetFileStream(fileName string) (io.ReadCloser, error)
}

func GetS3Client(cfg *config.Config) *S3Storage {
	once.Do(func() {

		s3Cfg, err := s3Config.LoadDefaultConfig(
			context.TODO(),
			s3Config.WithRegion("ap-south-1"),
			s3Config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(cfg.AwsAccessKey, cfg.AwsSecretKey, ""),
			),
		)
		if err != nil {
			panic(fmt.Sprintf("unable to load AWS SDK config, %v", err))
		}

		s3Client = s3.NewFromConfig(s3Cfg)
		fmt.Println("S3 client initialized")
	})
	return &S3Storage{
		client:           s3Client,
		reviewFileBucket: "hotelservice", // Assuming you have a bucket name in your config
	}
}

func (s *S3Storage) ListFiles(path string) ([]string, error) {

	var fileNames []string
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String("hotelservice"),
		Prefix: aws.String(path),
	}

	paginator := s3.NewListObjectsV2Paginator(s.client, input)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to list files: %w", err)
		}

		for _, obj := range page.Contents {
			//Skip "folders" (keys ending with "/")
			if *obj.Size != int64(0) {
				fileNames = append(fileNames, *obj.Key)
			}
		}
	}

	return fileNames, nil
}

// GetFileStream returns an io.ReadCloser to stream the file content
func (s *S3Storage) GetFileStream(fileName string) (io.ReadCloser, error) {
	output, err := s.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(s.reviewFileBucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get file %s: %w", fileName, err)
	}
	// Caller is responsible for closing the stream
	return output.Body, nil
}
