package config

import (
	"fmt"
	"maps"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type ConfigManager struct {
	mu     sync.RWMutex
	values map[string]any
}

func NewConfigManager() Config {
	cm := &ConfigManager{
		values: make(map[string]any),
	}
	return cm
}

func (cm *ConfigManager) Get(key string, defaultValue ...any) any {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if val, exists := cm.values[key]; exists {
		return val
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return nil
}

func (cm *ConfigManager) GetString(key string, defaultValue ...string) string {
	val := cm.Get(key)
	if val == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}
	return fmt.Sprintf("%v", val)
}

func (cm *ConfigManager) GetInt(key string, defaultValue ...int) int {
	val := cm.Get(key)
	if val == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	switch v := val.(type) {
	case int:
		return v
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

func (cm *ConfigManager) GetBool(key string, defaultValue ...bool) bool {
	val := cm.Get(key)
	if val == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}

	switch v := val.(type) {
	case bool:
		return v
	case string:
		return strings.ToLower(v) == "true" || v == "1" || strings.ToLower(v) == "yes"
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return false
}

func (cm *ConfigManager) GetFloat64(key string, defaultValue ...float64) float64 {
	val := cm.Get(key)
	if val == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	switch v := val.(type) {
	case float64:
		return v
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

func (cm *ConfigManager) GetDuration(key string, defaultValue ...time.Duration) time.Duration {
	val := cm.Get(key)
	if val == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	switch v := val.(type) {
	case time.Duration:
		return v
	case string:
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return 0
}

func (cm *ConfigManager) Set(key string, value any) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.values[key] = value
}

func (cm *ConfigManager) Has(key string) bool {
	return cm.Get(key) != nil
}

func (cm *ConfigManager) All() map[string]any {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	result := make(map[string]any)
	maps.Copy(result, cm.values)
	return result
}

func (cm *ConfigManager) LoadStruct(name string, cfg any) error {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// Handle nested structs and pointers to structs
		if field.Type.Kind() == reflect.Struct && field.Type != reflect.TypeOf(time.Duration(0)) {
			nestedName := name // Create a copy
			if nameTag := field.Tag.Get("name"); nameTag != "" {
				nestedName = fmt.Sprintf("%s.%s", name, nameTag)
			}
			if err := cm.LoadStruct(nestedName, fieldValue.Addr().Interface()); err != nil {
				return fmt.Errorf("%s: %w", field.Name, err)
			}
			continue
		} else if field.Type.Kind() == reflect.Pointer && field.Type.Elem().Kind() == reflect.Struct {
			if fieldValue.IsNil() {
				fieldValue.Set(reflect.New(field.Type.Elem()))
			}
			nestedName := name // Create a copy
			if nameTag := field.Tag.Get("name"); nameTag != "" {
				nestedName = fmt.Sprintf("%s.%s", name, nameTag)
			}
			if err := cm.LoadStruct(nestedName, fieldValue.Interface()); err != nil {
				return fmt.Errorf("%s: %w", field.Name, err)
			}
			continue
		}

		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue
		}

		parts := strings.Split(envTag, ",")
		envVar := parts[0]
		var defaultValue string
		var isRequired bool

		for _, part := range parts[1:] {
			switch {
			case part == "required":
				isRequired = true
			case strings.HasPrefix(part, "default="):
				defaultValue = strings.TrimPrefix(part, "default=")
			}
		}

		value, exists := os.LookupEnv(envVar)
		if !exists {
			if defaultValue != "" {
				value = defaultValue
			} else if isRequired {
				return fmt.Errorf("required config %q not set", envVar)
			} else {
				continue
			}
		}

		var fieldName string
		if nameTag := field.Tag.Get("name"); nameTag != "" {
			fieldName = nameTag
		} else {
			fieldName = field.Name
		}

		if err := cm.setFieldValue(fmt.Sprintf("%s.%s", name, fieldName), field, fieldValue, value); err != nil {
			return fmt.Errorf("%s: %w", field.Name, err)
		}
	}

	return nil
}

func (cm *ConfigManager) setFieldValue(
	name string,
	field reflect.StructField,
	fieldValue reflect.Value,
	value string,
) error {
	// Handle pointers to basic types
	if fieldValue.Kind() == reflect.Pointer {
		elemType := fieldValue.Type().Elem()
		switch elemType.Kind() {
		case reflect.String,
			reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64,
			reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64,
			reflect.Bool,
			reflect.Float32,
			reflect.Float64,
			reflect.Struct: // For time.Duration which is a named type based on int64
			if fieldValue.IsNil() {
				fieldValue.Set(reflect.New(elemType))
			}
			fieldValue = fieldValue.Elem()
		default:
			return fmt.Errorf("unsupported pointer type %q for field %q", elemType, field.Name)
		}
	}

	if !fieldValue.CanSet() {
		return fmt.Errorf("cannot set field %q", field.Name)
	}

	// Check if the field type is time.Duration (which is a named type based on int64)
	if field.Type == reflect.TypeOf(time.Duration(0)) {
		durationValue, err := time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("field %q: invalid duration value %q: %w", field.Name, value, err)
		}
		fieldValue.SetInt(int64(durationValue))
		cm.Set(name, durationValue)
		return nil
	}

	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(value)
		cm.Set(name, value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("field %q: invalid int value %q: %w", field.Name, value, err)
		}
		fieldValue.SetInt(intValue)
		cm.Set(name, value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("field %q: invalid uint value %q: %w", field.Name, value, err)
		}
		fieldValue.SetUint(uintValue)
		cm.Set(name, value)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("field %q: invalid bool value %q: %w", field.Name, value, err)
		}
		fieldValue.SetBool(boolValue)
		cm.Set(name, value)
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("field %q: invalid float value %q: %w", field.Name, value, err)
		}
		fieldValue.SetFloat(floatValue)
		cm.Set(name, value)
	default:
		return fmt.Errorf("unsupported type %q for field %q", fieldValue.Type(), field.Name)
	}
	return nil
}
