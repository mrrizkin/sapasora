package channel

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ChannelType identifies the kind of customer-facing channel. It describes the
// channel semantics, not the implementation used to connect to it.
type ChannelType string

const (
	ChannelTypeUnknown  ChannelType = "unknown"
	ChannelTypeWhatsApp ChannelType = "whatsapp"
	ChannelTypeWhatsapp ChannelType = ChannelTypeWhatsApp // compatibility spelling
	ChannelTypeTelegram ChannelType = "telegram"
	ChannelTypeSMS      ChannelType = "sms"
	ChannelTypeEmail    ChannelType = "email"
	ChannelTypeVoice    ChannelType = "voice"
	ChannelTypeCustom   ChannelType = "custom"
)

func (t ChannelType) String() string { return string(t) }

func (t ChannelType) Valid() bool {
	switch t {
	case ChannelTypeWhatsApp, ChannelTypeTelegram, ChannelTypeSMS, ChannelTypeEmail, ChannelTypeVoice, ChannelTypeCustom:
		return true
	default:
		return false
	}
}

func (t ChannelType) MarshalJSON() ([]byte, error) { return json.Marshal(t.String()) }

func (t *ChannelType) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*t = ChannelType(value)
	if !t.Valid() {
		return fmt.Errorf("invalid channel type %q", value)
	}
	return nil
}

// Provider identifies the adapter/provider family. Provider-specific payloads
// and credentials remain outside this domain package.
type Provider string

const (
	ProviderUnknown  Provider = "unknown"
	ProviderWhatsApp Provider = "whatsapp"
	ProviderWhatsapp Provider = ProviderWhatsApp // compatibility spelling
	ProviderTelegram Provider = "telegram"
	ProviderMeta     Provider = "meta"
	ProviderTwilio   Provider = "twilio"
	ProviderSMTP     Provider = "smtp"
	ProviderCustom   Provider = "custom"
)

func (p Provider) String() string { return string(p) }

func (p Provider) Valid() bool {
	switch p {
	case ProviderWhatsApp, ProviderTelegram, ProviderMeta, ProviderTwilio, ProviderSMTP, ProviderCustom:
		return true
	default:
		return false
	}
}

func (p Provider) MarshalJSON() ([]byte, error) { return json.Marshal(p.String()) }

func (p *Provider) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*p = Provider(value)
	if !p.Valid() {
		return fmt.Errorf("invalid channel provider %q", value)
	}
	return nil
}

// ChannelStatus is an explicit domain name for the lifecycle state already
// used by the adapter boundary. Keeping it as an alias prevents the account
// model and state machine from developing divergent status vocabularies.
type ChannelStatus = ConnectionState
type ChannelState = ConnectionState

const (
	ChannelStatusPending      = StatePending
	ChannelStatusConnecting   = StateConnecting
	ChannelStatusConnected    = StateConnected
	ChannelStatusDegraded     = StateDegraded
	ChannelStatusDisconnected = StateDisconnected
	ChannelStatusExpired      = StateExpired
	ChannelStatusError        = StateError
	ChannelStatePending       = StatePending
	ChannelStateConnecting    = StateConnecting
	ChannelStateConnected     = StateConnected
	ChannelStateDegraded      = StateDegraded
	ChannelStateDisconnected  = StateDisconnected
	ChannelStateExpired       = StateExpired
	ChannelStateError         = StateError
)

// CredentialKind identifies the storage/usage class of a credential without
// containing the credential itself.
type CredentialKind string

const (
	CredentialKindUnknown     CredentialKind = "unknown"
	CredentialKindOAuth       CredentialKind = "oauth"
	CredentialKindAPIKey      CredentialKind = "api_key"
	CredentialKindAccessToken CredentialKind = "access_token"
	CredentialKindSession     CredentialKind = "session"
	CredentialKindCustom      CredentialKind = "custom"
)

func (k CredentialKind) Valid() bool {
	switch k {
	case CredentialKindOAuth, CredentialKindAPIKey, CredentialKindAccessToken, CredentialKindSession, CredentialKindCustom:
		return true
	default:
		return false
	}
}

// CredentialReference points to a secret managed by a credential vault. It is
// safe to retain on a ChannelAccount because it never carries secret bytes.
type CredentialReference struct {
	ID        string         `json:"id"`
	Kind      CredentialKind `json:"kind"`
	Version   string         `json:"version,omitempty"`
	ExpiresAt *time.Time     `json:"expires_at,omitempty"`
}

func (r CredentialReference) Valid() error {
	if strings.TrimSpace(r.ID) == "" {
		return errors.New("credential reference id is required")
	}
	if !r.Kind.Valid() {
		return fmt.Errorf("invalid credential reference kind %q", r.Kind)
	}
	return nil
}

// SessionKind identifies a provider session representation without coupling
// the channel domain to a provider SDK.
type SessionKind string

