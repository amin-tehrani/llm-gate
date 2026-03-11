package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// ProviderConfig stores credentials for one provider.
type ProviderConfig struct {
	APIKey    string    `yaml:"api_key"`
	Active    bool      `yaml:"active"`
	UpdatedAt time.Time `yaml:"updated_at"`
}

// Config is the top-level configuration.
type Config struct {
	Version   int                       `yaml:"version"`
	Providers map[string]ProviderConfig `yaml:"providers"`
}

// ConfigDir returns the config directory path.
func ConfigDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "llm-gate")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "llm-gate")
}

// ConfigPath returns the full path to the config file.
func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.yaml")
}

// Load reads the config from disk. Returns a default config if the file doesn't exist.
func Load() (*Config, error) {
	cfg := &Config{
		Version:   1,
		Providers: make(map[string]ProviderConfig),
	}

	data, err := os.ReadFile(ConfigPath())
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if cfg.Providers == nil {
		cfg.Providers = make(map[string]ProviderConfig)
	}
	return cfg, nil
}

// Save writes the config to disk with secure permissions.
func (c *Config) Save() error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(ConfigPath(), data, 0600); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	return nil
}

// SetKey stores or updates an API key for a provider.
func (c *Config) SetKey(providerName, apiKey string) {
	pc := c.Providers[providerName]
	pc.APIKey = apiKey
	pc.UpdatedAt = time.Now()
	c.Providers[providerName] = pc
}

// RemoveKey deletes a provider's credentials.
func (c *Config) RemoveKey(providerName string) {
	delete(c.Providers, providerName)
}

// SetActive marks a provider as active or inactive.
func (c *Config) SetActive(providerName string, active bool) {
	pc := c.Providers[providerName]
	pc.Active = active
	c.Providers[providerName] = pc
}

// GetKey returns the API key for a provider, or empty string if not set.
func (c *Config) GetKey(providerName string) string {
	return c.Providers[providerName].APIKey
}

// IsActive returns whether a provider is currently active.
func (c *Config) IsActive(providerName string) bool {
	return c.Providers[providerName].Active
}

// IsConfigured returns whether a provider has a stored API key.
func (c *Config) IsConfigured(providerName string) bool {
	pc, ok := c.Providers[providerName]
	return ok && pc.APIKey != ""
}
