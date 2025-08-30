package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	envConfigVar   = "SHH_CONFIG"
	appDirName     = "shh"
	configFileName = "hosts.yaml"
)

func GetConfigPath() string {
	var p string

	if env := os.Getenv(envConfigVar); env != "" {
		p = expandTilde(env)
	} else if dir, err := os.UserConfigDir(); err == nil && dir != "" {
		p = filepath.Join(dir, appDirName, configFileName)
	} else if home, err := os.UserHomeDir(); err == nil && home != "" {
		p = filepath.Join(home, ".config", appDirName, configFileName)
	} else {
		p = filepath.Join(".", configFileName)
	}

	if err := EnsureConfigDir(p); err != nil {
		alt := filepath.Join(".", configFileName)
		if alt != p {
			_ = EnsureConfigDir(alt)
			return alt
		}
	}

	return p
}

func EnsureConfigDir(path string) error {
	if fi, err := os.Stat(path); err == nil {
		if fi.IsDir() {
			return fmt.Errorf("%s exists but is a directory", path)
		}
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	dir := filepath.Dir(path)
	if fi, err := os.Stat(dir); err == nil {
		if !fi.IsDir() {
			return fmt.Errorf("%s exists but is not a directory", dir)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	} else if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil
		}
		return err
	}
	defer f.Close()

	const defaultYAML = "version: \"0.0.1\"\n\nhosts: {}\n"
	if _, err := f.WriteString(defaultYAML); err != nil {
		return err
	}
	return f.Sync()
}

func expandTilde(p string) string {
	if p == "" || p[0] != '~' {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return p
	}
	if p == "~" {
		return home
	}
	if strings.HasPrefix(p, "~"+string(os.PathSeparator)) {
		return filepath.Join(home, p[2:])
	}
	return p
}
