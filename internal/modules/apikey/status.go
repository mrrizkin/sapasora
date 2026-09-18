package apikey

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type APIKeyStatus int // @name apikey.APIKeyStatus

const (
	APIKeyStatusActive APIKeyStatus = iota
	APIKeyStatusInactive
)

func (a APIKeyStatus) String() string {
	return [...]string{"active", "inactive"}[a]
}

func (a *APIKeyStatus) Scan(value any) error {
	switch value.(string) {
	case "active":
		*a = APIKeyStatusActive
	case "inactive":
		*a = APIKeyStatusInactive
	default:
		return errors.New("invalid api key status")
	}
	return nil
}

func (a APIKeyStatus) Value() (driver.Value, error) {
	return a.String(), nil
}

func (a APIKeyStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

func (a *APIKeyStatus) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return a.Scan(s)
}
