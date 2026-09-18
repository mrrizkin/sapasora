// Package mail provides in-memory mailer for testing purposes
package mail

import (
	"context"
	"fmt"
	"sync"
)

// MemoryMailer stores emails in memory, useful for testing
type MemoryMailer struct {
	emails []Email
	mutex  sync.RWMutex
}

// NewMemoryMailer creates a new memory mailer
func NewMemoryMailer() *MemoryMailer {
	return &MemoryMailer{
		emails: make([]Email, 0),
	}
}

// Deliver stores the email in memory
func (m *MemoryMailer) Deliver(email Email) error {
	return m.DeliverWithContext(context.Background(), email)
}

// DeliverWithContext stores the email in memory with context
func (m *MemoryMailer) DeliverWithContext(ctx context.Context, email Email) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.emails = append(m.emails, email)
	return nil
}

// GetEmails returns all sent emails
func (m *MemoryMailer) GetEmails() []Email {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Return a copy to prevent external modification
	result := make([]Email, len(m.emails))
	copy(result, m.emails)
	return result
}

// GetLastEmail returns the most recent email sent
func (m *MemoryMailer) GetLastEmail() (Email, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if len(m.emails) == 0 {
		return Email{}, fmt.Errorf("no emails in memory")
	}

	return m.emails[len(m.emails)-1], nil
}

// Clear removes all emails from memory
func (m *MemoryMailer) Clear() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.emails = make([]Email, 0)
}

// Len returns the number of emails stored
func (m *MemoryMailer) Len() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.emails)
}
