package domain

type Host struct {
	Host string `yaml:"host"`
	User string `yaml:"user,omitempty"`
	Port int    `yaml:"port,omitempty"`
}
