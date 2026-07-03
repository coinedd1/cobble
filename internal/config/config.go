package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Context   string `toml:"context"`
	Namespace string `toml:"namespace"`
	Container string `toml:"container"`
	Selector  string `toml:"selector"`
	DataDir   string `toml:"dataDir"`
	RconMode  string `toml:"rconMode"`
}

func defaults() Config {
	return Config{
		Namespace: "minecraft",
		Container: "minecraft",
		DataDir:   "/data",
		RconMode:  "exec",
	}
}

func Load() (Config, error) {
	cfg := defaults()

	path := os.Getenv("COBBLE CONFIG")
	if path == "" {
		dir, err := os.UserConfigDir()
		if err != nil {
			return cfg, err
		}
		path = filepath.Join(dir, "cobble", "config.toml")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
