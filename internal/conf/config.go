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
	confJson, err := os.ReadFile("config.json")
	if err != nil {
		return nil, fmt.Errorf("read config - %v", err)
	}

	err = json.Unmarshal(confJson, conf)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config - %v", err)
	}

	return conf, nil
}
