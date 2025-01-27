package main

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
	appDir = ""

	configValue = Config{
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
			Left:  "path",
			Right: "path",
		},
		Volume: 0.5,
	}
)

func defaultInit() error {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	appDir = filepath.Join(cacheDir, "clickclack")
	err = os.MkdirAll(appDir, 0755)
	if err != nil {
		return err
	}

	filePath := filepath.Join(appDir, "compose.json")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {

		err = os.MkdirAll(filepath.Dir(filePath), 0755)
		if err != nil {
			return err
		}

		jsonData, err := json.MarshalIndent(configValue, "", "  ")
		if err != nil {
			return err
		}

		err = os.WriteFile(filePath, jsonData, 0644)
		if err != nil {
			return err
		}

	} else if err != nil {
		return err
	} else {
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		err = json.Unmarshal(fileData, &configValue)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return err
		}
	}

	fmt.Println("Config Path:", appDir)

	return nil
}
