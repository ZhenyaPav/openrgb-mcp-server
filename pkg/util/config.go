package util

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	}
	OpenRGB struct {
		Host      string `yaml:"host"`
		Port      int    `yaml:"port"`
		TimeoutMS int    `yaml:"timeout_ms"`
	}
}

func LoadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var cfg Config

	if err := yaml.Unmarshal(b, &cfg); err != nil {
		fmt.Println(err)
		return nil, err
	}

	if cfg.OpenRGB.Host == "" {
		cfg.OpenRGB.Host = "localhost"
	}

	if cfg.OpenRGB.Port == 0 {
		cfg.OpenRGB.Port = 6742
	}
	if cfg.OpenRGB.TimeoutMS == 0 {
		cfg.OpenRGB.TimeoutMS = 2000
	}

	if cfg.Server.Name == "" {
		cfg.Server.Name = "OpenRGB MCP"
	}

	if cfg.Server.Version == "" {
		cfg.Server.Version = "1.0.0"
	}

	return &cfg, nil
}
