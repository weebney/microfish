package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type tomlconfig struct {
	UseDarkMode bool `toml:"useDarkMode"`
}

func handleConfig() {
	config, err := loadOrCreateConfig()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if !config.UseDarkMode {
		return
	}
}

func loadOrCreateConfig() (tomlconfig, error) {
	config := tomlconfig{}
	configDir, err := os.UserConfigDir()
	if err != nil {
		return config, err
	}

	configFile := filepath.Join(filepath.Join(configDir, "microfish"), "config.toml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(configDir, "microfish"), 0o755); err != nil {
			return config, err
		}
		if err := copyEmbeddedConfig(configFile); err != nil {
			return config, err
		}
		configFile = "./defaultconfig.toml"
	} else if err != nil {
		return config, err
	}

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		return config, err
	}

	return config, nil
}

func copyEmbeddedConfig(destination string) error {
	dest, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dest.Close()

	defaultConfig, err := os.ReadFile("./defaultconfig.toml")
	if err != nil {
		return err
	}

	_, err = io.WriteString(dest, string(defaultConfig))
	return err
}
