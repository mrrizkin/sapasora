package satpam

import (
	"reflect"
)

type Attribute struct {
	original any
	data     map[string]any
}

func NewAttributeCarrier(v any) *Attribute {
	if v == nil {
		return &Attribute{
			original: nil,
			data:     map[string]any{},
		}
	}

	if m, ok := v.(map[string]any); ok {
		return &Attribute{
			original: v,
			data:     m,
		}
	}

	data := structToMap(v)
	return &Attribute{
		original: v,
		data:     data,
	}
}

func structToMap(v any) map[string]any {
	result := make(map[string]any)

	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// Handle pointer
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return result
		}
		val = val.Elem()
		typ = typ.Elem()
	}

	// Only process structs
	if val.Kind() != reflect.Struct {
		return result
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get JSON tag name
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			continue
		}

		// Parse JSON tag (handle "name,omitempty" format)
		tagName := jsonTag
		for idx := 0; idx < len(jsonTag); idx++ {
			if jsonTag[idx] == ',' {
				tagName = jsonTag[:idx]
				break
			}
		}

		// Skip if json tag is "-"
		if tagName == "-" {
			continue
		}

		// Store the actual field value with its original type
		if fieldVal.CanInterface() {
			result[tagName] = fieldVal.Interface()
		}
	}

	return result
}

func (a *Attribute) GetAttribute(key string) any {
	if v, ok := a.data[key]; ok {
		return v
	}
	return nil
}

// Is checks if the original value has the same type as the provided value
func (a *Attribute) Is(target any) bool {
	if a.original == nil {
		return target == nil
	}

	originalType := reflect.TypeOf(a.original)
	targetType := reflect.TypeOf(target)

	return originalType == targetType
}

func (a *Attribute) GetString(key string) string {
	if v, ok := a.data[key].(string); ok {
		return v
	}
	return ""
}

func (a *Attribute) GetBool(key string) bool {
	if v, ok := a.data[key].(bool); ok {
		return v
	}
	return false
}

func (a *Attribute) GetFloat64(key string) float64 {
	if v, ok := a.data[key].(float64); ok {
		return v
	}
	return 0
}

func (a *Attribute) GetInt(key string) int {
	if v, ok := a.data[key].(int); ok {
		return v
	}
	return 0
}

func (a *Attribute) GetUint(key string) uint {
	if v, ok := a.data[key].(uint); ok {
		return v
	}
	return 0
}
