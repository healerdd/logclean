package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = "config.json"

func getConfigPath() string {
	exe, _ := os.Executable()
	return filepath.Join(filepath.Dir(exe), configFileName)
}

func LoadConfig() (*Config, error) {
	cfg := &Config{Tasks: []CleanTask{}}
	data, err := os.ReadFile(getConfigPath())
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
	err = json.Unmarshal(data, cfg)
	return cfg, err
}

func SaveConfig(cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(getConfigPath(), data, 0644)
}