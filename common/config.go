package common

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

const CONFIG_FILE_NAME string = ".todor_cfg.json"

type ConfigOptions struct {
	Ignore           []string
	DefaultOutputDir string
}

func getConfigFilePath() (string, error) {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	home_dir_path, err := filepath.Abs(home_dir)
	if err != nil {
		return "", err
	}

	return filepath.Join(home_dir_path, CONFIG_FILE_NAME), nil
}

func LoadConfig() (ConfigOptions, error) {
	config_file_path, err := getConfigFilePath()
	if err != nil {
		return ConfigOptions{}, nil
	}

	data, err := os.ReadFile(config_file_path)
	if err != nil {
		if os.IsNotExist(err) {
			default_config := DefaultConfig()
			fmt.Printf("No configuration file found at %s, generating defaults", config_file_path)
			if err := default_config.saveConfig(); err != nil {
				log.Fatal("Could not save default_config: ", err)
			}
			return default_config, nil
		}
		return ConfigOptions{}, err
	}

	var cfg ConfigOptions
	if err = json.Unmarshal(data, &cfg); err != nil {
		return ConfigOptions{}, err
	}

	return cfg, nil
}

func DefaultConfig() ConfigOptions {
	defaultConfig := ConfigOptions{
		Ignore:           []string{".git", "node_modules", ".next"},
		DefaultOutputDir: ".",
	}

	return defaultConfig
}

func (config *ConfigOptions) saveConfig() error {
	config_file_path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(config_file_path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fs.FileMode(DEFAULT_FILE_PERMISSIONS))
	if err != nil {
		return err
	}
	defer file.Close()

	json_data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	_, err = file.Write(json_data)
	if err != nil {
		return err
	}

	return nil
}

func (cfg *ConfigOptions) SetIgnore(ignoreList []string) {
	cfg.Ignore = ignoreList
	err := cfg.saveConfig()
	if err != nil {
		log.Fatal("Unable to save changes to ignore list: ", err)
	}
}

func (cfg *ConfigOptions) AddIgnore(ignore_item string) {
	cfg.Ignore = append(cfg.Ignore, ignore_item)
	err := cfg.saveConfig()
	if err != nil {
		log.Fatal("Unable to save changes to ignore list: ", err)
	}
}
