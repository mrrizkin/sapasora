// Package mail provides configuration for the mail service
package mail

import (
	"sapasora/platform/config"
	"fmt"
	"time"
)

// MailConfig holds all mail configuration settings
type MailConfig struct {
	Driver string           // Mail driver: smtp, file, memory
	SMTP   SMTPConfig       // SMTP configuration
	File   FileMailerConfig // File mailer configuration
	Queue  MailQueueConfig  // Queue configuration
}

// LoadConfig loads mail configuration from the application config
func LoadConfig(cfg config.Config) MailConfig {
	driver := cfg.GetString("mail.driver", "smtp")

	// SMTP Configuration
	smtpConfig := SMTPConfig{
		Host:       cfg.GetString("mail.smtp.host", "localhost"),
		Port:       cfg.GetInt("mail.smtp.port", 587),
		Username:   cfg.GetString("mail.smtp.username", ""),
		Password:   cfg.GetString("mail.smtp.password", ""),
		From:       cfg.GetString("mail.smtp.from", "noreply@example.com"),
		Encryption: cfg.GetString("mail.smtp.encryption", "tls"), // "ssl", "tls", or "none"
		SkipSSL:    cfg.GetBool("mail.smtp.skip_ssl", false),
	}

	// File Configuration
	fileConfig := FileMailerConfig{
		Path: cfg.GetString("mail.file.path", "./storage/mail"),
	}

	// Queue Configuration
	queueConfig := MailQueueConfig{
		Workers:    cfg.GetInt("mail.queue.workers", 1),
		BufferSize: cfg.GetInt("mail.queue.buffer_size", 100),
		MaxRetries: cfg.GetInt("mail.queue.max_retries", 3),
		Timeout:    cfg.GetDuration("mail.queue.timeout", 30*time.Second),
	}

	return MailConfig{
		Driver: driver,
		SMTP:   smtpConfig,
		File:   fileConfig,
		Queue:  queueConfig,
	}
}

// NewMailerFromConfig creates a mailer based on the configuration
func NewMailerFromConfig(cfg MailConfig) (Mailer, error) {
	switch cfg.Driver {
	case "smtp":
		return NewSMTPMailer(cfg.SMTP), nil
	case "file":
		return NewFileMailer(cfg.File), nil
	case "memory":
		return NewMemoryMailer(), nil
	default:
		return nil, fmt.Errorf("unsupported mail driver: %s", cfg.Driver)
	}
}

// NewMailServiceFromConfig creates a mail service based on the configuration
func NewMailServiceFromConfig(cfg config.Config) (MailService, error) {
	mailConfig := LoadConfig(cfg)

	mailer, err := NewMailerFromConfig(mailConfig)
	if err != nil {
		return nil, err
	}

	// Create and return mail service with queue
	service := NewMailService(mailer, mailConfig.Queue)
	return service, nil
}
