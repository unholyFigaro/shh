package hosts

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/unholyFigaro/shh/internal/config"
	"github.com/unholyFigaro/shh/internal/domain"
	"github.com/unholyFigaro/shh/internal/validation"
)

func AddHost(ctx context.Context, params map[string]any) error {
	err := validation.Validate(params, validation.HostSchema())
	if err != nil {
		return err
	}
	user := params["user"].(string)
	host := params["host"].(string)
	port := params["port"].(int)
	force := params["force"].(bool)
	name := params["name"].(string)

	if port == 0 {
		port = 22
	}

	path := config.GetConfigPath()
	cfg, err := config.LoadConfig(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg = &domain.Config{
				Version: "1.0",
				Hosts:   map[string]domain.Host{},
			}
			_ = os.MkdirAll(filepath.Dir(path), 0644)
		} else {
			return err
		}
	}
	if _, exist := cfg.Hosts[name]; exist && !force {
		return fmt.Errorf("host with name %q already exists. Use --force to overwrite", name)
	}

	cfg.Hosts[name] = domain.Host{
		Host: host,
		Port: port,
		User: user,
	}

	if err := config.SaveConfig(path, cfg); err != nil {
		return err
	}

	return nil
}
