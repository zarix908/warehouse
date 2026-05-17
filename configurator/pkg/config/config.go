package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ConfigsDir string `yaml:"configs_dir"`
	TargetFile string `yaml:"target"`
}

func Read(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if cfg.ConfigsDir == "" {
		return nil, fmt.Errorf("configs_dir is required")
	}

	if cfg.TargetFile == "" {
		return nil, fmt.Errorf("target is required")
	}

	return &cfg, nil
}
