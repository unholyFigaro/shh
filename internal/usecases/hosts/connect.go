package hosts

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/unholyFigaro/shh/internal/config"
)

func ConnectToHostByName(ctx context.Context, name string) error {
	path := config.GetConfigPath()
	cfg, err := config.LoadConfig(path)
	if err != nil {
		return err
	}
	h, ok := cfg.Hosts[name]
	if !ok {
		return fmt.Errorf("host not found: %s", name)
	}
	port := h.Port
	if port == 0 {
		port = 22
	}

	args := []string{"-p", strconv.Itoa(port), h.Host}

	if h.User != "" {
		args = append([]string{"-l", h.User}, args...)
	}
	cmd := exec.CommandContext(ctx, "ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
