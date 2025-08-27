package internal

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	s3Client *s3.Client
	once     sync.Once
)

func GetS3Client() *s3.Client {
	once.Do(func() {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			panic(fmt.Sprintf("unable to load AWS SDK config, %v", err))
		}
		s3Client = s3.NewFromConfig(cfg)
	})
	return s3Client
}

func DownloadFile(bucket, key, destPath string) error {
	client := GetS3Client()

	resp, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return fmt.Errorf("failed to get object: %w", err)
	}
	defer resp.Body.Close()

	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy object to file: %w", err)
	}

	fmt.Printf("File downloaded from s3://%s/%s → %s\n", bucket, key, destPath)
	return nil
}
