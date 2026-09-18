package idempotency

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

const (
	// MaxKeyLength is the maximum number of bytes accepted for an
	// Idempotency-Key. Keys are deliberately limited before they are scoped or
	// stored.
	MaxKeyLength = 255

	// DefaultTTL is used when a request does not specify an operation TTL.
	DefaultTTL = 24 * time.Hour

	// DefaultMaxTTL prevents an operation from occupying the idempotency store
	// indefinitely. A durable adapter may apply a stricter policy.
	DefaultMaxTTL = 7 * 24 * time.Hour

	// DefaultMaxPayloadBytes bounds the input used to calculate a fingerprint.
	DefaultMaxPayloadBytes = 1 << 20

	// DefaultMaxResultBytes bounds the original result retained for replay.
	DefaultMaxResultBytes = 1 << 20
)

var (
	ErrInvalidKey       = errors.New("idempotency key is invalid")
	ErrInvalidScope     = errors.New("idempotency scope is invalid")
	ErrInvalidTTL       = errors.New("idempotency ttl is invalid")
	ErrPayloadTooLarge  = errors.New("idempotency payload is too large")
	ErrResultTooLarge   = errors.New("idempotency result is too large")
	ErrConflict         = errors.New("idempotency key was used with a different request")
	ErrInProgress       = errors.New("idempotency request is already in progress")
	ErrNotOwner         = errors.New("idempotency claim is not owned")
	ErrExpired          = errors.New("idempotency claim has expired")
	ErrAlreadyCompleted = errors.New("idempotency claim is already completed")
)

// Scope identifies the caller and operation for an idempotency key. The same
// key may safely be reused in different scopes.
type Scope struct {
	WorkspaceID string
	ActorID     string
	Endpoint    string
}

// Request is the provider-neutral input to a claim operation. Payload is
// fingerprinted and is never retained by MemoryStore.
type Request struct {
	Key     string
	Scope   Scope
	Payload []byte
	TTL     time.Duration
}

// Result is the original operation result retained for duplicate replay. It
// is intentionally transport-neutral; HTTP middleware can map it to a
// response later.
type Result struct {
	Status  int
	Headers map[string]string
	Body    []byte
}

// ClaimStatus describes what a claim caller should do next.
type ClaimStatus uint8

const (
	ClaimAcquired ClaimStatus = iota + 1
	ClaimReplay
	ClaimInProgress
)

// Claim is an ownership handle returned by Claim. Only an acquired claim can
// be completed or released. The ownership token is private so callers cannot
// manufacture a valid completion handle.
type Claim struct {
	Status       ClaimStatus
	StoredResult *Result

	storeKey    string
	token       uint64
	fingerprint string
}

// Replay returns a defensive copy of the original result when this claim is a
// duplicate of a completed request.
func (c Claim) Replay() (Result, bool) {
	if c.Status != ClaimReplay || c.StoredResult == nil {
		return Result{}, false
	}
	return cloneResult(*c.StoredResult), true
}

// IsOwner reports whether this claim acquired execution ownership.
func (c Claim) IsOwner() bool { return c.Status == ClaimAcquired && c.token != 0 }

// Repository is the small boundary an eventual distributed or durable
// adapter can implement. MemoryStore is the only implementation in this
// foundation; adapter integration is deliberately deferred.
type Repository interface {
	Claim(Request) (Claim, error)
	Complete(Claim, Result) error
	Release(Claim) error
}

// Option configures MemoryStore.
type Option func(*MemoryStore)

// WithClock supplies a clock for deterministic expiry tests.
func WithClock(now func() time.Time) Option {
	return func(s *MemoryStore) {
		if now != nil {
			s.now = now
		}
	}
}

// WithDefaultTTL changes the TTL used when Request.TTL is zero.
func WithDefaultTTL(ttl time.Duration) Option {
	return func(s *MemoryStore) { s.defaultTTL = ttl }
}

// WithMaxTTL limits operation TTLs accepted by the store.
func WithMaxTTL(ttl time.Duration) Option {
	return func(s *MemoryStore) { s.maxTTL = ttl }
}

