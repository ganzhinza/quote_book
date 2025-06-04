package config

import (
	"encoding/json"
	"os"
)

type ServerConfig struct {
	Address      string `json:"address"`
	ReadTimeout  int    `json:"read_timeout"`  // В секундах
	WriteTimeout int    `json:"write_timeout"` // В секундах
	IdleTimeout  int    `json:"idle_timeout"`  // В секундах
}

type DatabaseConfig struct {
	Type string `json:"type"`
}

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
}

func MustLoad(fp string) (*Config, error) {
	file, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
