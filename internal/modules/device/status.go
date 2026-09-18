package device

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type DeviceStatus int // @name device.DeviceStatus

const (
	DeviceStatusInactive DeviceStatus = iota
	DeviceStatusActive
	DeviceStatusConnected
	DeviceStatusDisconnected
)

func (d DeviceStatus) String() string {
	return [...]string{"inactive", "active", "connected", "disconnected"}[d]
}

func (d *DeviceStatus) Scan(value any) error {
	switch value.(string) {
	case "inactive":
		*d = DeviceStatusInactive
	case "active":
		*d = DeviceStatusActive
	case "connected":
		*d = DeviceStatusConnected
	case "disconnected":
		*d = DeviceStatusDisconnected
	default:
		return errors.New("invalid device status")
	}
	return nil
}

func (d DeviceStatus) Value() (driver.Value, error) {
	return d.String(), nil
}

func (d DeviceStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *DeviceStatus) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return d.Scan(s)
}
