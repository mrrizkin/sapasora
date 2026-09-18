package config

import "time"

type Mail struct {
	Driver string `name:"driver" env:"MAIL_DRIVER,default=smtp"`

	SMTP struct {
		Host       string `name:"host"       env:"MAIL_SMTP_HOST,default=localhost"`
		Port       int    `name:"port"       env:"MAIL_SMTP_PORT,default=587"`
		Username   string `name:"username"   env:"MAIL_SMTP_USERNAME"`
		Password   string `name:"password"   env:"MAIL_SMTP_PASSWORD"`
		From       string `name:"from"       env:"MAIL_SMTP_FROM,default=noreply@localhost"`
		Encryption string `name:"encryption" env:"MAIL_SMTP_ENCRYPTION,default=tls"`
		SkipSSL    bool   `name:"skip_ssl"   env:"MAIL_SMTP_SKIP_SSL,default=false"`
	} `name:"smtp"`

	File struct {
		Path string `name:"path" env:"MAIL_FILE_PATH,default=./storage/mail"`
	} `name:"file"`

	Queue struct {
		Workers    int           `name:"workers"     env:"MAIL_QUEUE_WORKERS,default=1"`
		BufferSize int           `name:"buffer_size" env:"MAIL_QUEUE_BUFFER_SIZE,default=100"`
		MaxRetries int           `name:"max_retries" env:"MAIL_QUEUE_MAX_RETRIES,default=3"`
		Timeout    time.Duration `name:"timeout"     env:"MAIL_QUEUE_TIMEOUT,default=30s"`
	} `name:"queue"`
}
