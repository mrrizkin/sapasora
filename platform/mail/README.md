# Mail Service Documentation

## Overview

Sapasora's mail service is a comprehensive email delivery system inspired by Elixir Phoenix's Swoosh library. It provides a fluent interface for building emails, multiple delivery backends, and robust queuing capabilities.

## Features

- **Fluent Email Builder**: Easy-to-use chainable methods for constructing emails
- **Multiple Delivery Backends**: SMTP, file-based storage, and in-memory options
- **Email Queuing**: Asynchronous email delivery with retry logic
- **Attachments Support**: Attach files with proper MIME types
- **Configuration Support**: Environment-based configuration
- **Testing Utilities**: Memory mailer for testing

## Installation

The mail service is part of the Sapasora platform and doesn't require additional installation.

## Usage

### Basic Email Delivery

```go
import "sapasora/platform/mail"

// Create an email using the fluent builder
email := mail.NewMail().
    From("sender@example.com").
    To("recipient@example.com").
    Subject("Hello from Sapasora").
    Text("This is a plain text message.").
    HTML("<h1>This is an HTML message</h1>").
    Build()

// Create an SMTP mailer
mailer := mail.NewSMTPMailer(mail.SMTPConfig{
    Host:     "smtp.example.com",
    Port:     587,
    Username: "your-username",
    Password: "your-password",
    From:     "sender@example.com",
    Encryption: "tls",
})

// Send the email
err := mailer.Deliver(email)
if err != nil {
    // handle error
}
```

### Using Configuration

```go
import (
    "sapasora/platform/mail"
    "sapasora/platform/config"
)

// Load mail configuration from application config
cfg := config.Load() // assume this loads your config
mailService, err := mail.NewMailServiceFromConfig(cfg)
if err != nil {
    // handle error
}

// Create and send an email
email := mail.NewMail().
    To("user@example.com").
    Subject("Welcome!").
    HTML("<h1>Welcome to our service!</h1>").
    Build()

err = mailService.Deliver(email)
if err != nil {
    // handle error
}
```

### Email with Attachments

```go
// Create email with file attachment
email := mail.NewMail().
    From("sender@example.com").
    To("recipient@example.com").
    Subject("Report").
    Text("Please find the attached report.").
    AttachFile("./report.pdf", "application/pdf").
    Build()

// Or attach from bytes
content := []byte("file content")
email = mail.NewMail().
    From("sender@example.com").
    To("recipient@example.com").
    Subject("Data File").
    Text("Please find the attached data.").
    AttachBytes("data.txt", "text/plain", content).
    Build()
```

### Using Mail Queue

```go
// Create a mail service with queue support
mailer := mail.NewSMTPMailer(smtpConfig)
queueConfig := mail.MailQueueConfig{
    Workers:    2,
    BufferSize: 50,
    MaxRetries: 3,
    Timeout:    30 * time.Second,
}
mailService := mail.NewMailService(mailer, queueConfig)

// Queue an email for asynchronous delivery
email := mail.NewMail().
    To("user@example.com").
    Subject("Welcome!").
    HTML("<h1>Welcome!</h1>").
    Build()

err := mailService.Queue(email)
if err != nil {
    // handle error
}

// Get queue statistics
stats := mailService.GetStats()
fmt.Printf("Processed: %d, Failed: %d, Queue Size: %d\n",
    stats.TotalProcessed, stats.TotalFailed, stats.QueueSize)
```

### Development and Testing Mailers

```go
// File mailer for development - saves emails as JSON files
fileMailer := mail.NewFileMailer(mail.FileMailerConfig{
    Path: "./storage/mail",
})

// Memory mailer for testing
memoryMailer := mail.NewMemoryMailer()
email := mail.NewMail().
    To("test@example.com").
    Subject("Test").
    Text("This is a test").
    Build()

memoryMailer.Deliver(email)

// Retrieve sent emails in tests
sentEmails := memoryMailer.GetEmails()
lastEmail, err := memoryMailer.GetLastEmail()
```

## Mail Drivers

### SMTP Driver

Sends emails via SMTP server. Configuration options:

- `MAIL_SMTP_HOST` - SMTP server host (config key: `mail.smtp.host`)
- `MAIL_SMTP_PORT` - SMTP server port (config key: `mail.smtp.port`)
- `MAIL_SMTP_USERNAME` - Authentication username (config key: `mail.smtp.username`)
- `MAIL_SMTP_PASSWORD` - Authentication password (config key: `mail.smtp.password`)
- `MAIL_SMTP_FROM` - Default sender address (config key: `mail.smtp.from`)
- `MAIL_SMTP_ENCRYPTION` - Encryption type ("ssl", "tls", or "none") (config key: `mail.smtp.encryption`)
- `MAIL_SMTP_SKIP_SSL` - Skip SSL verification (for testing) (config key: `mail.smtp.skip_ssl`)

### File Driver

Saves emails as JSON files for development. Configuration options:

- `MAIL_FILE_PATH` - Directory to save email files (config key: `mail.file.path`)

### Memory Driver

Stores emails in memory, primarily for testing. Configured via the main driver setting:

- `MAIL_DRIVER` - Set to "memory" to use memory mailer (config key: `mail.driver`)

## Queue Configuration

- `MAIL_QUEUE_WORKERS` - Number of worker goroutines (config key: `mail.queue.workers`)
- `MAIL_QUEUE_BUFFER_SIZE` - Size of the queue channel (config key: `mail.queue.buffer_size`)
- `MAIL_QUEUE_MAX_RETRIES` - Maximum retry attempts for failed deliveries (config key: `mail.queue.max_retries`)
- `MAIL_QUEUE_TIMEOUT` - Timeout for email delivery operations (config key: `mail.queue.timeout`)

## Testing

The memory mailer makes testing easy:

```go
func TestEmailSending(t *testing.T) {
    // Create memory mailer for testing
    mailer := mail.NewMemoryMailer()

    // Create and send an email
    email := mail.NewMail().
        To("test@example.com").
        Subject("Test").
        Text("Test message").
        Build()

    err := mailer.Deliver(email)
    assert.NoError(t, err)

    // Verify email was stored
    emails := mailer.GetEmails()
    assert.Equal(t, 1, len(emails))
    assert.Equal(t, "test@example.com", emails[0].To[0])
}
```

## Best Practices

1. **Use Queuing in Production**: Always queue emails in production to avoid blocking HTTP requests
2. **Set Appropriate Timeouts**: Configure delivery timeouts to avoid hanging requests
3. **Test with Memory Mailer**: Use the memory mailer for unit tests
4. **Environment-Specific Configuration**: Use different drivers for different environments
5. **Monitor Queue Stats**: Regularly check queue statistics for failures
6. **Secure Credentials**: Store SMTP credentials in environment variables, not code
