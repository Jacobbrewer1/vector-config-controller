package vector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
)

// Config is the configuration supplied to vector. It is converted into JSON format before being written to a file.
type Config struct {
	internal internalConfig
}

// internalConfig is the internal representation of the vector configuration.
type internalConfig struct {
	SecretBackends map[string]map[string]any `json:"secret,omitempty"`
	Sources        map[string]map[string]any `json:"sources"`
	Transforms     map[string]map[string]any `json:"transforms,omitempty"`
	Sinks          map[string]map[string]any `json:"sinks"`
}

// NewConfig creates an empty vector configuration.
func NewConfig() *Config {
	return &Config{
		internal: internalConfig{
			SecretBackends: make(map[string]map[string]any),
			Sources:        make(map[string]map[string]any),
			Transforms:     make(map[string]map[string]any),
			Sinks:          make(map[string]map[string]any),
		},
	}
}

// AddSecretBackend adds the specified configuration as a secret backend under key.
//
// The backendName is what Vector will refer to when using the secret backend.
// https://vector.dev/highlights/2022-07-07-secrets-management/
func (c *Config) AddSecretBackend(backendName string, cfg map[string]any) {
	if _, ok := c.internal.SecretBackends[backendName]; ok {
		panic(fmt.Sprintf("secret backend key '%s' already added to configuration", backendName))
	}

	c.internal.SecretBackends[backendName] = cfg
}

// AddSourceUntyped adds the specified configuration as a vector source under key.
func (c *Config) AddSourceUntyped(key string, cfg map[string]any) {
	if _, ok := c.internal.Sources[key]; ok {
		panic(fmt.Sprintf("source key '%s' already added to configuration", key))
	}

	c.internal.Sources[key] = cfg
}

// AddTransformUntyped adds the specified configuration as a vector transform under key.
func (c *Config) AddTransformUntyped(key string, cfg map[string]any) {
	if _, ok := c.internal.Transforms[key]; ok {
		panic(fmt.Sprintf("transforms key '%s' already added to configuration", key))
	}

	c.internal.Transforms[key] = cfg
}

// AddSinkUntyped adds the specified configuration as a vector sink under key.
func (c *Config) AddSinkUntyped(key string, cfg map[string]any) {
	if _, ok := c.internal.Sinks[key]; ok {
		panic(fmt.Sprintf("sinks key '%s' already added to configuration", key))
	}

	c.internal.Sinks[key] = cfg
}

// JSON returns the JSON representation of the configuration.
func (c *Config) JSON() (string, error) {
	result := bytes.NewBuffer(nil)
	if err := json.NewEncoder(result).Encode(c.internal); err != nil {
		return "", fmt.Errorf("error encoding config: %w", err)
	}
	return result.String(), nil
}

// Sources returns a copy of the current sources.
func (c *Config) Sources() map[string]map[string]any {
	result := make(map[string]map[string]any)
	maps.Copy(result, c.internal.Sources)
	return result
}

// Transforms returns a copy of the current transforms.
func (c *Config) Transforms() map[string]map[string]any {
	result := make(map[string]map[string]any)
	maps.Copy(result, c.internal.Transforms)
	return result
}

// Sinks returns a copy of the current sinks.
func (c *Config) Sinks() map[string]map[string]any {
	result := make(map[string]map[string]any)
	maps.Copy(result, c.internal.Sinks)
	return result
}
