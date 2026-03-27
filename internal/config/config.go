package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	User     UserConfig     `yaml:"user"`
	Settings SettingsConfig `yaml:"settings"`
	Wakeup   WakeupConfig   `yaml:"wakeup"`
	Layout   []LayoutEntry  `yaml:"layout"`
}

type UserConfig struct {
	Name     string         `yaml:"name"`
	Greeting string         `yaml:"greeting"`
	Extra    map[string]any `yaml:",inline"`
}

type SettingsConfig struct {
	MasterCurrency string `yaml:"master_currency"`
}

type WakeupConfig struct {
	Enabled     bool   `yaml:"enabled"`
	MinGapHours int    `yaml:"min_gap_hours"`
	CutoffHour  int    `yaml:"cutoff_hour"`
	Terminal    string `yaml:"terminal"`
}

type LayoutEntry struct {
	Module string         `yaml:"module"`
	Title  string         `yaml:"title"`
	Config map[string]any `yaml:"config"`
	Feeds  []any          `yaml:"feeds"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read config %q: %w", path, err)
	}

	cfg := &Config{}
	cfg.Wakeup.Enabled = true
	cfg.Wakeup.MinGapHours = 6
	cfg.Wakeup.CutoffHour = 15
	cfg.Wakeup.Terminal = "kitty"

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("cannot parse config %q: %w", path, err)
	}
	return cfg, nil
}
