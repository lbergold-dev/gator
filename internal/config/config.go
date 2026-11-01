package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = "config.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// SetUser writes the current user name to ~/.config/gator/config.json.
// It returns an error if the file cannot be written.
func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user
	return write(*c)
}

// Read loads the configuration from ~/.config/gator/config.json.
// Returns an error if the file does not exist or contains invalid JSON.
func Read() (*Config, error) {
	cfgDir, err := getConfigDir()
	if err != nil {
		return nil, fmt.Errorf("get config dir: %w", err)
	}

	jsonCfg, err := os.ReadFile(filepath.Join(cfgDir, configFileName))
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	cfg := &Config{}

	if err := json.Unmarshal(jsonCfg, cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}

func write(cfg Config) error {
	cfgDir, err := getConfigDir()
	if err != nil {
		return fmt.Errorf("get config dir: %w", err)
	}

	if err := os.MkdirAll(cfgDir, 0700); err != nil {
		return fmt.Errorf("create config file path: %w", err)
	}

	jsonCfg, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(filepath.Join(cfgDir, configFileName), jsonCfg, 0600); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}

	return nil
}

func getConfigDir() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get user home dir: %w", err)
	}

	return filepath.Join(userHomeDir, ".config", "gator"), nil

}
