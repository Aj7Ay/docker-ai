package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	SkipCleanupWarning bool `json:"skip_cleanup_warning"`
}

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".docker-ai-config.json"), nil
}

func LoadConfig() (Config, error) {
	var config Config
	configPath, err := GetConfigPath()
	if err != nil {
		return config, err
	}

	f, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{SkipCleanupWarning: false}, nil // Default config
		}
		return config, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}

func SaveConfig(config Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	f, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(config)
} 