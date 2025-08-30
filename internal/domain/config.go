package domain

type Config struct {
	Version string          `yaml:"version"`
	Hosts   map[string]Host `yaml:"hosts"`
}
