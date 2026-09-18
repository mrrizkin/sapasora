package devicetoken

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type DeviceTokenStatus int // @name devicetoken.DeviceTokenStatus

const (
	DeviceTokenStatusActive DeviceTokenStatus = iota
	DeviceTokenStatusInactive
)

func (d DeviceTokenStatus) String() string {
	return [...]string{"active", "inactive"}[d]
}

func (d *DeviceTokenStatus) Scan(value any) error {
	switch value.(string) {
	case "active":
		*d = DeviceTokenStatusActive
	case "inactive":
		*d = DeviceTokenStatusInactive
	default:
		return errors.New("invalid device token status")
	}
	return nil
}

func (d DeviceTokenStatus) Value() (driver.Value, error) {
	return d.String(), nil
}

func (d DeviceTokenStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *DeviceTokenStatus) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return d.Scan(s)
}
