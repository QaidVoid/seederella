package config

import (
	"fmt"
	"os"
	"strings"

	"encoding/json"
	"gopkg.in/yaml.v3"
)

type ColumnConfig struct {
	// Either a faker generator or a fixed value
	Faker     string      `json:"faker,omitempty" yaml:"faker,omitempty"`
	Value     interface{} `json:"value,omitempty" yaml:"value,omitempty"`
	Reference string      `json:"reference,omitempty" yaml:"reference,omitempty"`
	SameAs    string      `json:"same_as,omitempty" yaml:"same_as,omitempty"`
	Transform string      `json:"transform,omitempty" yaml:"transform,omitempty"`
}

type TableConfig struct {
	Count  int                     `json:"count" yaml:"count"`
	Fields map[string]ColumnConfig `json:"fields" yaml:"fields"`
}

type Config struct {
	Driver string                 `json:"driver" yaml:"driver"`
	DSN    string                 `json:"dsn" yaml:"dsn"`
	Schema string                 `json:"schema,omitempty" yaml:"schema,omitempty"`
	Tables map[string]TableConfig `json:"tables" yaml:"tables"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	cfg := &Config{}
	switch {
	case strings.HasSuffix(path, ".yaml"), strings.HasSuffix(path, ".yml"):
		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}
	case strings.HasSuffix(path, ".json"):
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported config format: %s", path)
	}

	return cfg, nil
}
