package main

import (
	"fmt"
	"log"
	"time"

	"sapasora/platform/mail"
)

func main() {
	fmt.Println("sapasora Mail Service Examples")

	// Example 1: Basic email with SMTP
	fmt.Println("\n1. Basic Email with SMTP:")
	exampleBasicSMTP()

	// Example 2: Email with attachments
	fmt.Println("\n2. Email with Attachments:")
	exampleWithAttachments()

	// Example 3: Using the memory mailer for testing
	fmt.Println("\n3. Using Memory Mailer:")
	exampleMemoryMailer()

	// Example 4: Using mail queue
	fmt.Println("\n4. Using Mail Queue:")
	exampleMailQueue()

	// Example 5: File mailer for development
	fmt.Println("\n5. Using File Mailer:")
	exampleFileMailer()
}

func exampleBasicSMTP() {
	// For this example, we'll show the code but not execute it to avoid sending real emails
	fmt.Println("Creating basic email...")

	email := mail.NewMail().
		From("sender@example.com").
		To("recipient@example.com").
		Subject("Hello from Sapasora").
		Text("This is a simple plain text message.").
		HTML("<h1>This is an HTML message</h1>").
		Build()

	fmt.Printf("Email created: %s\n", email.Subject)

	// In a real application, you would create an SMTP mailer like this:
	/*
		mailer := mail.NewSMTPMailer(mail.SMTPConfig{
			Host:       "smtp.example.com",
			Port:       587,
			Username:   "your-username",
			Password:   "your-password",
			From:       "sender@example.com",
			Encryption: "tls",
		})

		err := mailer.Deliver(email)
		if err != nil {
			log.Printf("Error sending email: %v", err)
		}
	*/
}

func exampleWithAttachments() {
	fmt.Println("Creating email with attachments...")

	// Create email with attachment
	email := mail.NewMail().
		From("sender@example.com").
		To("recipient@example.com").
		Subject("Email with Attachment").
		Text("Please find the attached file.").
		HTML("<p>Please find the attached file.</p>").
		AttachBytes("report.txt", "text/plain", []byte("This is a sample report")).
		Build()

	fmt.Printf("Email created with %d attachment(s)\n", len(email.Attachments))
}

func exampleMemoryMailer() {
	fmt.Println("Using memory mailer for testing...")

	// Create memory mailer
	memoryMailer := mail.NewMemoryMailer()

	// Create and send an email
	email := mail.NewMail().
		From("test@example.com").
		To("user@example.com").
		Subject("Test Email").
		Text("This is a test email for memory mailer").
		Build()

	err := memoryMailer.Deliver(email)
	if err != nil {
		log.Printf("Error delivering email: %v", err)
		return
	}

	// Check how many emails are stored
	emailCount := memoryMailer.Len()
	fmt.Printf("Emails in memory: %d\n", emailCount)

	// Get the last email sent
	lastEmail, err := memoryMailer.GetLastEmail()
	if err != nil {
		log.Printf("Error getting last email: %v", err)
		return
	}

	fmt.Printf("Last email subject: %s\n", lastEmail.Subject)

	// Clear emails
	memoryMailer.Clear()
	fmt.Printf("After clearing, emails in memory: %d\n", memoryMailer.Len())
}

func exampleMailQueue() {
	fmt.Println("Using mail queue...")

	// Create a mock mailer (in real usage, this would be SMTP, file, etc.)
	mockMailer := mail.NewMemoryMailer()

	// Create queue configuration
	queueConfig := mail.MailQueueConfig{
		Workers:    2,
		BufferSize: 10,
		MaxRetries: 3,
		Timeout:    10 * time.Second,
	}

	// Create mail service with queue
	mailService := mail.NewMailService(mockMailer, queueConfig)

	// Create and queue an email
	email := mail.NewMail().
		From("queue@example.com").
		To("recipient@example.com").
		Subject("Queued Email").
		Text("This email is being sent through the queue").
		Build()

	err := mailService.Queue(email)
	if err != nil {
		log.Printf("Error queuing email: %v", err)
		return
	}

	// Get queue statistics (only available on MailServiceWithQueue)
	if queueService, ok := mailService.(*mail.MailServiceWithQueue); ok {
		stats := queueService.GetStats()
		fmt.Printf("Queue stats - Processed: %d, Failed: %d, Size: %d\n",
			stats.TotalProcessed, stats.TotalFailed, stats.QueueSize)

		// Simulate some processing time
		time.Sleep(100 * time.Millisecond)

		// Get updated stats
		stats = queueService.GetStats()
		fmt.Printf("Updated queue stats - Processed: %d, Failed: %d, Size: %d\n",
			stats.TotalProcessed, stats.TotalFailed, stats.QueueSize)
	} else {
		fmt.Println("Mail service is not a queue service")
	}
}

func exampleFileMailer() {
	fmt.Println("Using file mailer...")

	// Create file mailer (emails will be saved to ./storage/mail)
	fileMailer := mail.NewFileMailer(mail.FileMailerConfig{
		Path: "./storage/mail",
	})

	// Create and send an email
	email := mail.NewMail().
		From("file@example.com").
		To("recipient@example.com").
		Subject("File Mailer Email").
		Text("This email will be saved to a file").
		HTML("<p>This email will be saved to a file</p>").
		Build()

	err := fileMailer.Deliver(email)
	if err != nil {
		log.Printf("Error delivering email to file: %v", err)
		return
	}

	fmt.Println("Email saved to file successfully")
}
