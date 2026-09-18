package config

import "time"

type Storage struct {
	driver string `name:"driver" env:"STORAGE_DRIVER,default=local"`

	Local struct {
		Path string `name:"path" env:"STORAGE_LOCAL_PATH,default=./storage"`
	} `name:"local"`

	Minio struct {
		Bucket          string `name:"bucket"   env:"STORAGE_MINIO_BUCKET"`
		Endpoint        string `name:"endpoint" env:"STORAGE_MINIO_ENDPOINT"`
		Region          string `name:"region"   env:"STORAGE_MINIO_REGION"`
		Token           string `name:"token"    env:"STORAGE_MINIO_TOKEN"`
		Secure          bool   `name:"secure"   env:"STORAGE_MINIO_SECURE,default=false"`
		Reset           bool   `name:"reset"    env:"STORAGE_MINIO_RESET,default=false"`
		AccessKeyID     string `name:"access"   env:"STORAGE_MINIO_ACCESS"`
		SecretAccessKey string `name:"secret"   env:"STORAGE_MINIO_SECRET"`
	} `name:"minio"`

	S3 struct {
		Bucket          string        `name:"bucket"       env:"STORAGE_S3_BUCKET"`
		Endpoint        string        `name:"endpoint"     env:"STORAGE_S3_ENDPOINT"`
		Region          string        `name:"region"       env:"STORAGE_S3_REGION"`
		RequestTimeout  time.Duration `name:"timeout"      env:"STORAGE_S3_TIMEOUT"`
		Reset           bool          `name:"reset"        env:"STORAGE_S3_RESET"`
		MaxAttempts     int           `name:"max_attempts" env:"STORAGE_S3_MAX_ATTEMPTS,default=3"`
		AccessKey       string        `name:"access"       env:"STORAGE_S3_ACCESS"`
		SecretAccessKey string        `name:"secret"       env:"STORAGE_S3_SECRET"`
	} `name:"s3"`
}