const (
	SessionKindUnknown SessionKind = "unknown"
	SessionKindBrowser SessionKind = "browser"
	SessionKindDevice  SessionKind = "device"
	SessionKindOAuth   SessionKind = "oauth"
	SessionKindRuntime SessionKind = "runtime"
	SessionKindCustom  SessionKind = "custom"
)

func (k SessionKind) Valid() bool {
	switch k {
	case SessionKindBrowser, SessionKindDevice, SessionKindOAuth, SessionKindRuntime, SessionKindCustom:
		return true
	default:
		return false
	}
}

// ProviderCredentialReference is the explicit provider-facing name for a
// credential reference without changing its opaque-reference semantics.
type ProviderCredentialReference = CredentialReference

// SessionReference points to provider session state. Session material is owned
// by the provider/session store and is intentionally not embedded here.
type SessionReference struct {
	ID         string      `json:"id"`
	Kind       SessionKind `json:"kind"`
	ExpiresAt  *time.Time  `json:"expires_at,omitempty"`
	LastSeenAt *time.Time  `json:"last_seen_at,omitempty"`
}

// ProviderSessionReference is the explicit provider-facing name for a
// session reference.
type ProviderSessionReference = SessionReference

func (r SessionReference) Valid() error {
	if strings.TrimSpace(r.ID) == "" {
		return errors.New("session reference id is required")
	}
	if !r.Kind.Valid() {
		return fmt.Errorf("invalid session reference kind %q", r.Kind)
	}
	return nil
}

// ConnectionEventType identifies an append-only lifecycle observation.
type ConnectionEventType string

const (
	ConnectionEventUnknown      ConnectionEventType = "unknown"
	ConnectionEventCreated      ConnectionEventType = "created"
	ConnectionEventConnecting   ConnectionEventType = "connecting"
	ConnectionEventConnected    ConnectionEventType = "connected"
	ConnectionEventDegraded     ConnectionEventType = "degraded"
	ConnectionEventDisconnected ConnectionEventType = "disconnected"
	ConnectionEventExpired      ConnectionEventType = "expired"
	ConnectionEventError        ConnectionEventType = "error"
)

func (t ConnectionEventType) Valid() bool {
	switch t {
	case ConnectionEventCreated, ConnectionEventConnecting, ConnectionEventConnected, ConnectionEventDegraded, ConnectionEventDisconnected, ConnectionEventExpired, ConnectionEventError:
		return true
	default:
		return false
	}
}

// ChannelConnectionEvent is a descriptive alias for the connection event
// value model.
type ChannelConnectionEvent = ConnectionEvent

// RatePolicy limits provider operations within a rolling/fixed window. Zero
// limits mean "not configured"; negative values are invalid.
type RatePolicy struct {
	RequestsPerWindow int           `json:"requests_per_window"`
	Window            time.Duration `json:"window"`
	Burst             int           `json:"burst,omitempty"`
}

func (p RatePolicy) Valid() error {
	if p.RequestsPerWindow < 0 || p.Burst < 0 {
		return errors.New("rate policy limits cannot be negative")
	}
	if p.RequestsPerWindow > 0 && p.Window <= 0 {
		return errors.New("rate policy window is required when a limit is configured")
	}
	return nil
}

// QuotaPolicy contains longer-lived usage ceilings. Zero means unlimited/not
// configured and is distinct from a provider's actual usage counters.
type QuotaPolicy struct {
	DailyLimit   int64 `json:"daily_limit,omitempty"`
	MonthlyLimit int64 `json:"monthly_limit,omitempty"`
}

func (p QuotaPolicy) Valid() error {
	if p.DailyLimit < 0 || p.MonthlyLimit < 0 {
		return errors.New("quota limits cannot be negative")
	}
	return nil
}

// QuotaRatePolicy groups channel quota and rate controls in one value object.
type QuotaRatePolicy struct {
	Quota QuotaPolicy `json:"quota"`
	Rate  RatePolicy  `json:"rate"`
}

func (p QuotaRatePolicy) Valid() error {
	if err := p.Quota.Valid(); err != nil {
		return err
	}
	return p.Rate.Valid()
}

// ChannelQuotaRatePolicy is a descriptive alias for the grouped policy.
type ChannelQuotaRatePolicy = QuotaRatePolicy

// ChannelQuotaPolicy and ChannelRatePolicy are descriptive aliases for code
// that wants the channel context in a field/type name.
type ChannelQuotaPolicy = QuotaPolicy
type ChannelRatePolicy = RatePolicy

