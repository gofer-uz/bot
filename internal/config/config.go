package config

import (
	"os"
	"strings"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

// Config represents application configuration.

type Config struct {
	Server   ServerConfig   `toml:"server"`
	Telegram TelegramConfig `toml:"telegram"`
}

type ServerConfig struct {
	PublicURL string `toml:"public_url"`
	Listen    string `toml:"listen"`
}

type TelegramConfig struct {
	Token string `toml:"token"`
}

var (
	cfg     Config
	cfgOnce sync.Once
)

// Load reads configuration from the provided path (toml format).
func Load(path string) (Config, error) {
	var err error
	cfgOnce.Do(func() {
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			err = readErr
			return
		}
		if unmarshalErr := toml.Unmarshal(data, &cfg); unmarshalErr != nil {
			err = unmarshalErr
			return
		}
		if err == nil {
			// trim whitespace which often sneaks into toml strings
			cfg.Server.PublicURL = strings.TrimSpace(cfg.Server.PublicURL)
		}
	})
	return cfg, err
}
