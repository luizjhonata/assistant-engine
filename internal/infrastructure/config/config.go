package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	dirName         = ".assistant-engine"
	configFileName  = "config.json"
	DefaultTime     = "09:00"
	DefaultDelayH   = 24
)

type WebhookType string

const (
	WebhookWorkflow WebhookType = "workflow"
	WebhookClassic  WebhookType = "classic"
)

type Config struct {
	WebhookURL    string      `json:"webhook_url"`
	WebhookType   WebhookType `json:"webhook_type"`
	MentionID     string      `json:"mention_id"`
	MentionName   string      `json:"mention_name"`
	DefaultTime   string      `json:"default_time"`
	DefaultDelayH int         `json:"default_delay_hours"`
}

func Load() (Config, error) {
	path, err := configPath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("reading config file %s: %w", path, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parsing config file: %w", err)
	}

	if cfg.DefaultTime == "" {
		cfg.DefaultTime = DefaultTime
	}
	if cfg.DefaultDelayH == 0 {
		cfg.DefaultDelayH = DefaultDelayH
	}
	if cfg.WebhookType == "" {
		cfg.WebhookType = WebhookWorkflow
	}

	return cfg, nil
}

func DataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolving home directory: %w", err)
	}

	dir := filepath.Join(home, dirName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("creating data directory: %w", err)
	}

	return dir, nil
}

func configPath() (string, error) {
	dir, err := DataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, configFileName), nil
}
