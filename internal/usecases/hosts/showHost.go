package hosts

import (
	"fmt"
	"io"

	"github.com/unholyFigaro/shh/internal/config"
	"github.com/unholyFigaro/shh/internal/domain"
	"github.com/unholyFigaro/shh/internal/ui"
)

func ShowHostsByName(w io.Writer, names []string) error {
	notFound := make([]string, 0)
	found := make(map[string]domain.Host)
	for _, name := range names {
		host, err := FindHostByName(name)
		if err != nil {
			notFound = append(notFound, name)
			continue
		}
		found[name] = *host
	}
	ui.PrintHosts(w, found)
	return nil
}

func FindHostByName(name string) (*domain.Host, error) {
	cfg, err := config.LoadConfig(config.GetConfigPath())
	if err != nil {
		return nil, err
	}

	if host, exist := cfg.Hosts[name]; exist {
		return &host, nil
	}

	return nil, fmt.Errorf("host with name %q not found", name)
}
