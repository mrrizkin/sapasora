package storage

import (
	"github.com/gofiber/storage/minio"
)

type Minio struct {
	*minio.Storage
}

func NewMinio(config *StorageConfig) *Minio {
	minioConfig := config.Minio
	store := minio.New(minio.Config{
		Bucket:   minioConfig.Bucket,
		Endpoint: minioConfig.Endpoint,
		Region:   minioConfig.Region,
		Token:    minioConfig.Token,
		Secure:   minioConfig.Secure,
		Reset:    minioConfig.Reset,
		Credentials: minio.Credentials{
			AccessKeyID:     minioConfig.AccessKeyID,
			SecretAccessKey: minioConfig.SecretAccessKey,
		},
	})
	return &Minio{
		Storage: store,
	}
}
