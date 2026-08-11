package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type AppConfig struct {
	Email                 string  `toml:"email"`
	Password              string  `toml:"password"`
	Region                string  `toml:"region"`
	UpdateIntervalMinutes int     `toml:"update_interval_minutes"`
	MinGlucose            float64 `toml:"min_glucose"`
	MaxGlucose            float64 `toml:"max_glucose"`
}

type ThemeConfig struct {
	Bg       string `toml:"bg"`
	Fg       string `toml:"fg"`
	Accent1  string `toml:"accent1"`
	Accent2  string `toml:"accent2"`
	Good     string `toml:"good"`
	Warning  string `toml:"warning"`
	Critical string `toml:"critical"`
}

type Config struct {
	App   AppConfig   `toml:"app"`
	Theme ThemeConfig `toml:"theme"`
}

const defaultTOML = `[app]
email = "tu@email.com"
password = "tu_password"
region = "cl" # Opciones: cl (Chile/LatAm), eu (Europa), us (EEUU), global
update_interval_minutes = 5
min_glucose = 70.0
max_glucose = 180.0

[theme]
# Tokyo Night
bg = "#1a1b26"
fg = "#c0caf5"
accent1 = "#7aa2f7"
accent2 = "#bb9af7"
good = "#9ece6a"
warning = "#e0af68"
critical = "#f7768e"
`

func Load() (Config, error) {
	var cfg Config
	home, err := os.UserHomeDir()
	if err != nil {
		return cfg, fmt.Errorf("no se pudo obtener el directorio home: %w", err)
	}

	configDir := filepath.Join(home, ".config", "libre-cli")
	configPath := filepath.Join(configDir, "config.toml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			return cfg, fmt.Errorf("error creando directorio de config: %w", err)
		}
		err = os.WriteFile(configPath, []byte(defaultTOML), 0644)
		if err != nil {
			return cfg, fmt.Errorf("error creando config.toml por defecto: %w", err)
		}
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, fmt.Errorf("error leyendo config.toml: %w", err)
	}
	err = toml.Unmarshal(data, &cfg)
	return cfg, err
}
