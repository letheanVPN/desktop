package config

import (
	"fmt"
	"reflect"
	"strings"
)

// Config holds the resolved paths and user-configurable settings for the application.
type Config struct {
	// --- Dynamic Paths (not stored in config.json) ---
	DataDir       string `json:"-"`
	ConfigDir     string `json:"-"`
	CacheDir      string `json:"-"`
	WorkspacesDir string `json:"-"`
	RootDir       string `json:"-"`
	UserHomeDir   string `json:"-"`
	IsNew         bool   `json:"-"` // Flag indicating if the config was newly created.

	// --- Storable Settings (persisted in config.json) ---
	DefaultRoute string   `json:"defaultRoute,omitempty"`
	Features     []string `json:"features,omitempty"`
	Language     string   `json:"language,omitempty"`
}

// Key retrieves a configuration value by its key. It checks JSON tags and field names (case-insensitive).
func (c *Config) Key(key string) (interface{}, error) {
	// Use reflection to inspect the struct fields.
	val := reflect.ValueOf(c).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name

		// Check the field name first.
		if strings.EqualFold(fieldName, key) {
			return val.Field(i).Interface(), nil
		}

		// Then check the `json` tag.
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			jsonName := strings.Split(jsonTag, ",")[0]
			if strings.EqualFold(jsonName, key) {
				return val.Field(i).Interface(), nil
			}
		}
	}

	return nil, fmt.Errorf("key '%s' not found in config", key)
}
