// Package mail provides mail service for the application
package mail

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/smtp"
	"net/textproto"
	"path/filepath"
	"strings"
	"time"
)

// SMTPConfig holds configuration for SMTP mailer
type SMTPConfig struct {
	Host       string
	Port       int
	Username   string
	Password   string
	From       string
	Encryption string // "ssl", "tls", or "none"
	SkipSSL    bool   // Skip SSL verification (useful for testing)
}

// SMTPMailer implements the Mailer interface for sending emails via SMTP
type SMTPMailer struct {
	config SMTPConfig
}

// NewSMTPMailer creates a new SMTP mailer with the given configuration
func NewSMTPMailer(config SMTPConfig) *SMTPMailer {
	return &SMTPMailer{
		config: config,
	}
}

// Deliver sends an email via SMTP
func (s *SMTPMailer) Deliver(email Email) error {
	return s.DeliverWithContext(context.Background(), email)
}

// DeliverWithContext sends an email via SMTP with context
func (s *SMTPMailer) DeliverWithContext(ctx context.Context, email Email) error {
	// Set up authentication
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	// Create the message
	message, err := buildMessage(email)
	if err != nil {
		return fmt.Errorf("failed to build email message: %w", err)
	}

	// Extract recipient emails
	recipients := make([]string, 0, len(email.To)+len(email.Cc)+len(email.Bcc))
	recipients = append(recipients, email.To...)
	recipients = append(recipients, email.Cc...)
	recipients = append(recipients, email.Bcc...)

	// Connect to the SMTP server
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	var conn smtpClient
	switch s.config.Encryption {
	case "ssl":
		conn, err = newSSLClient(addr, s.config.SkipSSL)
	case "tls":
		conn, err = newTLSClient(addr, s.config.Host, s.config.SkipSSL)
	default:
		conn, err = newPlainClient(addr)
	}
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	// Authenticate
	if err := conn.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}

	// Set sender
	from := email.From
	if from == "" {
		from = s.config.From
	}
	if err := conn.Mail(from); err != nil {
		return fmt.Errorf("SMTP MAIL command failed: %w", err)
	}

	// Set recipients
	for _, recipient := range recipients {
		if err := conn.Rcpt(recipient); err != nil {
			return fmt.Errorf("SMTP RCPT command failed for %s: %w", recipient, err)
		}
	}

	// Send data
	writer, err := conn.Data()
	if err != nil {
		return fmt.Errorf("SMTP DATA command failed: %w", err)
	}

	_, err = writer.Write(message)
	if err != nil {
		writer.Close()
		return fmt.Errorf("failed to write email data: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close email data writer: %w", err)
	}

	return nil
}

// smtpClient interface for different SMTP client types
type smtpClient interface {
	Close() error
	Mail(from string) error
	Rcpt(to string) error
	Data() (io.WriteCloser, error)
	Auth(a smtp.Auth) error
}

// Helper functions to build the email message
func buildMessage(email Email) ([]byte, error) {
	var body bytes.Buffer

	// Write headers
	if email.From != "" {
		body.WriteString(fmt.Sprintf("From: %s\r\n", email.From))
	}

	if len(email.To) > 0 {
		body.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(email.To, ", ")))
	}

	if len(email.Cc) > 0 {
		body.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(email.Cc, ", ")))
	}

	body.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))

	// Add custom headers
	for key, value := range email.Headers {
		body.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	// Determine content type based on whether we have HTML or text
	var contentType string
	if email.HTML != "" && email.Text != "" {
		contentType = "multipart/mixed"
	} else if email.HTML != "" {
		contentType = "text/html"
	} else {
		contentType = "text/plain"
	}

	// If we have attachments or mixed content, we need multipart
	if len(email.Attachments) > 0 || (email.HTML != "" && email.Text != "") {
		return buildMultipartMessage(email)
	}

	// Simple case: just text or HTML
	body.WriteString("MIME-Version: 1.0\r\n")
	if contentType == "text/html" {
		body.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		body.WriteString("\r\n")
		body.WriteString(email.HTML)
	} else {
		body.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		body.WriteString("\r\n")
		body.WriteString(email.Text)
	}

	return body.Bytes(), nil
}

