package s3

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/Arpitmovers/reviewservice/internal/config"
	logger "github.com/Arpitmovers/reviewservice/internal/logging"
	"github.com/aws/aws-sdk-go-v2/aws"
	s3Config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
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
			s3Config.WithRegion(cfg.AwsRegion),
			s3Config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(cfg.AwsAccessKey, cfg.AwsSecretKey, ""),
			),
		)
		if err != nil {
			logger.Logger.Fatal("unable to load AWS SDK config",
				zap.Error(err),
				zap.String("region", cfg.AwsRegion),
			)
		}

		s3Client = s3.NewFromConfig(s3Cfg)
		logger.Logger.Info("S3 client initialized",
			zap.String("region", cfg.AwsRegion),
			zap.String("bucket", "hotelservice"),
		)
	})

	return &S3Storage{
		client:           s3Client,
		reviewFileBucket: "hotelservice",
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
			logger.Logger.Error("failed to list files from S3",
				zap.String("bucket", s.reviewFileBucket),
				zap.String("prefix", path),
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to list files: %w", err)
		}

		for _, obj := range page.Contents {
			if *obj.Size != int64(0) { // skip empty "folders"
				fileNames = append(fileNames, *obj.Key)
			}
		}
	}

	logger.Logger.Info("S3 files listed successfully",
		zap.String("bucket", s.reviewFileBucket),
		zap.String("prefix", path),
		zap.Int("file_count", len(fileNames)),
	)

	return fileNames, nil
}

func (s *S3Storage) GetFileStream(fileName string) (io.ReadCloser, error) {
	output, err := s.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(s.reviewFileBucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		logger.Logger.Error("failed to get S3 object",
			zap.String("bucket", s.reviewFileBucket),
			zap.String("file", fileName),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get file %s: %w", fileName, err)
	}

	logger.Logger.Info("S3 file stream opened",
		zap.String("bucket", s.reviewFileBucket),
		zap.String("file", fileName),
	)

	return output.Body, nil
}
