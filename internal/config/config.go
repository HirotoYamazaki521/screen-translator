package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = "config.json"

// Region はキャプチャ範囲を表す。null の場合は全画面キャプチャ。
type Region struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Config はアプリケーション設定を表す。
type Config struct {
	DisplayIndex    int     `json:"display_index"`
	IntervalSeconds int     `json:"interval_seconds"`
	CaptureRegion   *Region `json:"capture_region"`
	LanguageFrom    string  `json:"language_from"`
	LanguageTo      string  `json:"language_to"`
	FontSize        int     `json:"font_size"`
	APIKey          string  `json:"api_key"`
}

// Default はデフォルト設定を返す。
func Default() Config {
	return Config{
		DisplayIndex:    0,
		IntervalSeconds: 5,
		CaptureRegion:   nil,
		LanguageFrom:    "en",
		LanguageTo:      "ja",
		FontSize:        18,
	}
}

// configDir は設定ファイルのディレクトリパスを返す。
func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".screen-translator"), nil
}

// Load は設定ファイルを読み込む。ファイルが存在しない場合はデフォルト値を返す。
func Load() (Config, error) {
	dir, err := configDir()
	if err != nil {
		return Default(), err
	}

	path := filepath.Join(dir, configFileName)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Default(), nil
	}
	if err != nil {
		return Default(), fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Default(), fmt.Errorf("failed to parse config file: %w", err)
	}
	return cfg, nil
}

// Save は設定ファイルに書き込む。ディレクトリが存在しない場合は作成する。
func Save(cfg Config) error {
	dir, err := configDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}

	path := filepath.Join(dir, configFileName)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}
