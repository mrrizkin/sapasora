// Package mail provides file-based mailer for development
package mail

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileMailer saves emails to files for development and debugging
type FileMailer struct {
	path string
}

// FileMailerConfig holds configuration for file mailer
type FileMailerConfig struct {
	Path string // Directory where emails will be saved
}

// NewFileMailer creates a new file-based mailer
func NewFileMailer(config FileMailerConfig) *FileMailer {
	// Ensure the directory exists
	if config.Path == "" {
		config.Path = "./storage/mail" // Default path
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(config.Path, 0755); err != nil {
		panic(fmt.Sprintf("failed to create mail directory: %v", err))
	}

	return &FileMailer{
		path: config.Path,
	}
}

// Deliver saves the email to a file
func (f *FileMailer) Deliver(email Email) error {
	return f.DeliverWithContext(context.Background(), email)
}

// DeliverWithContext saves the email to a file with context
func (f *FileMailer) DeliverWithContext(ctx context.Context, email Email) error {
	// Create a filename based on timestamp and subject
	timestamp := time.Now().Format("20060102_150405")
	sanitizedSubject := sanitizeFilename(email.Subject)
	filename := fmt.Sprintf("email_%s_%s.json", timestamp, sanitizedSubject)

	filepath := filepath.Join(f.path, filename)

	// Create email payload for file
	emailPayload := struct {
		Timestamp time.Time `json:"timestamp"`
		Email     Email     `json:"email"`
	}{
		Timestamp: time.Now(),
		Email:     email,
	}

	// Marshal email to JSON
	data, err := json.MarshalIndent(emailPayload, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal email to JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write email to file: %w", err)
	}

	return nil
}

// Helper function to sanitize filename
func sanitizeFilename(s string) string {
	// Replace invalid characters for filenames
	replacements := []string{
		"/",
		"_",
		"\\",
		"_",
		":",
		"_",
		"*",
		"_",
		"?",
		"_",
		"\"",
		"_",
		"<",
		"_",
		">",
		"_",
		"|",
		"_",
		" ",
		"_",
	}

	result := s
	for i := 0; i < len(replacements); i += 2 {
		result = strings.ReplaceAll(result, replacements[i], replacements[i+1])
	}

	// Limit length to prevent issues with filesystems
	if len(result) > 100 {
		result = result[:100]
	}

	return result
}
