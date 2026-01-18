package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HistoryFile       string   `yaml:"history_file"`
	MaxHistory        int      `yaml:"max_history"`
	DefaultSearchPath string   `yaml:"default_search_path"`
	IgnoreDirs        []string `yaml:"ignore_dirs"`
}

func defaultConfig() Config {
	home, _ := os.UserHomeDir()
	return Config{
		HistoryFile:       filepath.Join(home, ".zsh_history"),
		MaxHistory:        50,
		DefaultSearchPath: ".",
		IgnoreDirs:        []string{"node_modules", ".git", "vendor"},
	}
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".sakibox", "config.yaml"), nil
}

func Load() (Config, error) {
	path, err := configPath()
	if err != nil {
		return Config{}, err
	}

	cfg := defaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return Config{}, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func EnsureConfig() error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); err == nil {
		return nil
	}

	cfg := defaultConfig()
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
