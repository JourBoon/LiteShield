package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type ClientYAML struct {
	ID         string `yaml:"id"`
	BackendURL string `yaml:"backendURL"`
	DefaultTTL string `yaml:"defaultTTL"`
	RateRPS    int    `yaml:"rateRPS"`
}

type FileConfig struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Clients []ClientYAML `yaml:"clients"`
}

func LoadConfig(path string) (*FileConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg FileConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *ClientYAML) TTL() time.Duration {
	d, _ := time.ParseDuration(c.DefaultTTL)
	if d == 0 {
		d = 60 * time.Second
	}
	return d
}
