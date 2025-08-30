package config

import (
	"os"

	"github.com/unholyFigaro/shh/internal/domain"
	yaml "gopkg.in/yaml.v2"
)

func LoadConfig(path string) (*domain.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config domain.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func SaveConfig(path string, config *domain.Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
