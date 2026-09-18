// Package hash provides utilities for password hashing using Argon2.
package hash

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Error definitions
var (
	ErrInvalidArgon2Hash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleArgon2Version = errors.New("incompatible version of argon2")
)

// Argon2Params represents the parameters used for Argon2id hashing.
type Argon2Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	KeyLen      uint32
	SaltLen     uint32
}

// DefaultArgon2Params provides recommended default values for Argon2id hashing.
var DefaultArgon2Params = Argon2Params{
	Memory:      1 << 16, // 64MB
	Iterations:  1,
	Parallelism: 1,
	KeyLen:      32,
	SaltLen:     16,
}

// Argon2 creates a new Argon2id hash from a plaintext password.
// Optional custom parameters can be provided, otherwise defaults are used.
func Argon2(plaintext string, customParams ...Argon2Params) (string, error) {
	params := DefaultArgon2Params
	if len(customParams) > 0 {
		params = mergeArgon2Params(customParams[0])
	}

	// Generate salt
	salt := NanoID(int(params.SaltLen))

	// Hash the password
	hash := argon2.IDKey(
		[]byte(plaintext),
		[]byte(salt),
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLen,
	)

	// Encode salt and hash using base64
	encodedSalt := base64.RawStdEncoding.EncodeToString([]byte(salt))
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	// Format the encoded hash
	encodedFull := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		params.Memory,
		params.Iterations,
		params.Parallelism,
		encodedSalt,
		encodedHash,
	)

	return encodedFull, nil
}

// Argon2Verify checks if a plaintext password matches an encoded hash.
func Argon2Verify(plaintext, encodedHash string) bool {
	params, salt, hash, err := decodeArgon2Hash(encodedHash)
	if err != nil {
		return false
	}

	// Generate a new hash with the same parameters and salt
	newHash := argon2.IDKey(
		[]byte(plaintext),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLen,
	)

	// Compare hashes in constant time to prevent timing attacks
	return subtle.ConstantTimeCompare(hash, newHash) == 1
}

// decodeArgon2Hash extracts the parameters, salt, and hash from an encoded hash string.
func decodeArgon2Hash(encodedHash string) (*Argon2Params, []byte, []byte, error) {
	// Split the hash into parts
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, ErrInvalidArgon2Hash
	}

	// Check version
	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleArgon2Version
	}

	// Extract parameters
	params := &Argon2Params{}
	if _, err := fmt.Sscanf(
		parts[3],
		"m=%d,t=%d,p=%d",
		&params.Memory,
		&params.Iterations,
		&params.Parallelism,
	); err != nil {
		return nil, nil, nil, err
	}

	// Decode salt
	salt, err := base64.RawStdEncoding.Strict().DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}
	params.SaltLen = uint32(len(salt))

	// Decode hash
	hash, err := base64.RawStdEncoding.Strict().DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, err
	}
	params.KeyLen = uint32(len(hash))

	return params, salt, hash, nil
}

// mergeArgon2Params combines custom parameters with defaults for any unspecified values.
func mergeArgon2Params(custom Argon2Params) Argon2Params {
	result := custom

	// Apply defaults for any zero values
	if result.Memory == 0 {
		result.Memory = DefaultArgon2Params.Memory
	}
	if result.Iterations == 0 {
		result.Iterations = DefaultArgon2Params.Iterations
	}
	if result.Parallelism == 0 {
		result.Parallelism = DefaultArgon2Params.Parallelism
	}
	if result.KeyLen == 0 {
		result.KeyLen = DefaultArgon2Params.KeyLen
	}
	if result.SaltLen == 0 {
		result.SaltLen = DefaultArgon2Params.SaltLen
	}

	return result
}
