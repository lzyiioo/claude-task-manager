package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	Isolation struct {
		Default     string `json:"default"`
		DockerImage string `json:"dockerImage"`
	} `json:"isolation"`

	Logging struct {
		EventRetentionDays int    `json:"eventRetentionDays"`
		MaxEventFileSize   int    `json:"maxEventFileSizeMB"`
		LogLevel           string `json:"logLevel"`
	} `json:"logging"`

	Permission struct {
		AutoAllowRead   bool     `json:"autoAllowRead"`
		AutoAllowGlob   bool     `json:"autoAllowGlob"`
		BlockedCommands []string `json:"blockedCommands"`
	} `json:"permission"`

	Concurrency struct {
		MaxTasks      int `json:"maxTasks"`
		MaxCpuPercent int `json:"maxCpuPercent"`
	} `json:"concurrency"`

	// Runtime fields (not serialized)
	configPath string
}

// Default returns a new config with default values
func Default() *Config {
	var cfg Config
	if err := json.Unmarshal([]byte(DefaultConfigJSON), &cfg); err != nil {
		panic("invalid default config: " + err.Error())
	}
	return &cfg
}

// Load reads config from the default path, creating it if not exists
func Load() (*Config, error) {
	configPath, err := configFilePath()
	if err != nil {
		return nil, err
	}

	// Create default config if not exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg := Default()
		cfg.configPath = configPath
		if err := cfg.Save(); err != nil {
			return nil, err
		}
		return cfg, nil
	}

	return LoadFrom(configPath)
}

// LoadFrom reads config from a specific path
func LoadFrom(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := Default()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	cfg.configPath = path
	return cfg, nil
}

// Save writes the config to disk
func (c *Config) Save() error {
	if c.configPath == "" {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(c.configPath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.configPath, data, 0644)
}

// configFilePath returns the path to the config file
func configFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".claude-task-manager", "config.json"), nil
}
