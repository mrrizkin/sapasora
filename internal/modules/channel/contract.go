package channel

import (
	"context"
	"time"
)

// ConnectionState is the provider-neutral lifecycle state reported by an adapter.
type ConnectionState string

const (
	StateUnknown      ConnectionState = "unknown"
	StatePending      ConnectionState = "pending"
	StateConnecting   ConnectionState = "connecting"
	StateConnected    ConnectionState = "connected"
	StateDegraded     ConnectionState = "degraded"
	StateDisconnected ConnectionState = "disconnected"
	StateExpired      ConnectionState = "expired"
	StateError        ConnectionState = "error"
)

// AddressKind identifies the semantic kind of a channel address without
// importing a provider's address or payload type.
type AddressKind string

const (
	AddressUnknown  AddressKind = "unknown"
	AddressPhone    AddressKind = "phone"
	AddressUsername AddressKind = "username"
	AddressEmail    AddressKind = "email"
	AddressExternal AddressKind = "external"
)

// Address is a normalized provider-neutral recipient or actor address.
type Address struct {
	Kind  AddressKind
	Value string
}

// ConnectRequest contains references and metadata needed to establish a
// connection. CredentialRef is a reference, never a credential value.
type ConnectRequest struct {
	AccountID     string
	Address       Address
	CredentialRef string
	Metadata      map[string]string
}

// DisconnectRequest describes why a connection is being closed.
type DisconnectRequest struct {
	Reason string
}

// HealthStatus is a point-in-time provider health result.
type HealthStatus struct {
	State     ConnectionState
	Healthy   bool
	Detail    string
	CheckedAt time.Time
}

// Capability is a provider-neutral capability identifier.
type Capability string

const (
	CapabilitySendText        Capability = "send.text"
	CapabilitySendMedia       Capability = "send.media"
	CapabilityGetContact      Capability = "get.contact"
	CapabilityGetUser         Capability = "get.user"
	CapabilityReceiveEvent    Capability = "receive.event"
	CapabilityValidateAddress Capability = "validate.address"
)

// CapabilitySet is an immutable-by-convention snapshot of adapter features.
type CapabilitySet struct {
	Provider string
	Items    []Capability
}

// Has reports whether the capability is present.
func (c CapabilitySet) Has(capability Capability) bool {
	for _, item := range c.Items {
		if item == capability {
			return true
		}
	}
	return false
}

// List returns a copy of the advertised capabilities.
func (c CapabilitySet) List() []Capability {
	return append([]Capability(nil), c.Items...)
}

// MessageKind identifies a provider-neutral outbound or inbound message kind.
type MessageKind string

const (
	MessageText     MessageKind = "text"
	MessageImage    MessageKind = "image"
	MessageAudio    MessageKind = "audio"
	MessageVideo    MessageKind = "video"
	MessageDocument MessageKind = "document"
	MessageSticker  MessageKind = "sticker"
	MessageLocation MessageKind = "location"
	MessageContact  MessageKind = "contact"
)

// Attachment describes media by provider-neutral references and metadata.
type Attachment struct {
	URL      string
	MIMEType string
	Filename string
}

// Message is the provider-neutral message shape at the adapter boundary.
type Message struct {
	Kind        MessageKind
	Body        string
	Attachments []Attachment
	Metadata    map[string]string
}

// SendRequest is an outbound message request.
type SendRequest struct {
	To             Address
	Message        Message
	IdempotencyKey string
}

// DeliveryStatus is the initial provider-neutral send result state.
type DeliveryStatus string

const (
	DeliveryAccepted DeliveryStatus = "accepted"
	DeliverySent     DeliveryStatus = "sent"
	DeliveryFailed   DeliveryStatus = "failed"
)

// SendResult identifies a provider-accepted outbound message.
type SendResult struct {
	MessageID  string
	Status     DeliveryStatus
	AcceptedAt time.Time
}

// Contact is a normalized channel contact.
type Contact struct {
	Address     Address
	DisplayName string
	FirstName   string
	LastName    string
}

// User is a normalized channel user identity.
type User struct {
	Address     Address
	DisplayName string
	Username    string
}

// GetContactRequest requests one contact by normalized address.
type GetContactRequest struct {
	Address Address
}

// GetUserRequest requests one user by normalized address.
type GetUserRequest struct {
	Address Address
}

// RawEvent is an opaque provider event envelope. Payload remains owned by the
// provider adapter and is not exposed as a provider-specific Go type.
type RawEvent struct {
	ID         string
	Type       string
	Payload    []byte
	Headers    map[string]string
	ReceivedAt time.Time
}

// InboundEventKind identifies a normalized inbound event.
type InboundEventKind string

const (
	InboundMessageEvent InboundEventKind = "message"
	InboundStatusEvent  InboundEventKind = "status"
	InboundContactEvent InboundEventKind = "contact"
)

// InboundEvent is the provider-neutral event shape consumed by the channel
// domain.
type InboundEvent struct {
	ID         string
	Kind       InboundEventKind
	Sender     Address
	Contact    *Contact
	Message    *Message
	OccurredAt time.Time
	Metadata   map[string]string
}

// Connect establishes a provider connection.
type Connect interface {
	Connect(context.Context, ConnectRequest) (ConnectionState, error)
}

// Disconnect gracefully closes a provider connection.
type Disconnect interface {
	Disconnect(context.Context, DisconnectRequest) error
}

// Health reports provider and connection health.
type Health interface {
	Health(context.Context) (HealthStatus, error)
}

// Capabilities reports the features supported by the adapter.
type Capabilities interface {
	Capabilities() CapabilitySet
}

// Send sends a provider-neutral message.
type Send interface {
	Send(context.Context, SendRequest) (SendResult, error)
}

// GetContact retrieves a normalized contact.
type GetContact interface {
	GetContact(context.Context, GetContactRequest) (Contact, error)
}

// GetUser retrieves a normalized user identity.
type GetUser interface {
	GetUser(context.Context, GetUserRequest) (User, error)
}

// ReceiveEvent receives one opaque provider event envelope. ErrNoEvent may be
// returned when a non-blocking source currently has no event.
type ReceiveEvent interface {
	ReceiveEvent(context.Context) (RawEvent, error)
}

// ValidateAddress validates and normalizes a channel address.
type ValidateAddress interface {
	ValidateAddress(context.Context, Address) (Address, error)
}

// NormalizeError maps provider failures to the typed provider error catalog.
type NormalizeError interface {
	NormalizeError(error) *ProviderError
}

// NormalizeInboundEvent maps an opaque provider event to a domain event.
type NormalizeInboundEvent interface {
	NormalizeInboundEvent(context.Context, RawEvent) (InboundEvent, error)
}

// Adapter is a convenience aggregate for adapters that implement every
// capability. Providers may implement only the capability interfaces they
// support; existing providers are not required to implement Adapter.
type Adapter interface {
	Connect
	Disconnect
	Health
	Capabilities
	Send
	GetContact
	GetUser
	ReceiveEvent
	ValidateAddress
	NormalizeError
	NormalizeInboundEvent
}