// WithMaxPayloadBytes limits bytes accepted for request fingerprinting.
func WithMaxPayloadBytes(limit int) Option {
	return func(s *MemoryStore) { s.maxPayloadBytes = limit }
}

// WithMaxResultBytes limits bytes retained for original-result replay.
func WithMaxResultBytes(limit int) Option {
	return func(s *MemoryStore) { s.maxResultBytes = limit }
}

type record struct {
	fingerprint string
	expiresAt   time.Time
	token       uint64
	result      *Result
}

// MemoryStore provides atomic process-local claims with TTL expiry. It is not
// a cross-process idempotency guarantee; use a future adapter at a deployment
// boundary when that guarantee is required.
type MemoryStore struct {
	mu sync.Mutex

	records map[string]record
	nextID  uint64
	now     func() time.Time

	defaultTTL      time.Duration
	maxTTL          time.Duration
	maxPayloadBytes int
	maxResultBytes  int
}

// NewMemoryStore constructs the process-local foundation store.
func NewMemoryStore(options ...Option) *MemoryStore {
	store := &MemoryStore{
		records:         make(map[string]record),
		now:             time.Now,
		defaultTTL:      DefaultTTL,
		maxTTL:          DefaultMaxTTL,
		maxPayloadBytes: DefaultMaxPayloadBytes,
		maxResultBytes:  DefaultMaxResultBytes,
	}
	for _, option := range options {
		if option != nil {
			option(store)
		}
	}
	return store
}

// New is a short alias for NewMemoryStore.
func New(options ...Option) *MemoryStore { return NewMemoryStore(options...) }

// ParseKey validates and bounds an Idempotency-Key without including the key
// in any returned error. Keys are opaque visible-ASCII values; whitespace and
// control characters are rejected to avoid header/parser ambiguity.
func ParseKey(raw string) (string, error) {
	if len(raw) == 0 || len(raw) > MaxKeyLength {
		return "", ErrInvalidKey
	}
	for i := 0; i < len(raw); i++ {
		if raw[i] < 0x21 || raw[i] > 0x7e {
			return "", ErrInvalidKey
		}
	}
	return raw, nil
}

// Fingerprint returns a stable SHA-256 fingerprint. Valid JSON is normalized
// before hashing, so insignificant whitespace and object key ordering do not
// create a false conflict. Non-JSON payloads are hashed as supplied.
func Fingerprint(payload []byte) (string, error) {
	return fingerprint(payload, DefaultMaxPayloadBytes)
}

func fingerprint(payload []byte, maxBytes int) (string, error) {
	if maxBytes <= 0 || len(payload) > maxBytes {
		return "", ErrPayloadTooLarge
	}

	canonical := payload
	if json.Valid(payload) {
		decoder := json.NewDecoder(bytes.NewReader(payload))
		decoder.UseNumber()
		var value any
		if err := decoder.Decode(&value); err == nil {
			if normalized, err := json.Marshal(value); err == nil {
				canonical = normalized
			}
		}
	}

	digest := sha256.Sum256(canonical)
	return hex.EncodeToString(digest[:]), nil
}

// Claim atomically reserves a scoped key. A first caller receives
// ClaimAcquired. A concurrent caller with the same fingerprint receives
// ClaimInProgress and ErrInProgress. After completion, the same request gets
// ClaimReplay and the original result; a different fingerprint gets ErrConflict.
func (s *MemoryStore) Claim(request Request) (Claim, error) {
	key, fp, ttl, err := s.normalizeRequest(request)
	if err != nil {
		return Claim{}, err
	}

	now := s.now()
	s.mu.Lock()
	defer s.mu.Unlock()

	if existing, ok := s.records[key]; ok {
		if !now.Before(existing.expiresAt) {
			delete(s.records, key)
		} else {
			if existing.fingerprint != fp {
				return Claim{}, ErrConflict
			}
			if existing.result != nil {
				result := cloneResult(*existing.result)
				return Claim{Status: ClaimReplay, StoredResult: &result, storeKey: key, fingerprint: fp}, nil
			}
			return Claim{Status: ClaimInProgress, storeKey: key, fingerprint: fp}, ErrInProgress
		}
	}

	s.nextID++
	s.records[key] = record{
		fingerprint: fp,
		expiresAt:   now.Add(ttl),
		token:       s.nextID,
	}
	return Claim{Status: ClaimAcquired, storeKey: key, token: s.nextID, fingerprint: fp}, nil
}

