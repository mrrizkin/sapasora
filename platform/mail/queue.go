// Package mail provides mail queue functionality
package mail

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MailQueue represents a queue for sending emails asynchronously
type MailQueue struct {
	mailer  Mailer
	queue   chan QueuedEmail
	workers int
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	mu      sync.RWMutex
	stats   QueueStats
}

// QueuedEmail represents an email in the queue
type QueuedEmail struct {
	Email      Email
	Timestamp  time.Time
	Attempts   int
	MaxRetries int
}

// QueueStats holds statistics about the mail queue
type QueueStats struct {
	TotalProcessed int64
	TotalFailed    int64
	QueueSize      int
}

// MailQueueConfig holds configuration for the mail queue
type MailQueueConfig struct {
	Workers    int           // Number of worker goroutines
	BufferSize int           // Size of the queue channel
	MaxRetries int           // Maximum number of retry attempts
	Timeout    time.Duration // Timeout for email delivery
}

// NewMailQueue creates a new mail queue
func NewMailQueue(mailer Mailer, config MailQueueConfig) *MailQueue {
	if config.Workers <= 0 {
		config.Workers = 1
	}
	if config.BufferSize <= 0 {
		config.BufferSize = 100
	}
	if config.MaxRetries <= 0 {
		config.MaxRetries = 3
	}
	if config.Timeout <= 0 {
		config.Timeout = 30 * time.Second
	}

	ctx, cancel := context.WithCancel(context.Background())

	queue := &MailQueue{
		mailer:  mailer,
		queue:   make(chan QueuedEmail, config.BufferSize),
		workers: config.Workers,
		ctx:     ctx,
		cancel:  cancel,
		stats:   QueueStats{QueueSize: config.BufferSize},
	}

	// Start worker goroutines
	for i := 0; i < config.Workers; i++ {
		queue.wg.Add(1)
		go queue.worker()
	}

	return queue
}

// worker processes emails from the queue
func (mq *MailQueue) worker() {
	defer mq.wg.Done()

	for {
		select {
		case email := <-mq.queue:
			// Process the email
			err := mq.processEmail(email)

			mq.mu.Lock()
			if err != nil {
				mq.stats.TotalFailed++
				// If we haven't exceeded max retries, try again
				if email.Attempts < email.MaxRetries {
					newEmail := QueuedEmail{
						Email:      email.Email,
						Timestamp:  time.Now(),
						Attempts:   email.Attempts + 1,
						MaxRetries: email.MaxRetries,
					}
					// Put it back in the queue with a short delay
					go func() {
						time.Sleep(time.Duration(email.Attempts) * time.Second)
						mq.queue <- newEmail
					}()
				}
			} else {
				mq.stats.TotalProcessed++
			}
			mq.mu.Unlock()

		case <-mq.ctx.Done():
			return
		}
	}
}

// processEmail delivers an email using the configured mailer
func (mq *MailQueue) processEmail(qe QueuedEmail) error {
	// Create a context with timeout for email delivery
	ctx, cancel := context.WithTimeout(mq.ctx, 30*time.Second)
	defer cancel()

	return mq.mailer.DeliverWithContext(ctx, qe.Email)
}

// Queue adds an email to the queue
func (mq *MailQueue) Queue(email Email) error {
	return mq.QueueWithContext(context.Background(), email)
}

// QueueWithContext adds an email to the queue with context
func (mq *MailQueue) QueueWithContext(ctx context.Context, email Email) error {
	queuedEmail := QueuedEmail{
		Email:      email,
		Timestamp:  time.Now(),
		Attempts:   0,
		MaxRetries: 3, // Default max retries
	}

	select {
	case mq.queue <- queuedEmail:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-mq.ctx.Done():
		return fmt.Errorf("mail queue is shutting down")
	}
}

// Stop gracefully shuts down the mail queue
func (mq *MailQueue) Stop() {
	mq.cancel()
	mq.wg.Wait()
	close(mq.queue)
}

// GetStats returns current queue statistics
func (mq *MailQueue) GetStats() QueueStats {
	mq.mu.RLock()
	defer mq.mu.RUnlock()

	stats := mq.stats
	stats.QueueSize = len(mq.queue)
	return stats
}

// MailServiceWithQueue implements the MailService interface with queue support
type MailServiceWithQueue struct {
	mailer Mailer
	queue  *MailQueue
}

// NewMailService creates a new mail service with queue support
func NewMailService(mailer Mailer, config MailQueueConfig) MailService {
	queue := NewMailQueue(mailer, config)

	return &MailServiceWithQueue{
		mailer: mailer,
		queue:  queue,
	}
}

// Deliver sends an email immediately via the underlying mailer
func (ms *MailServiceWithQueue) Deliver(email Email) error {
	return ms.mailer.Deliver(email)
}

// DeliverWithContext sends an email immediately via the underlying mailer with context
func (ms *MailServiceWithQueue) DeliverWithContext(ctx context.Context, email Email) error {
	return ms.mailer.DeliverWithContext(ctx, email)
}

// Queue sends an email via the queue
func (ms *MailServiceWithQueue) Queue(email Email) error {
	return ms.queue.Queue(email)
}

// QueueWithContext sends an email via the queue with context
func (ms *MailServiceWithQueue) QueueWithContext(ctx context.Context, email Email) error {
	return ms.queue.QueueWithContext(ctx, email)
}

// Stop stops the mail service queue
func (ms *MailServiceWithQueue) Stop() {
	ms.queue.Stop()
}

// GetStats returns queue statistics
func (ms *MailServiceWithQueue) GetStats() QueueStats {
	return ms.queue.GetStats()
}
