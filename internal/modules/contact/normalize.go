package contact

import (
	"errors"
	"strings"
	"unicode"
)

var (
	errInvalidPhone    = errors.New("invalid phone address")
	errInvalidEmail    = errors.New("invalid email address")
	errInvalidUsername = errors.New("invalid username address")
)

// NormalizePhone returns an international E.164-style identity. A country
// code is required because the domain has no locale/provider context with
// which to interpret local numbers. Formatting punctuation is discarded.
func NormalizePhone(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" || (!strings.HasPrefix(value, "+") && !strings.HasPrefix(value, "00")) {
		return "", errInvalidPhone
	}
	var digits strings.Builder
	for index, r := range value {
		switch {
		case r >= '0' && r <= '9':
			digits.WriteRune(r)
		case r == '+' && index == 0:
			// The plus is validated after separators are removed.
		case r == ' ' || r == '\t' || r == '-' || r == '(' || r == ')' || r == '.' || r == '/':
		default:
			return "", errInvalidPhone
		}
	}
	result := digits.String()
	if strings.HasPrefix(value, "00") {
		result = strings.TrimPrefix(result, "00")
	}
	if len(result) < 7 || len(result) > 15 || result == "" || result[0] == '0' {
		return "", errInvalidPhone
	}
	return "+" + result, nil
}

// NormalizeEmail trims and lowercases an email address. It intentionally does
// not apply provider-specific rules such as Gmail dot removal or plus-tag
// stripping.
func NormalizeEmail(value string) (string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" || strings.ContainsAny(value, " \t\r\n") || strings.Count(value, "@") != 1 {
		return "", errInvalidEmail
	}
	parts := strings.SplitN(value, "@", 2)
	local, domain := parts[0], parts[1]
	if local == "" || domain == "" || len(local) > 64 || len(domain) > 253 || strings.Contains(local, "..") || strings.Contains(domain, "..") {
		return "", errInvalidEmail
	}
	if strings.HasPrefix(local, ".") || strings.HasSuffix(local, ".") || strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return "", errInvalidEmail
	}
	for _, r := range domain {
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '-') {
			return "", errInvalidEmail
		}
	}
	return local + "@" + domain, nil
}

// NormalizeUsername lowercases a provider-neutral username and removes one
// optional leading @. Provider namespaces must be supplied separately when
// the same spelling can identify different users on different systems.
func NormalizeUsername(value string) (string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.TrimPrefix(value, "@")
	if value == "" || len(value) > 128 {
		return "", errInvalidUsername
	}
	for _, r := range value {
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '_' || r == '-') {
			return "", errInvalidUsername
		}
	}
	return value, nil
}

// NormalizeAddress applies the canonicalizer for kind and uses the explicit
// namespace in the uniqueness key. Blank namespaces use the provider-neutral
// "global" namespace.
func NormalizeAddress(kind AddressKind, namespace, value string) (AddressIdentity, error) {
	namespace = strings.ToLower(strings.TrimSpace(namespace))
	if namespace == "" {
		namespace = "global"
	}
	if strings.ContainsAny(namespace, " \t\r\n") {
		return AddressIdentity{}, errors.New("invalid address namespace")
	}
	var (
		normalized string
		err        error
	)
	switch kind {
	case AddressKindPhone:
		normalized, err = NormalizePhone(value)
	case AddressKindEmail:
		normalized, err = NormalizeEmail(value)
	case AddressKindUsername:
		normalized, err = NormalizeUsername(value)
	default:
		return AddressIdentity{}, errors.New("invalid address kind")
	}
	if err != nil {
		return AddressIdentity{}, err
	}
	return AddressIdentity{Kind: kind, Namespace: namespace, Value: normalized}, nil
}

// CanonicalPhone, CanonicalEmail, and CanonicalUsername are descriptive
// aliases for callers that prefer canonical terminology.
func CanonicalPhone(value string) (string, error)    { return NormalizePhone(value) }
func CanonicalEmail(value string) (string, error)    { return NormalizeEmail(value) }
func CanonicalUsername(value string) (string, error) { return NormalizeUsername(value) }