// Complete stores the operation's original result for replay. Completion is
// idempotent only through replay claims; an ownership handle cannot be reused.
func (s *MemoryStore) Complete(claim Claim, result Result) error {
	if !claim.IsOwner() {
		return ErrNotOwner
	}
	if err := validateResult(result, s.maxResultBytes); err != nil {
		return err
	}

	now := s.now()
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.records[claim.storeKey]
	if !ok || !now.Before(existing.expiresAt) {
		if ok {
			delete(s.records, claim.storeKey)
		}
		return ErrExpired
	}
	if existing.token != claim.token || existing.fingerprint != claim.fingerprint {
		return ErrNotOwner
	}
	if existing.result != nil {
		return ErrAlreadyCompleted
	}

	stored := cloneResult(result)
	existing.result = &stored
	s.records[claim.storeKey] = existing
	return nil
}

// Release gives up a pending claim after an operation fails before producing
// an original result. Completed records remain replayable until their TTL.
func (s *MemoryStore) Release(claim Claim) error {
	if !claim.IsOwner() {
		return ErrNotOwner
	}

	now := s.now()
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.records[claim.storeKey]
	if !ok || !now.Before(existing.expiresAt) {
		if ok {
			delete(s.records, claim.storeKey)
		}
		return ErrExpired
	}
	if existing.token != claim.token || existing.fingerprint != claim.fingerprint {
		return ErrNotOwner
	}
	if existing.result != nil {
		return ErrAlreadyCompleted
	}
	delete(s.records, claim.storeKey)
	return nil
}

func (s *MemoryStore) normalizeRequest(request Request) (string, string, time.Duration, error) {
	key, err := ParseKey(request.Key)
	if err != nil {
		return "", "", 0, err
	}
	if err := request.Scope.validate(); err != nil {
		return "", "", 0, err
	}

	ttl := request.TTL
	if ttl == 0 {
		ttl = s.defaultTTL
	}
	if ttl <= 0 || s.maxTTL <= 0 || ttl > s.maxTTL {
		return "", "", 0, ErrInvalidTTL
	}

	fp, err := fingerprint(request.Payload, s.maxPayloadBytes)
	if err != nil {
		return "", "", 0, err
	}
	return scopedKey(request.Scope, key), fp, ttl, nil
}

func (s Scope) validate() error {
	if !validReference(s.WorkspaceID, MaxKeyLength) ||
		!validReference(s.ActorID, MaxKeyLength) ||
		!validReference(s.Endpoint, MaxKeyLength*2) {
		return ErrInvalidScope
	}
	return nil
}

func validReference(value string, maxBytes int) bool {
	if len(value) == 0 || len(value) > maxBytes {
		return false
	}
	for i := 0; i < len(value); i++ {
		if value[i] < 0x21 || value[i] > 0x7e {
			return false
		}
	}
	return true
}

func scopedKey(scope Scope, key string) string {
	hash := sha256.New()
	for _, value := range []string{scope.WorkspaceID, scope.ActorID, scope.Endpoint, key} {
		var length [8]byte
		binary.BigEndian.PutUint64(length[:], uint64(len(value)))
		_, _ = hash.Write(length[:])
		_, _ = hash.Write([]byte(value))
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func validateResult(result Result, maxBytes int) error {
	if maxBytes <= 0 || len(result.Body) > maxBytes {
		return ErrResultTooLarge
	}
	return nil
}

func cloneResult(result Result) Result {
	clone := Result{Status: result.Status}
	if result.Body != nil {
		clone.Body = append([]byte(nil), result.Body...)
	}
	if result.Headers != nil {
		clone.Headers = make(map[string]string, len(result.Headers))
		for key, value := range result.Headers {
			clone.Headers[key] = value
		}
	}
	return clone
}
