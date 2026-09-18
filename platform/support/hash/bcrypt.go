package hash

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// BCrypt-specific errors
var (
	ErrInvalidBcryptHash = errors.New("the encoded bcrypt hash is not in the correct format")
	ErrBcryptCostTooLow  = errors.New("bcrypt cost parameter is too low")
	ErrBcryptCostTooHigh = errors.New("bcrypt cost parameter is too high")
)

// BCryptParams represents the parameters used for bcrypt hashing.
type BCryptParams struct {
	// Cost determines the computational complexity of the hashing.
	// Must be between bcrypt.MinCost and bcrypt.MaxCost (4-31).
	Cost int
}

// DefaultBCryptParams provides recommended default values for bcrypt hashing.
var DefaultBCryptParams = BCryptParams{
	Cost: 12, // Good balance between security and performance
}

// BCrypt creates a new bcrypt hash from a plaintext password.
// Optional custom parameters can be provided, otherwise defaults are used.
func BCrypt(plaintext string, customParams ...BCryptParams) (string, error) {
	params := DefaultBCryptParams
	if len(customParams) > 0 {
		params = validateBCryptParams(customParams[0])
	}

	// Verify cost is within acceptable range
	if params.Cost < bcrypt.MinCost {
		return "", ErrBcryptCostTooLow
	}
	if params.Cost > bcrypt.MaxCost {
		return "", ErrBcryptCostTooHigh
	}

	// Generate hash with specified cost
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), params.Cost)
	if err != nil {
		return "", err
	}

	// Format the encoded hash to include our package prefix
	encodedHash := fmt.Sprintf("$bcrypt$%s", string(hash))
	return encodedHash, nil
}

// BCryptVerify checks if a plaintext password matches an encoded hash.
func BCryptVerify(plaintext, encodedHash string) bool {
	// Extract the actual bcrypt hash
	hash, err := extractBCryptHash(encodedHash)
	if err != nil {
		return false
	}

	// Compare the password with the hash
	err = bcrypt.CompareHashAndPassword(hash, []byte(plaintext))
	return err == nil
}

// GetBCryptCost extracts the cost parameter from an encoded hash.
func GetBCryptCost(encodedHash string) (int, error) {
	hash, err := extractBCryptHash(encodedHash)
	if err != nil {
		return 0, err
	}

	// Parse the cost from the hash
	cost, err := bcrypt.Cost(hash)
	if err != nil {
		return 0, err
	}

	return cost, nil
}

// extractBCryptHash extracts the raw bcrypt hash from our encoded format.
func extractBCryptHash(encodedHash string) ([]byte, error) {
	// Check if the hash has our prefix
	if strings.HasPrefix(encodedHash, "$bcrypt$") {
		// Remove our prefix
		return []byte(encodedHash[8:]), nil
	}

	// If it doesn't have our prefix, check if it's a valid bcrypt hash already
	if strings.HasPrefix(encodedHash, "$2a$") ||
		strings.HasPrefix(encodedHash, "$2b$") ||
		strings.HasPrefix(encodedHash, "$2y$") {
		return []byte(encodedHash), nil
	}

	return nil, ErrInvalidBcryptHash
}

// validateBCryptParams ensures the provided parameters are within acceptable ranges.
func validateBCryptParams(params BCryptParams) BCryptParams {
	result := params

	// Apply default cost if unspecified or out of range
	if result.Cost < bcrypt.MinCost || result.Cost > bcrypt.MaxCost {
		result.Cost = DefaultBCryptParams.Cost
	}

	return result
}

// NeedsBCryptRehash checks if a hash should be regenerated based on current security standards.
// Returns true if the hash uses outdated parameters.
func NeedsBCryptRehash(encodedHash string, customParams ...BCryptParams) (bool, error) {
	// Use default params if none provided
	params := DefaultBCryptParams
	if len(customParams) > 0 {
		params = validateBCryptParams(customParams[0])
	}

	// Get current cost from the hash
	currentCost, err := GetBCryptCost(encodedHash)
	if err != nil {
		return false, err
	}

	// Check if the hash needs to be upgraded
	return currentCost < params.Cost, nil
}
