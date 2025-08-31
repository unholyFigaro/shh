package hosts

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/unholyFigaro/shh/internal/config"
	"github.com/unholyFigaro/shh/internal/domain"
)

func ConnectToHostByName(ctx context.Context, name, jumpHostName string) error {
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
	var args []string
	if jumpHostName != "" {
		jumpHost, ok := cfg.Hosts[jumpHostName]
		if !ok {
			return fmt.Errorf("jump host not found: %s", jumpHostName)
		}
		args = append(args, "-J", jumpSpec(jumpHost))
	}

	args = append(args, "-p", strconv.Itoa(port), h.Host)
	if h.User != "" {
		args = append([]string{"-l", h.User}, args...)
	}

	fmt.Printf("%v\n", args)
	cmd := exec.CommandContext(ctx, "ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func jumpSpec(jhost domain.Host) string {
	conString := jhost.Host
	if jhost.User != "" {
		conString = jhost.User + "@" + conString
	}
	if jhost.Port == 0 {
		return conString
	}
	return fmt.Sprintf("%s:%d", conString, jhost.Port)
}
