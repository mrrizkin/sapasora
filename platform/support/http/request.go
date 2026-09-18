// Package http provides utilities for HTTP requests
package http

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ParseQueryParams parses query parameters into a struct with support for nested and embedded types
func ParseQueryParams(c *fiber.Ctx, dest any) error {
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Pointer || destVal.Elem().Kind() != reflect.Struct {
		return errors.New("destination must be a pointer to a struct")
	}

	structVal := destVal.Elem()
	return parseStruct(c, structVal, "")
}

// parseStruct recursively parses struct fields, handling nested and embedded types
func parseStruct(c *fiber.Ctx, structVal reflect.Value, prefix string) error {
	structType := structVal.Type()

	for i := 0; i < structVal.NumField(); i++ {
		field := structVal.Field(i)
		fieldType := structType.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		// Handle embedded/anonymous fields
		if fieldType.Anonymous {
			if err := handleEmbeddedField(c, field, prefix); err != nil {
				return err
			}
			continue
		}

		// Handle regular fields with query tags
		if err := handleTaggedField(c, field, fieldType, prefix); err != nil {
			return err
		}
	}

	return nil
}

// handleEmbeddedField processes embedded/anonymous struct fields
func handleEmbeddedField(c *fiber.Ctx, field reflect.Value, prefix string) error {
	switch field.Kind() {
	case reflect.Struct:
		// Recursively parse embedded struct
		return parseStruct(c, field, prefix)
	case reflect.Pointer:
		if field.Type().Elem().Kind() == reflect.Struct {
			// Initialize pointer if nil
			if field.IsNil() {
				field.Set(reflect.New(field.Type().Elem()))
			}
			// Recursively parse embedded pointer to struct
			return parseStruct(c, field.Elem(), prefix)
		}
	}
	return nil
}

// handleTaggedField processes fields with query tags
func handleTaggedField(
	c *fiber.Ctx,
	field reflect.Value,
	fieldType reflect.StructField,
	prefix string,
) error {
	queryTag := fieldType.Tag.Get("query")
	if queryTag == "" {
		// Check if it's a nested struct without a query tag
		if field.Kind() == reflect.Struct ||
			(field.Kind() == reflect.Pointer && field.Type().Elem().Kind() == reflect.Struct) {
			return handleNestedStruct(c, field, fieldType, prefix)
		}
		return nil
	}

	tagInfo := parseQueryTag(queryTag)
	if tagInfo.name == "" {
		return nil
	}

	// Build parameter name with prefix for nested structures
	paramName := buildParamName(prefix, tagInfo.name)
	paramValue := c.Query(paramName)

	// Validate required parameters
	if err := validateRequired(paramValue, tagInfo, paramName); err != nil {
		return err
	}

	// Apply default value if needed
	if paramValue == "" && tagInfo.defaultValue != "" {
		paramValue = tagInfo.defaultValue
	}

	// Skip if no value to set
	if paramValue == "" {
		return nil
	}

	// Set the field value
	return setFieldValue(field, paramValue, paramName)
}

// handleNestedStruct processes nested structs without query tags
func handleNestedStruct(
	c *fiber.Ctx,
	field reflect.Value,
	fieldType reflect.StructField,
	prefix string,
) error {
	// Use field name as prefix for nested queries (e.g., "user.name", "address.city")
	nestedPrefix := buildParamName(prefix, strings.ToLower(fieldType.Name))

	switch field.Kind() {
	case reflect.Struct:
		return parseStruct(c, field, nestedPrefix)
	case reflect.Pointer:
		if field.Type().Elem().Kind() == reflect.Struct {
			// Check if any nested parameters exist before initializing
			if hasNestedParams(c, nestedPrefix) {
				if field.IsNil() {
					field.Set(reflect.New(field.Type().Elem()))
				}
				return parseStruct(c, field.Elem(), nestedPrefix)
			}
		}
	}
	return nil
}

// QueryTagInfo holds parsed query tag information
type QueryTagInfo struct {
	name         string
	required     bool
	defaultValue string
}

