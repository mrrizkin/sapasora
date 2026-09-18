package channel

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	// ErrNoEvent indicates that a non-blocking event source has no event ready.
	ErrNoEvent = errors.New("no channel event available")
	// ErrInvalidAddress indicates that an address cannot be normalized.
	ErrInvalidAddress = errors.New("invalid channel address")
	// ErrUnsupported indicates that a capability is not supported by an adapter.
	ErrUnsupported = errors.New("channel capability unsupported")
	// ErrNotFound indicates that a requested contact or user does not exist.
	ErrNotFound = errors.New("channel resource not found")
)

// ErrorCategory is the stable category used by callers for retry, auth, and
// user-action decisions. Provider error strings must not be used for control
// flow.
type ErrorCategory string

const (
	ErrorCategoryUnknown        ErrorCategory = "unknown"
	ErrorCategoryAuthentication ErrorCategory = "authentication"
	ErrorCategoryAuthorization  ErrorCategory = "authorization"
	ErrorCategoryInvalidAddress ErrorCategory = "invalid_address"
	ErrorCategoryInvalidRequest ErrorCategory = "invalid_request"
	ErrorCategoryUnsupported    ErrorCategory = "unsupported"
	ErrorCategoryRateLimited    ErrorCategory = "rate_limited"
	ErrorCategoryTemporary      ErrorCategory = "temporary"
	ErrorCategoryUnavailable    ErrorCategory = "unavailable"
	ErrorCategoryTimeout        ErrorCategory = "timeout"
	ErrorCategoryNotFound       ErrorCategory = "not_found"
	ErrorCategoryConflict       ErrorCategory = "conflict"
	ErrorCategoryInternal       ErrorCategory = "internal"
)

// ProviderError is the normalized, safe representation of an adapter failure.
// Cause is retained for errors.Is/errors.As and must not be serialized or
// logged as an unredacted provider error.
type ProviderError struct {
	Provider   string
	Category   ErrorCategory
	Code       string
	Message    string
	Retryable  bool
	RetryAfter time.Duration
	Cause      error
}

func (e *ProviderError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Code != "" {
		return e.Code
	}
	if e.Category != "" {
		return string(e.Category)
	}
	return string(ErrorCategoryUnknown)
}

func (e *ProviderError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// IsCategory reports whether err is a normalized provider error in category.
func IsCategory(err error, category ErrorCategory) bool {
	var providerErr *ProviderError
	return errors.As(err, &providerErr) && providerErr.Category == category
}

// ErrorMapping maps one provider error code to a stable platform category.
type ErrorMapping struct {
	ProviderCode string
	Code         string
	Category     ErrorCategory
	Retryable    bool
	RetryAfter   time.Duration
}

// ErrorMappingCatalog is a deterministic provider-error mapping catalog.
// Lookups are exact and the fallback is used for unknown provider codes.
type ErrorMappingCatalog struct {
	mappings map[string]ErrorMapping
	fallback ErrorMapping
}

// NewErrorMappingCatalog creates a catalog. Duplicate provider codes are
// rejected so a mapping cannot silently change based on input order.
func NewErrorMappingCatalog(entries []ErrorMapping) (ErrorMappingCatalog, error) {
	catalog := ErrorMappingCatalog{
		mappings: make(map[string]ErrorMapping, len(entries)),
		fallback: ErrorMapping{
			Code:      "provider_error",
			Category:  ErrorCategoryUnknown,
			Retryable: false,
		},
	}
	for _, entry := range entries {
		if entry.ProviderCode == "" {
			return ErrorMappingCatalog{}, errors.New("provider error mapping code is required")
		}
		if _, exists := catalog.mappings[entry.ProviderCode]; exists {
			return ErrorMappingCatalog{}, fmt.Errorf("duplicate provider error mapping: %s", entry.ProviderCode)
		}
		if entry.Code == "" {
			return ErrorMappingCatalog{}, fmt.Errorf("normalized error code is required for %s", entry.ProviderCode)
		}
		if entry.Category == "" {
			return ErrorMappingCatalog{}, fmt.Errorf("error category is required for %s", entry.ProviderCode)
		}
		catalog.mappings[entry.ProviderCode] = entry
	}
	return catalog, nil
}

// WithFallback returns a catalog with a replacement unknown-code mapping.
func (c ErrorMappingCatalog) WithFallback(fallback ErrorMapping) (ErrorMappingCatalog, error) {
	if fallback.Code == "" || fallback.Category == "" {
		return ErrorMappingCatalog{}, errors.New("fallback error mapping requires code and category")
	}
	copyCatalog := ErrorMappingCatalog{
		mappings: make(map[string]ErrorMapping, len(c.mappings)),
		fallback: fallback,
	}
	for code, mapping := range c.mappings {
		copyCatalog.mappings[code] = mapping
	}
	return copyCatalog, nil
}

// Lookup returns the exact mapping or the catalog fallback and whether the
// provider code was explicitly registered.
func (c ErrorMappingCatalog) Lookup(providerCode string) (ErrorMapping, bool) {
	mapping, ok := c.mappings[providerCode]
	if ok {
		return mapping, true
	}
	if c.fallback.Code == "" {
		return ErrorMapping{Code: "provider_error", Category: ErrorCategoryUnknown}, false
	}
	return c.fallback, false
}

// Normalize creates a safe typed error from a provider code and optional cause.
func (c ErrorMappingCatalog) Normalize(provider, providerCode string, cause error) *ProviderError {
	if providerCode == "" && cause != nil {
		var normalized *ProviderError
		if errors.As(cause, &normalized) {
			return normalized
		}
	}
	mapping, _ := c.Lookup(providerCode)
	message := mapping.Code
	if cause != nil && mapping.Category == ErrorCategoryUnknown {
		message = "provider operation failed"
	}
	return &ProviderError{
		Provider:   provider,
		Category:   mapping.Category,
		Code:       mapping.Code,
		Message:    message,
		Retryable:  mapping.Retryable,
		RetryAfter: mapping.RetryAfter,
		Cause:      cause,
	}
}

// DefaultErrorMappingCatalog provides stable mappings for common adapter
// failures. Provider implementations can extend it with provider codes.
func DefaultErrorMappingCatalog() ErrorMappingCatalog {
	catalog, err := NewErrorMappingCatalog([]ErrorMapping{
		{ProviderCode: "invalid_address", Code: "invalid_address", Category: ErrorCategoryInvalidAddress},
		{ProviderCode: "unsupported", Code: "unsupported", Category: ErrorCategoryUnsupported},
		{ProviderCode: "rate_limited", Code: "rate_limited", Category: ErrorCategoryRateLimited, Retryable: true},
		{ProviderCode: "timeout", Code: "timeout", Category: ErrorCategoryTimeout, Retryable: true},
		{ProviderCode: "authentication", Code: "authentication_failed", Category: ErrorCategoryAuthentication},
	})
	if err != nil {
		panic("invalid default channel error mapping catalog: " + err.Error())
	}
	return catalog
}

func normalizeFakeError(err error) *ProviderError {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return DefaultErrorMappingCatalog().Normalize("fake", "timeout", err)
	}
	if errors.Is(err, ErrInvalidAddress) {
		return DefaultErrorMappingCatalog().Normalize("fake", "invalid_address", err)
	}
	if errors.Is(err, ErrUnsupported) {
		return DefaultErrorMappingCatalog().Normalize("fake", "unsupported", err)
	}
	return DefaultErrorMappingCatalog().Normalize("fake", "", err)
}
