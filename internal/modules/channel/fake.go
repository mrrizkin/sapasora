package channel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
)

const fakeProvider = "fake"

// FakeAdapter is a deterministic, in-memory implementation of Adapter for
// contract tests. It never performs network I/O, uses fixed message IDs, and
// returns queued events in insertion order.
type FakeAdapter struct {
	mu sync.Mutex

	// Err is returned by operations that do not have an operation-specific
	// error configured. It is intentionally exported for concise test setup.
	Err error

	provider          string
	capabilities      CapabilitySet
	capabilitiesSet   bool
	state             ConnectionState
	contacts          map[string]Contact
	users             map[string]User
	events            []RawEvent
	operationErrors   map[string]error
	calls             []string
	sent              []SendResult
	nextMessageNumber int
}

var _ Adapter = (*FakeAdapter)(nil)

// NewFakeAdapter constructs a ready-to-use deterministic fake.
func NewFakeAdapter() *FakeAdapter {
	return &FakeAdapter{
		provider:          fakeProvider,
		state:             StateDisconnected,
		contacts:          make(map[string]Contact),
		users:             make(map[string]User),
		operationErrors:   make(map[string]error),
		nextMessageNumber: 1,
	}
}

func (f *FakeAdapter) ensureInitializedLocked() {
	if f.provider == "" {
		f.provider = fakeProvider
	}
	if f.state == StateUnknown {
		f.state = StateDisconnected
	}
	if f.contacts == nil {
		f.contacts = make(map[string]Contact)
	}
	if f.users == nil {
		f.users = make(map[string]User)
	}
	if f.operationErrors == nil {
		f.operationErrors = make(map[string]error)
	}
	if f.nextMessageNumber == 0 {
		f.nextMessageNumber = 1
	}
}