// ChannelAccount is the provider-neutral tenant-owned channel connection.
// ID is internal-only; PublicID is the only identifier intended for external
// references. TenantID must be included in every repository lookup.
type ChannelAccount struct {
	ID       uint64 `json:"-"`
	PublicID string `json:"id"`
	TenantID string `json:"-"`

	Name       string        `json:"name"`
	Type       ChannelType   `json:"type"`
	Provider   Provider      `json:"provider"`
	ExternalID string        `json:"external_id,omitempty"`
	Status     ChannelStatus `json:"status"`

	Capabilities CapabilitySet        `json:"capabilities"`
	Credential   *CredentialReference `json:"credential,omitempty"`
	Session      *SessionReference    `json:"session,omitempty"`
	Policy       QuotaRatePolicy      `json:"policy"`

	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"-"`
	LastConnectedAt *time.Time `json:"last_connected_at,omitempty"`
	LastEventAt     *time.Time `json:"last_event_at,omitempty"`
	LastError       string     `json:"last_error,omitempty"`
}

// Valid checks invariants that are independent of persistence.
func (a ChannelAccount) Valid() error {
	if strings.TrimSpace(a.TenantID) == "" {
		return fmt.Errorf("%w: tenant id is required", ErrInvalidChannelAccount)
	}
	if strings.TrimSpace(a.PublicID) == "" {
		return fmt.Errorf("%w: public id is required", ErrInvalidChannelAccount)
	}
	if strings.TrimSpace(a.Name) == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidChannelAccount)
	}
	if !a.Type.Valid() {
		return fmt.Errorf("%w: invalid channel type %q", ErrInvalidChannelAccount, a.Type)
	}
	if !a.Provider.Valid() {
		return fmt.Errorf("%w: invalid provider %q", ErrInvalidChannelAccount, a.Provider)
	}
	if a.Status == StateUnknown || !isLifecycleState(a.Status) {
		return fmt.Errorf("%w: invalid channel status %q", ErrInvalidChannelAccount, a.Status)
	}
	if a.Credential != nil {
		if err := a.Credential.Valid(); err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidChannelAccount, err)
		}
	}
	if a.Session != nil {
		if err := a.Session.Valid(); err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidChannelAccount, err)
		}
	}
	if err := a.Policy.Valid(); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidChannelAccount, err)
	}
	return nil
}

// ConnectionEvent is an append-only, tenant-scoped observation of a channel
// lifecycle operation. Error text and metadata must be redacted by callers.
type ConnectionEvent struct {
	ID                     uint64 `json:"-"`
	PublicID               string `json:"id"`
	TenantID               string `json:"-"`
	ChannelAccountPublicID string `json:"channel_id"`

	Type       ConnectionEventType `json:"type"`
	FromState  ConnectionState     `json:"from_state,omitempty"`
	ToState    ConnectionState     `json:"to_state"`
	Reason     string              `json:"reason,omitempty"`
	ErrorCode  string              `json:"error_code,omitempty"`
	OccurredAt time.Time           `json:"occurred_at"`
	Metadata   map[string]string   `json:"metadata,omitempty"`
}

func (e ConnectionEvent) Valid() error {
	if strings.TrimSpace(e.TenantID) == "" {
		return fmt.Errorf("%w: tenant id is required", ErrInvalidConnectionEvent)
	}
	if strings.TrimSpace(e.ChannelAccountPublicID) == "" {
		return fmt.Errorf("%w: channel public id is required", ErrInvalidConnectionEvent)
	}
	if !e.Type.Valid() {
		return fmt.Errorf("%w: invalid event type %q", ErrInvalidConnectionEvent, e.Type)
	}
	if !isLifecycleState(e.ToState) {
		return fmt.Errorf("%w: invalid destination state %q", ErrInvalidConnectionEvent, e.ToState)
	}
	if e.FromState != "" && e.FromState != StateUnknown && !isLifecycleState(e.FromState) {
		return fmt.Errorf("%w: invalid source state %q", ErrInvalidConnectionEvent, e.FromState)
	}
	return nil
}

func cloneTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copyValue := *value
	return &copyValue
}

func cloneCredential(value *CredentialReference) *CredentialReference {
	if value == nil {
		return nil
	}
	copyValue := *value
	copyValue.ExpiresAt = cloneTime(value.ExpiresAt)
	return &copyValue
}

func cloneSession(value *SessionReference) *SessionReference {
	if value == nil {
		return nil
	}
	copyValue := *value
	copyValue.ExpiresAt = cloneTime(value.ExpiresAt)
	copyValue.LastSeenAt = cloneTime(value.LastSeenAt)
	return &copyValue
}

func cloneChannelAccount(value *ChannelAccount) *ChannelAccount {
	if value == nil {
		return nil
	}
	copyValue := *value
	copyValue.Capabilities = CapabilitySet{Provider: value.Capabilities.Provider, Items: value.Capabilities.List()}
	copyValue.Credential = cloneCredential(value.Credential)
	copyValue.Session = cloneSession(value.Session)
	copyValue.DeletedAt = cloneTime(value.DeletedAt)
	copyValue.LastConnectedAt = cloneTime(value.LastConnectedAt)
	copyValue.LastEventAt = cloneTime(value.LastEventAt)
	return &copyValue
}

func cloneConnectionEvent(value *ConnectionEvent) *ConnectionEvent {
	if value == nil {
		return nil
	}
	copyValue := *value
	copyValue.Metadata = cloneStringMap(value.Metadata)
	return &copyValue
}
