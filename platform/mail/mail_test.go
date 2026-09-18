package mail

import (
	"testing"
)

func TestMailBuilder(t *testing.T) {
	// Test creating a basic email
	email := NewMail().
		From("sender@example.com").
		To("recipient@example.com").
		Subject("Test Subject").
		Text("Test text").
		HTML("<p>Test HTML</p>").
		Build()

	if email.From != "sender@example.com" {
		t.Errorf("Expected From to be 'sender@example.com', got '%s'", email.From)
	}

	if len(email.To) != 1 || email.To[0] != "recipient@example.com" {
		t.Errorf("Expected To to contain 'recipient@example.com', got %v", email.To)
	}

	if email.Subject != "Test Subject" {
		t.Errorf("Expected Subject to be 'Test Subject', got '%s'", email.Subject)
	}

	if email.Text != "Test text" {
		t.Errorf("Expected Text to be 'Test text', got '%s'", email.Text)
	}

	if email.HTML != "<p>Test HTML</p>" {
		t.Errorf("Expected HTML to be '<p>Test HTML</p>', got '%s'", email.HTML)
	}
}

func TestMailBuilderMultipleRecipients(t *testing.T) {
	// Test adding multiple recipients
	email := NewMail().
		To("recipient1@example.com", "recipient2@example.com").
		Cc("cc@example.com").
		Bcc("bcc@example.com").
		Subject("Multi-recipients").
		Text("Test").
		Build()

	if len(email.To) != 2 {
		t.Errorf("Expected 2 To recipients, got %d", len(email.To))
	}

	if email.To[0] != "recipient1@example.com" || email.To[1] != "recipient2@example.com" {
		t.Errorf(
			"Expected To recipients to be ['recipient1@example.com', 'recipient2@example.com'], got %v",
			email.To,
		)
	}

	if len(email.Cc) != 1 || email.Cc[0] != "cc@example.com" {
		t.Errorf("Expected Cc to contain 'cc@example.com', got %v", email.Cc)
	}

	if len(email.Bcc) != 1 || email.Bcc[0] != "bcc@example.com" {
		t.Errorf("Expected Bcc to contain 'bcc@example.com', got %v", email.Bcc)
	}
}

func TestMemoryMailer(t *testing.T) {
	mailer := NewMemoryMailer()

	email := NewMail().
		To("test@example.com").
		Subject("Test").
		Text("Test message").
		Build()

	err := mailer.Deliver(email)
	if err != nil {
		t.Fatalf("Error delivering email: %v", err)
	}

	if mailer.Len() != 1 {
		t.Errorf("Expected 1 email in memory, got %d", mailer.Len())
	}

	sentEmails := mailer.GetEmails()
	if len(sentEmails) != 1 {
		t.Errorf("Expected 1 sent email, got %d", len(sentEmails))
	}

	if sentEmails[0].Subject != "Test" {
		t.Errorf("Expected subject 'Test', got '%s'", sentEmails[0].Subject)
	}

	lastEmail, err := mailer.GetLastEmail()
	if err != nil {
		t.Fatalf("Error getting last email: %v", err)
	}

	if lastEmail.Subject != "Test" {
		t.Errorf("Expected last email subject 'Test', got '%s'", lastEmail.Subject)
	}
}
