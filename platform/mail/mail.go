// Package mail provides mail service for the application
package mail

import (
	"context"
)

// Email represents an email message
type Email struct {
	From        string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	Text        string
	HTML        string
	Attachments []Attachment
	Headers     map[string]string
}

// Attachment represents an email attachment
type Attachment struct {
	Filename    string
	ContentType string
	Content     []byte
	Inline      bool // Whether this is an inline attachment (e.g., for embedded images)
}

// Mailer interface defines the contract for email delivery
type Mailer interface {
	Deliver(email Email) error
	DeliverWithContext(ctx context.Context, email Email) error
}

// MailBuilder provides a fluent interface for building emails
type MailBuilder struct {
	email Email
}

// NewMail creates a new mail builder
func NewMail() *MailBuilder {
	return &MailBuilder{
		email: Email{
			To:          []string{},
			Cc:          []string{},
			Bcc:         []string{},
			Attachments: []Attachment{},
			Headers:     make(map[string]string),
		},
	}
}

// From sets the sender of the email
func (mb *MailBuilder) From(from string) *MailBuilder {
	mb.email.From = from
	return mb
}

// To adds recipients to the email
func (mb *MailBuilder) To(to ...string) *MailBuilder {
	mb.email.To = append(mb.email.To, to...)
	return mb
}

// Cc adds carbon copy recipients to the email
func (mb *MailBuilder) Cc(cc ...string) *MailBuilder {
	mb.email.Cc = append(mb.email.Cc, cc...)
	return mb
}

// Bcc adds blind carbon copy recipients to the email
func (mb *MailBuilder) Bcc(bcc ...string) *MailBuilder {
	mb.email.Bcc = append(mb.email.Bcc, bcc...)
	return mb
}

// Subject sets the subject of the email
func (mb *MailBuilder) Subject(subject string) *MailBuilder {
	mb.email.Subject = subject
	return mb
}

// Text sets the plain text body of the email
func (mb *MailBuilder) Text(text string) *MailBuilder {
	mb.email.Text = text
	return mb
}

// HTML sets the HTML body of the email
func (mb *MailBuilder) HTML(html string) *MailBuilder {
	mb.email.HTML = html
	return mb
}

// AttachFile adds an attachment from a file path
func (mb *MailBuilder) AttachFile(filename string, contentType string) *MailBuilder {
	// This will be implemented to read the file and add it as an attachment
	// For now, we'll just set up the structure
	attachment := Attachment{
		Filename:    filename,
		ContentType: contentType,
		Inline:      false,
	}
	mb.email.Attachments = append(mb.email.Attachments, attachment)
	return mb
}

// AttachBytes adds an attachment from bytes
func (mb *MailBuilder) AttachBytes(
	filename string,
	contentType string,
	content []byte,
) *MailBuilder {
	attachment := Attachment{
		Filename:    filename,
		ContentType: contentType,
		Content:     content,
		Inline:      false,
	}
	mb.email.Attachments = append(mb.email.Attachments, attachment)
	return mb
}

// AddHeader adds a custom header to the email
func (mb *MailBuilder) AddHeader(key, value string) *MailBuilder {
	mb.email.Headers[key] = value
	return mb
}

// Build returns the constructed email
func (mb *MailBuilder) Build() Email {
	return mb.email
}

// MailService interface for the main mail service
type MailService interface {
	Mailer
	Queue(email Email) error
	QueueWithContext(ctx context.Context, email Email) error
}

// KangPaket interface remains for compatibility
type KangPaket interface {
	Mailer
}