func buildMultipartMessage(email Email) ([]byte, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	boundary := writer.Boundary()

	body.WriteString("MIME-Version: 1.0\r\n")
	body.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", boundary))
	body.WriteString("\r\n")

	// If we have both text and HTML, create an alternative part
	if email.Text != "" && email.HTML != "" {
		altWriter := multipart.NewWriter(nil)
		altBoundary := altWriter.Boundary()

		// Write the alternative part header
		body.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		body.WriteString(
			fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s\r\n", altBoundary),
		)
		body.WriteString("\r\n")

		// Write text part
		body.WriteString(fmt.Sprintf("--%s\r\n", altBoundary))
		body.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		body.WriteString("\r\n")
		body.WriteString(email.Text)
		body.WriteString("\r\n")

		// Write HTML part
		body.WriteString(fmt.Sprintf("--%s\r\n", altBoundary))
		body.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		body.WriteString("\r\n")
		body.WriteString(email.HTML)
		body.WriteString("\r\n")

		// Close alternative boundary
		body.WriteString(fmt.Sprintf("--%s--\r\n", altBoundary))
	} else if email.Text != "" {
		// Just text
		body.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		body.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		body.WriteString("\r\n")
		body.WriteString(email.Text)
		body.WriteString("\r\n")
	} else if email.HTML != "" {
		// Just HTML
		body.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		body.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		body.WriteString("\r\n")
		body.WriteString(email.HTML)
		body.WriteString("\r\n")
	}

	// Add attachments
	for _, attachment := range email.Attachments {
		body.WriteString(fmt.Sprintf("--%s\r\n", boundary))

		h := make(textproto.MIMEHeader)
		if attachment.Inline {
			h.Set("Content-Type", attachment.ContentType)
			h.Set("Content-ID", fmt.Sprintf("<%s>", attachment.Filename))
			h.Set(
				"Content-Disposition",
				fmt.Sprintf("inline; filename=\"%s\"", filepath.Base(attachment.Filename)),
			)
		} else {
			h.Set("Content-Type", attachment.ContentType)
			h.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(attachment.Filename)))
		}
		h.Set("Content-Transfer-Encoding", "base64")

		// Write the headers manually
		for header, values := range h {
			for _, value := range values {
				body.WriteString(fmt.Sprintf("%s: %s\r\n", header, value))
			}
		}
		body.WriteString("\r\n")

		// Write the attachment content
		body.Write(attachment.Content)
		body.WriteString("\r\n")
	}

	// Close the multipart
	body.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	return body.Bytes(), nil
}

// PlainClient wraps the standard net/smtp.Client
type PlainClient struct {
	client *smtp.Client
	conn   net.Conn
}

func newPlainClient(addr string) (*PlainClient, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	client, err := smtp.NewClient(conn, "")
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &PlainClient{
		client: client,
		conn:   conn,
	}, nil
}

func (pc *PlainClient) Close() error {
	pc.client.Quit()
	return pc.conn.Close()
}

func (pc *PlainClient) Mail(from string) error {
	return pc.client.Mail(from)
}

func (pc *PlainClient) Rcpt(to string) error {
	return pc.client.Rcpt(to)
}

func (pc *PlainClient) Data() (io.WriteCloser, error) {
	return pc.client.Data()
}

func (pc *PlainClient) Auth(a smtp.Auth) error {
	return pc.client.Auth(a)
}

// TLSClient for STARTTLS connections
type TLSClient struct {
	client *smtp.Client
	conn   net.Conn
}

func newTLSClient(addr, host string, skipSSL bool) (*TLSClient, error) {
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return nil, err
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Start TLS
	tlsConfig := &tls.Config{
		ServerName: host,
	}
	if skipSSL {
		tlsConfig.InsecureSkipVerify = true
	}

	if err = client.StartTLS(tlsConfig); err != nil {
		client.Quit()
		conn.Close()
		return nil, err
	}

	return &TLSClient{
		client: client,
		conn:   conn,
	}, nil
}

func (tc *TLSClient) Close() error {
	tc.client.Quit()
	return tc.conn.Close()
}

func (tc *TLSClient) Mail(from string) error {
	return tc.client.Mail(from)
}

func (tc *TLSClient) Rcpt(to string) error {
	return tc.client.Rcpt(to)
}

func (tc *TLSClient) Data() (io.WriteCloser, error) {
	return tc.client.Data()
}

func (tc *TLSClient) Auth(a smtp.Auth) error {
	return tc.client.Auth(a)
}

// SSLClient for SSL/TLS connections
type SSLClient struct {
	client *smtp.Client
	conn   *tls.Conn
}

func newSSLClient(addr string, skipSSL bool) (*SSLClient, error) {
	tlsConfig := &tls.Config{}
	if skipSSL {
		tlsConfig.InsecureSkipVerify = true
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return nil, err
	}

	host, _, _ := net.SplitHostPort(addr)
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &SSLClient{
		client: client,
		conn:   conn,
	}, nil
}

func (sc *SSLClient) Close() error {
	sc.client.Quit()
	return sc.conn.Close()
}

func (sc *SSLClient) Mail(from string) error {
	return sc.client.Mail(from)
}

func (sc *SSLClient) Rcpt(to string) error {
	return sc.client.Rcpt(to)
}

func (sc *SSLClient) Data() (io.WriteCloser, error) {
	return sc.client.Data()
}

func (sc *SSLClient) Auth(a smtp.Auth) error {
	return sc.client.Auth(a)
}
