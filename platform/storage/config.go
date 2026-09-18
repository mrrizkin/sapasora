package storage

import (
	"sapasora/platform/config"
	"time"
)

type StorageConfig struct {
	Local LocalConfig
	Minio MinioConfig
	S3    S3Config
}

func NewStorageConfig(cfg config.Config) *StorageConfig {
	return &StorageConfig{
		Local: LocalConfig{
			Path: cfg.GetString("storage.local.path"),
		},
		Minio: MinioConfig{
			Bucket:          cfg.GetString("storage.minio.bucket"),
			Endpoint:        cfg.GetString("storage.minio.endpoint"),
			Region:          cfg.GetString("storage.minio.region"),
			Token:           cfg.GetString("storage.minio.token"),
			Secure:          cfg.GetBool("storage.minio.secure"),
			Reset:           cfg.GetBool("storage.minio.reset"),
			AccessKeyID:     cfg.GetString("storage.minio.access"),
			SecretAccessKey: cfg.GetString("storage.minio.secret"),
		},
		S3: S3Config{
			Bucket:          cfg.GetString("storage.s3.bucket"),
			Endpoint:        cfg.GetString("storage.s3.endpoint"),
			Region:          cfg.GetString("storage.s3.region"),
			RequestTimeout:  cfg.GetDuration("storage.s3.timeout"),
			Reset:           cfg.GetBool("storage.s3.reset"),
			MaxAttempts:     cfg.GetInt("storage.s3.max_attempts"),
			AccessKey:       cfg.GetString("storage.s3.access"),
			SecretAccessKey: cfg.GetString("storage.s3.secret"),
		},
	}
}

type MinioConfig struct {
	Bucket          string
	Endpoint        string
	Region          string
	Token           string
	Secure          bool
	Reset           bool
	AccessKeyID     string
	SecretAccessKey string
}

type S3Config struct {
	Bucket          string
	Endpoint        string
	Region          string
	RequestTimeout  time.Duration
	Reset           bool
	MaxAttempts     int
	AccessKey       string
	SecretAccessKey string
}

type LocalConfig struct {
	Path string
}
