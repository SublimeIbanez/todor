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

type DefaultIgnoreOptions string
type UseGitIgnoreOptions int

const (
	Git         DefaultIgnoreOptions = ".git"
	GitIgnore   DefaultIgnoreOptions = ".gitignore"
	NodeModules DefaultIgnoreOptions = "node_modules"
	Next        DefaultIgnoreOptions = ".next"
	Png         DefaultIgnoreOptions = ".png"
	Bmp         DefaultIgnoreOptions = ".bmp"
	Jpeg        DefaultIgnoreOptions = ".jpeg"
	Svc         DefaultIgnoreOptions = ".svc"
	Eps         DefaultIgnoreOptions = ".eps"

	// Fucking go I swear this is your biggest hindrance
	InvalidUseGitIgnore UseGitIgnoreOptions = 0
	UseGitIgnore        UseGitIgnoreOptions = 1
	DoNotUseGitIgnore   UseGitIgnoreOptions = 2
)

type ConfigOptions struct {
	Ignore           []string
	UseGitIgnore     UseGitIgnoreOptions
	DefaultOutputDir string
}

// Sets defaults for outdated configuration files and saves them if necessary
func (cfg *ConfigOptions) setDefaults() {
	must_save := false

	if cfg.Ignore == nil {
		must_save = true
		cfg.Ignore = []string{string(Git), string(GitIgnore), string(Next), string(NodeModules), string(Bmp), string(Jpeg), string(Svc), string(Png), string(Eps)}
	}
	if cfg.DefaultOutputDir == "" {
		must_save = true
		cfg.DefaultOutputDir = "."
	}
	if cfg.UseGitIgnore == InvalidUseGitIgnore {
		must_save = true
		cfg.UseGitIgnore = UseGitIgnore
	}

	if must_save {
		fmt.Println("Updating config file to new version")
		cfg.saveConfig()
	}
}

func DefaultConfig() ConfigOptions {
	defaultConfig := ConfigOptions{
		Ignore:           []string{string(Git), string(GitIgnore), string(Next), string(NodeModules), string(Bmp), string(Jpeg), string(Svc), string(Png), string(Eps)},
		UseGitIgnore:     UseGitIgnore,
		DefaultOutputDir: ".",
	}

	return defaultConfig
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

	cfg.setDefaults()

	return cfg, nil
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
