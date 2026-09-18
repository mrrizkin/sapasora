// Package mail provides utility functions for email composition
package mail

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"slices"
)

// RenderTemplate renders an HTML template to string (placeholder implementation)
// This would integrate with your templating system
func RenderTemplate(templateName string, data any) (string, error) {
	// This is a placeholder - would integrate with actual templating system
	// like the templ package used in the project
	return "", fmt.Errorf("template rendering not implemented yet")
}

// ReadFileAttachment reads a file and returns it as an attachment
func ReadFileAttachment(filePath string) (Attachment, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return Attachment{}, fmt.Errorf("failed to read file: %w", err)
	}

	// Determine content type from file extension
	contentType := "application/octet-stream" // default

	// For images
	if isImageFile(filePath) {
		contentType = fmt.Sprintf("image/%s", getFileExtension(filePath))
	} else if isTextFile(filePath) {
		contentType = "text/plain"
	} else if isHTMLFile(filePath) {
		contentType = "text/html"
	} else if isPDFFile(filePath) {
		contentType = "application/pdf"
	} else if isZipFile(filePath) {
		contentType = "application/zip"
	}

	return Attachment{
		Filename:    filePath,
		ContentType: contentType,
		Content:     content,
		Inline:      false,
	}, nil
}

// ReadFileAttachmentAsInline reads a file and returns it as an inline attachment
func ReadFileAttachmentAsInline(filePath, cid string) (Attachment, error) {
	attachment, err := ReadFileAttachment(filePath)
	if err != nil {
		return Attachment{}, err
	}

	// If a CID (Content-ID) is provided, set it as inline
	if cid != "" {
		attachment.Inline = true
		attachment.Filename = cid
	}

	return attachment, nil
}

// Helper functions to determine file types
func isImageFile(filename string) bool {
	ext := getFileExtension(filename)
	imageExts := []string{"jpg", "jpeg", "png", "gif", "bmp", "webp"}
	return slices.Contains(imageExts, ext)
}

func isTextFile(filename string) bool {
	ext := getFileExtension(filename)
	textExts := []string{"txt", "md", "csv"}
	return slices.Contains(textExts, ext)
}

func isHTMLFile(filename string) bool {
	ext := getFileExtension(filename)
	return ext == "html" || ext == "htm"
}

func isPDFFile(filename string) bool {
	ext := getFileExtension(filename)
	return ext == "pdf"
}

func isZipFile(filename string) bool {
	ext := getFileExtension(filename)
	return ext == "zip"
}

func getFileExtension(filename string) string {
	for i := len(filename) - 1; i >= 0 && filename[i] != '/'; i-- {
		if filename[i] == '.' {
			return filename[i+1:]
		}
	}
	return ""
}

// MailTemplate represents a reusable email template
type MailTemplate struct {
	Subject string
	HTML    string
	Text    string
}

// ApplyTemplate applies a template to a mail builder
func (mb *MailBuilder) ApplyTemplate(template MailTemplate) *MailBuilder {
	mb.email.Subject = template.Subject
	mb.email.HTML = template.HTML
	mb.email.Text = template.Text
	return mb
}

// AttachFromReader adds an attachment from an io.Reader
func (mb *MailBuilder) AttachFromReader(
	reader io.Reader,
	filename, contentType string,
) *MailBuilder {
	content, err := io.ReadAll(reader)
	if err != nil {
		// In a real implementation, we might want to handle this differently
		// For now, we'll add a placeholder
		return mb
	}

	attachment := Attachment{
		Filename:    filename,
		ContentType: contentType,
		Content:     content,
		Inline:      false,
	}

	mb.email.Attachments = append(mb.email.Attachments, attachment)
	return mb
}

// AttachFromBase64 adds an attachment from base64 encoded data
func (mb *MailBuilder) AttachFromBase64(filename, contentType, base64Data string) *MailBuilder {
	content, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		// In a real implementation, we might want to handle this differently
		return mb
	}

	attachment := Attachment{
		Filename:    filename,
		ContentType: contentType,
		Content:     content,
		Inline:      false,
	}

	mb.email.Attachments = append(mb.email.Attachments, attachment)
	return mb
}

// InlineImage adds an inline image to the email
func (mb *MailBuilder) InlineImage(imagePath, cid string) *MailBuilder {
	attachment, err := ReadFileAttachmentAsInline(imagePath, cid)
	if err != nil {
		// In a real implementation, we might want to handle this differently
		return mb
	}

	mb.email.Attachments = append(mb.email.Attachments, attachment)
	return mb
}
