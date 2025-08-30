package hosts

import (
	"context"
	"io"

	"github.com/unholyFigaro/shh/internal/config"
	"github.com/unholyFigaro/shh/internal/ui"
)

func ListHosts(ctx context.Context, w io.Writer) error {
	path := config.GetConfigPath()

	data, err := config.LoadConfig(path)
	if err != nil {
		return err
	}

	return ui.PrintHosts(w, data.Hosts)
}
