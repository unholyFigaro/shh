package completion

import (
	"sort"
	"strings"

	"github.com/unholyFigaro/shh/internal/config"
)

func HostNamesByPrefix(prefix string) ([]string, error) {
	cfg, err := config.LoadConfig(config.GetConfigPath())
	if cfg == nil || err != nil || cfg.Hosts == nil {
		return nil, nil
	}

	out := make([]string, 0, len(cfg.Hosts))
	for name := range cfg.Hosts {
		if prefix == "" || strings.HasPrefix(name, prefix) {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out, nil
}
