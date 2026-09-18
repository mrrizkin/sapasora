package storage

import (
	"github.com/gofiber/storage/s3/v2"
)

type S3 struct {
	*s3.Storage
}

func NewS3(config *StorageConfig) *S3 {
	s3Config := config.S3
	store := s3.New(s3.Config{
		Bucket:         s3Config.Bucket,
		Endpoint:       s3Config.Endpoint,
		Region:         s3Config.Region,
		RequestTimeout: s3Config.RequestTimeout,
		Reset:          s3Config.Reset,
		MaxAttempts:    s3Config.MaxAttempts,
		Credentials: s3.Credentials{
			AccessKey:       s3Config.AccessKey,
			SecretAccessKey: s3Config.SecretAccessKey,
		},
	})
	return &S3{
		Storage: store,
	}
}
