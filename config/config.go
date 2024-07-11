package config

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"os"

	"github.com/shamil/weather/internal/infrastructure/database"
)

type Config struct {
	Server `yaml:"server"`
}

type Server struct {
	Token    string       `yaml:"token"`
	Database database.Opt `yaml:"database"`
}

func New(filepath string) (*Config, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	d := yaml.NewDecoder(bytes.NewReader(content))
	if err = d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}
