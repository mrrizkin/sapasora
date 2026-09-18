package device

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type DeviceType int // @name device.DeviceType

const (
	DeviceTypeWhatsapp DeviceType = iota
	DeviceTypeTelegram
)

func (d DeviceType) String() string {
	return [...]string{"whatsapp", "telegram"}[d]
}

func (d *DeviceType) Scan(value any) error {
	switch value.(string) {
	case "whatsapp":
		*d = DeviceTypeWhatsapp
	case "telegram":
		*d = DeviceTypeTelegram
	default:
		return errors.New("invalid device type")
	}
	return nil
}

func (d DeviceType) Value() (driver.Value, error) {
	return d.String(), nil
}

func (d DeviceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *DeviceType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	return d.Scan(s)
}
