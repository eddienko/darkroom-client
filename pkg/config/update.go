package config

import (
	"fmt"
	"reflect"
	"strings"
)

// UpdateField updates a field in the Config struct by name (case-insensitive)
// value must be a string; it converts to the field type if needed.
func (c *Config) UpdateField(fieldName string, value string) error {
	v := reflect.ValueOf(c).Elem() // pointer to struct
	t := v.Type()

	// Find field by name (case-insensitive)
	var field reflect.Value
	found := false
	for i := 0; i < t.NumField(); i++ {
		if strings.EqualFold(t.Field(i).Name, fieldName) {
			field = v.Field(i)
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("field %q not found in Config", fieldName)
	}

	// Only support string fields for now
	if field.Kind() != reflect.String {
		return fmt.Errorf("field %q is not a string field", fieldName)
	}

	if !field.CanSet() {
		return fmt.Errorf("cannot set field %q", fieldName)
	}

	field.SetString(value)
	return nil
}
