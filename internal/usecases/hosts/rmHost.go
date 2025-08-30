package hosts

import (
	"fmt"
	"io"

	"github.com/unholyFigaro/shh/internal/config"
)

func RemoveHost(w io.Writer, name string) error {
	err := removeHostByName(name)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Host %q removed successfully!\n", name)
	return nil
}

func removeHostByName(name string) error {
	path := config.GetConfigPath()
	cfg, err := config.LoadConfig(path)
	if err != nil {
		return err
	}

	if _, exists := cfg.Hosts[name]; !exists {
		return fmt.Errorf("host with name %q does not exist", name)
	}

	delete(cfg.Hosts, name)
	if err := config.SaveConfig(path, cfg); err != nil {
		return err
	}

	return nil
}