// parseQueryTag parses query tag and extracts options
func parseQueryTag(queryTag string) QueryTagInfo {
	tagParts := strings.Split(queryTag, ",")
	if len(tagParts) == 0 {
		return QueryTagInfo{}
	}

	info := QueryTagInfo{name: tagParts[0]}

	// Process tag options
	for _, option := range tagParts[1:] {
		switch {
		case option == "required":
			info.required = true
		case strings.HasPrefix(option, "default="):
			info.defaultValue = strings.TrimPrefix(option, "default=")
		}
	}

	return info
}

// buildParamName constructs parameter name with prefix
func buildParamName(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "." + name
}

// validateRequired checks if required parameters are present
func validateRequired(paramValue string, tagInfo QueryTagInfo, paramName string) error {
	if paramValue == "" && tagInfo.required && tagInfo.defaultValue == "" {
		return fmt.Errorf("missing required query parameter: %s", paramName)
	}
	return nil
}

// hasNestedParams checks if any query parameters exist for a nested struct
func hasNestedParams(c *fiber.Ctx, prefix string) bool {
	queries := c.Queries()
	prefixDot := prefix + "."

	for key := range queries {
		if strings.HasPrefix(key, prefixDot) {
			return true
		}
	}
	return false
}

// setFieldValue sets the field value using various methods
func setFieldValue(field reflect.Value, paramValue, paramName string) error {
	// Try Scan method first (for SQL nullable types)
	if err := trySccanMethod(field, paramValue); err == nil {
		return nil
	}

	// Try UnmarshalText method for custom types
	if err := tryUnmarshalText(field, paramValue); err == nil {
		return nil
	}

	// Handle standard types
	return setStandardType(field, paramValue, paramName)
}

// trySccanMethod attempts to use the Scan method
func trySccanMethod(field reflect.Value, paramValue string) error {
	scanMethod := field.Addr().MethodByName("Scan")
	if !scanMethod.IsValid() {
		return errors.New("scan method not available")
	}

	result := scanMethod.Call([]reflect.Value{reflect.ValueOf(paramValue)})
	if len(result) > 0 && !result[0].IsNil() {
		return result[0].Interface().(error)
	}
	return nil
}

// tryUnmarshalText attempts to use the UnmarshalText method
func tryUnmarshalText(field reflect.Value, paramValue string) error {
	unmarshalMethod := field.Addr().MethodByName("UnmarshalText")
	if !unmarshalMethod.IsValid() {
		return errors.New("unmarshal method not available")
	}

	result := unmarshalMethod.Call([]reflect.Value{reflect.ValueOf([]byte(paramValue))})
	if len(result) > 0 && !result[0].IsNil() {
		return result[0].Interface().(error)
	}
	return nil
}

// setStandardType handles built-in Go types
func setStandardType(field reflect.Value, paramValue, paramName string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(paramValue)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(paramValue, 10, field.Type().Bits())
		if err != nil {
			return fmt.Errorf(
				"invalid int value %q for parameter %s: %w",
				paramValue,
				paramName,
				err,
			)
		}
		field.SetInt(intValue)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(paramValue, 10, field.Type().Bits())
		if err != nil {
			return fmt.Errorf(
				"invalid uint value %q for parameter %s: %w",
				paramValue,
				paramName,
				err,
			)
		}
		field.SetUint(uintValue)

	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(paramValue, field.Type().Bits())
		if err != nil {
			return fmt.Errorf(
				"invalid float value %q for parameter %s: %w",
				paramValue,
				paramName,
				err,
			)
		}
		field.SetFloat(floatValue)

	case reflect.Bool:
		boolValue, err := strconv.ParseBool(paramValue)
		if err != nil {
			return fmt.Errorf(
				"invalid bool value %q for parameter %s: %w",
				paramValue,
				paramName,
				err,
			)
		}
		field.SetBool(boolValue)

	case reflect.Pointer:
		// Handle pointer types
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		return setStandardType(field.Elem(), paramValue, paramName)

	default:
		return fmt.Errorf("unsupported type %q for parameter %s", field.Type().String(), paramName)
	}

	return nil
}