// Provider returns the stable fake provider identifier.
func (f *FakeAdapter) Provider() string {
	if f == nil {
		return ""
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.ensureInitializedLocked()
	return f.provider
}

// SetError configures an operation-specific error. An empty operation changes
// the default Err used by all operations.
func (f *FakeAdapter) SetError(operation string, err error) {
	if f == nil {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.ensureInitializedLocked()
	if operation == "" {
		f.Err = err
		return
	}
	f.operationErrors[operation] = err
}

// SetCapabilities replaces the advertised capabilities.
func (f *FakeAdapter) SetCapabilities(capabilities CapabilitySet) {
	if f == nil {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.capabilities = CapabilitySet{Provider: capabilities.Provider, Items: capabilities.List()}
	f.capabilitiesSet = true
}

// SetContact makes a contact available to GetContact.
func (f *FakeAdapter) SetContact(contact Contact) {
	if f == nil {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.ensureInitializedLocked()
	f.contacts[addressKey(contact.Address)] = contact
}

// SetUser makes a user available to GetUser.
func (f *FakeAdapter) SetUser(user User) {
	if f == nil {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.ensureInitializedLocked()
	f.users[addressKey(user.Address)] = user
}

// QueueEvent appends an event to the deterministic receive queue.
func (f *FakeAdapter) QueueEvent(event RawEvent) {
	if f == nil {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.ensureInitializedLocked()
	f.events = append(f.events, cloneRawEvent(event))
}

// Calls returns a copy of operation names in invocation order.
func (f *FakeAdapter) Calls() []string {
	if f == nil {
		return nil
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string(nil), f.calls...)
}

// LastCall returns the last operation name, or an empty string.
func (f *FakeAdapter) LastCall() string {
	calls := f.Calls()
	if len(calls) == 0 {
		return ""
	}
	return calls[len(calls)-1]
}

// Sent returns a copy of successful send results.
func (f *FakeAdapter) Sent() []SendResult {
	if f == nil {
		return nil
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]SendResult(nil), f.sent...)
}

func (f *FakeAdapter) record(ctx context.Context, operation string) error {
	if f == nil {
		return errors.New("fake channel adapter is nil")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.ensureInitializedLocked()
	f.calls = append(f.calls, operation)
	if err, ok := f.operationErrors[operation]; ok {
		return err
	}
	return f.Err
}

func (f *FakeAdapter) Connect(ctx context.Context, _ ConnectRequest) (ConnectionState, error) {
	if err := f.record(ctx, "connect"); err != nil {
		return StateError, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.state = StateConnected
	return f.state, nil
}

func (f *FakeAdapter) Disconnect(ctx context.Context, _ DisconnectRequest) error {
	if err := f.record(ctx, "disconnect"); err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.state = StateDisconnected
	return nil
}

func (f *FakeAdapter) Health(ctx context.Context) (HealthStatus, error) {
	if err := f.record(ctx, "health"); err != nil {
		return HealthStatus{State: StateError}, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.ensureInitializedLocked()
	return HealthStatus{State: f.state, Healthy: f.state == StateConnected}, nil
}

func (f *FakeAdapter) Capabilities() CapabilitySet {
	if f == nil {
		return CapabilitySet{}
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.ensureInitializedLocked()
	if f.capabilitiesSet {
		return CapabilitySet{Provider: f.capabilities.Provider, Items: f.capabilities.List()}
	}
	return CapabilitySet{
		Provider: f.provider,
		Items: []Capability{
			CapabilitySendText,
			CapabilityGetContact,
			CapabilityGetUser,
			CapabilityReceiveEvent,
			CapabilityValidateAddress,
		},
	}
}

func (f *FakeAdapter) Send(ctx context.Context, _ SendRequest) (SendResult, error) {
	if err := f.record(ctx, "send"); err != nil {
		return SendResult{Status: DeliveryFailed}, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.ensureInitializedLocked()
	result := SendResult{
		MessageID: fmt.Sprintf("fake-message-%d", f.nextMessageNumber),
		Status:    DeliveryAccepted,
	}
	f.nextMessageNumber++
	f.sent = append(f.sent, result)
	return result, nil
}

func (f *FakeAdapter) GetContact(ctx context.Context, request GetContactRequest) (Contact, error) {
	if err := f.record(ctx, "get_contact"); err != nil {
		return Contact{}, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.ensureInitializedLocked()
	if contact, ok := f.contacts[addressKey(request.Address)]; ok {
		return contact, nil
	}
	return Contact{Address: request.Address}, nil
}

func (f *FakeAdapter) GetUser(ctx context.Context, request GetUserRequest) (User, error) {
	if err := f.record(ctx, "get_user"); err != nil {
		return User{}, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.ensureInitializedLocked()
	if user, ok := f.users[addressKey(request.Address)]; ok {
		return user, nil
	}
	return User{Address: request.Address}, nil
}

func (f *FakeAdapter) ReceiveEvent(ctx context.Context) (RawEvent, error) {
	if err := f.record(ctx, "receive_event"); err != nil {
		return RawEvent{}, err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.ensureInitializedLocked()
	if len(f.events) == 0 {
		return RawEvent{}, ErrNoEvent
	}
	event := cloneRawEvent(f.events[0])
	f.events = f.events[1:]
	return event, nil
}

func (f *FakeAdapter) ValidateAddress(ctx context.Context, address Address) (Address, error) {
	if err := f.record(ctx, "validate_address"); err != nil {
		return Address{}, err
	}
	address.Value = strings.TrimSpace(address.Value)
	if address.Value == "" {
		return Address{}, fmt.Errorf("%w: value is empty", ErrInvalidAddress)
	}
	if address.Kind == "" || address.Kind == AddressUnknown {
		address.Kind = AddressExternal
	}
	return address, nil
}

func (f *FakeAdapter) NormalizeError(err error) *ProviderError {
	if f != nil {
		_ = f.record(context.Background(), "normalize_error")
	}
	return normalizeFakeError(err)
}

func (f *FakeAdapter) NormalizeInboundEvent(ctx context.Context, raw RawEvent) (InboundEvent, error) {
	if err := f.record(ctx, "normalize_inbound_event"); err != nil {
		return InboundEvent{}, err
	}
	event := InboundEvent{}
	if len(raw.Payload) > 0 {
		if err := json.Unmarshal(raw.Payload, &event); err != nil {
			return InboundEvent{}, fmt.Errorf("decode normalized inbound event: %w", err)
		}
	}
	if event.ID == "" {
		event.ID = raw.ID
	}
	if event.Kind == "" {
		event.Kind = InboundEventKind(raw.Type)
	}
	if event.OccurredAt.IsZero() {
		event.OccurredAt = raw.ReceivedAt
	}
	if event.Metadata == nil {
		event.Metadata = map[string]string{}
	}
	return event, nil
}

func addressKey(address Address) string {
	return string(address.Kind) + ":" + address.Value
}

func cloneRawEvent(event RawEvent) RawEvent {
	return RawEvent{
		ID:         event.ID,
		Type:       event.Type,
		Payload:    append([]byte(nil), event.Payload...),
		Headers:    cloneStringMap(event.Headers),
		ReceivedAt: event.ReceivedAt,
	}
}

func cloneStringMap(values map[string]string) map[string]string {
	if values == nil {
		return nil
	}
	clone := make(map[string]string, len(values))
	for key, value := range values {
		clone[key] = value
	}
	return clone
}
