// Package mail provides a complete mail service similar to Elixir Phoenix's Swoosh
// It includes multiple mailer backends, queuing capabilities, and a fluent builder API
package mail

// Version of the mail package
const Version = "1.0.0"

// Common functions and types are available directly from this package:
//  - NewMail() - Creates a new email builder
//  - Email - The email type
//  - Attachment - For email attachments
//  - Mailer - Interface for mail delivery
//  - MailService - Interface with queuing support
//
// Available mailers:
//  - NewSMTPMailer() - For sending emails via SMTP
//  - NewFileMailer() - For saving emails to files (development)
//  - NewMemoryMailer() - For testing
//
// Available services:
//  - NewMailService() - Mail service with queuing
//  - NewMailServiceFromConfig() - Mail service from configuration
//
// Configuration:
//  - LoadConfig() - Load mail configuration
//  - NewMailerFromConfig() - Create mailer from configuration
//  - NewMailServiceFromConfig() - Create mail service from configuration
