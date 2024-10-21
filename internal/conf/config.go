package conf

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Port          string `json:"port"`
	StorePath     string `json:"store_path"`
	StoreDriver   string `json:"store_driver"`
	MigrationPath string `json:"migration_path"`
	Version       string `json:"version"`
	LogFile       string `json:"log_file"`
}

func Load() (*Config, error) {
	conf := &Config{}
	confJSON, err := os.ReadFile("config.json")
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	err = json.Unmarshal(confJSON, conf)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return conf, nil
}
