package sound

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Keyboard struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"keyboard"`
	Mouse struct {
		Left  string `json:"left"`
		Right string `json:"right"`
	} `json:"mouse"`
	Volume float64 `json:"volume"`
}

var (
	appDir          = ""
	left_click_url  = "https://akocdw82ai.ufs.sh/f/Jk6mQ2VBlE6tLk5IKluEC9coqerdXTUMpmgu6VvIWanSiKHh"
	right_click_url = "https://akocdw82ai.ufs.sh/f/Jk6mQ2VBlE6tZpCpGUgjRlfiMK6r04kEQNc59egnYLduJA3w"
)

func NewDefaultConfig() Config {
	return Config{
		Keyboard: struct {
			Type string `json:"type"`
			ID   string `json:"id"`
		}{
			Type: "mechavibes",
			ID:   "1200000000001",
		},
		Mouse: struct {
			Left  string `json:"left"`
			Right string `json:"right"`
		}{
			Left:  filepath.Join(appDir, "left_click.mp3"),
			Right: filepath.Join(appDir, "right_click.mp3"),
		},
		Volume: 0.5,
	}
}

func (c *Config) Save(filePath string) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	jsonData, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (c *Config) Load(filePath string) error {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if err := json.Unmarshal(fileData, c); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

func InitConfig(volume float64, soundPack, lcPath, rcPath string) (*Config, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache directory: %w", err)
	}

	appDir = filepath.Join(cacheDir, "clickclack")
	filePath := filepath.Join(appDir, "compose.json")

	config := NewDefaultConfig()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := config.Save(filePath); err != nil {
			return nil, fmt.Errorf("failed to create config file: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to check file: %w", err)
	} else {
		if err := config.Load(filePath); err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
	}

	downloadFile(left_click_url, config.Mouse.Left)
	downloadFile(right_click_url, config.Mouse.Right)

	update := false
	if volume != 0.0 {
		config.Volume = volume
		update = true
	}
	if soundPack != "" {
		config.Keyboard.ID = soundPack
		update = true
	}
	if lcPath != "" {
		config.Mouse.Left = lcPath
		update = true
	}
	if rcPath != "" {
		config.Mouse.Right = rcPath
		update = true
	}

	if update {
		if err := config.Save(filePath); err != nil {
			return nil, fmt.Errorf("failed to update config file: %w", err)
		}
	}

	return &config, nil
}
