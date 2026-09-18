package hash

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

// Error definitions
var (
	ErrInvalidPBKDF2Hash = errors.New("the encoded pbkdf2 hash is not in the correct format")
	ErrPBKDF2Iterations  = errors.New("iterations must be a positive integer")
)

// PBKDF2Params represents the parameters used for PBKDF2 hashing.
type PBKDF2Params struct {
	Iterations int
	KeyLen     int
	SaltLen    int
}

// DefaultPBKDF2Params provides recommended default values for PBKDF2 hashing.
var DefaultPBKDF2Params = PBKDF2Params{
	Iterations: 10000,
	KeyLen:     32,
	SaltLen:    16,
}

// PBKDF2 creates a new PBKDF2 hash from a plaintext password.
// Optional custom parameters can be provided, otherwise defaults are used.
func PBKDF2(plaintext string, customParams ...PBKDF2Params) (string, error) {
	params := DefaultPBKDF2Params
	if len(customParams) > 0 {
		params = mergePBKDF2Params(customParams[0])
	}

	// Generate salt
	salt := NanoID(int(params.SaltLen))

	// Hash the password
	hash := pbkdf2.Key(
		[]byte(plaintext),
		[]byte(salt),
		params.Iterations,
		params.KeyLen,
		sha256.New,
	)

	// Encode salt and hash using base64
	encodedSalt := base64.RawStdEncoding.EncodeToString([]byte(salt))
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	// Format the encoded hash
	encodedFull := fmt.Sprintf(
		"$pbkdf2$v=1$iter=%d,keyLen=%d$%s$%s",
		params.Iterations,
		params.KeyLen,
		encodedSalt,
		encodedHash,
	)

	return encodedFull, nil
}

// PBKDF2Verify checks if a plaintext password matches an encoded hash.
func PBKDF2Verify(plaintext, encodedHash string) bool {
	params, salt, hash, err := decodePBKDF2Hash(encodedHash)
	if err != nil {
		return false
	}

	// Generate a new hash with the same parameters and salt
	newHash := pbkdf2.Key(
		[]byte(plaintext),
		salt,
		params.Iterations,
		params.KeyLen,
		sha256.New,
	)

	// Compare hashes in constant time to prevent timing attacks
	return subtle.ConstantTimeCompare(hash, newHash) == 1
}

// decodePBKDF2Hash extracts the parameters, salt, and hash from an encoded hash string.
func decodePBKDF2Hash(encodedHash string) (*PBKDF2Params, []byte, []byte, error) {
	// Split the hash into parts
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 5 {
		return nil, nil, nil, ErrInvalidPBKDF2Hash
	}

	// Check version
	if parts[0] != "pbkdf2" {
		return nil, nil, nil, ErrInvalidPBKDF2Hash
	}

	// Extract parameters
	params := &PBKDF2Params{}
	if _, err := fmt.Sscanf(
		parts[1],
		"iter=%d,keyLen=%d",
		&params.Iterations,
		&params.KeyLen,
	); err != nil {
		return nil, nil, nil, err
	}

	// Decode salt
	salt, err := base64.RawStdEncoding.Strict().DecodeString(parts[2])
	if err != nil {
		return nil, nil, nil, err
	}
	params.SaltLen = len(salt)

	// Decode hash
	hash, err := base64.RawStdEncoding.Strict().DecodeString(parts[3])
	if err != nil {
		return nil, nil, nil, err
	}
	params.KeyLen = len(hash)

	return params, salt, hash, nil
}

// mergePBKDF2Params combines custom parameters with defaults for any unspecified values.
func mergePBKDF2Params(custom PBKDF2Params) PBKDF2Params {
	result := custom

	// Apply defaults for any zero values
	if result.Iterations == 0 {
		result.Iterations = DefaultPBKDF2Params.Iterations
	}
	if result.KeyLen == 0 {
		result.KeyLen = DefaultPBKDF2Params.KeyLen
	}
	if result.SaltLen == 0 {
		result.SaltLen = DefaultPBKDF2Params.SaltLen
	}

	return result
}
