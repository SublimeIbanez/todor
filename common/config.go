package common

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const CONFIG_FILE_NAME string = ".todor_cfg.json"

type DefaultBlacklistFile string
type UseGitIgnoreOptions int

const (
	GitIgnore DefaultBlacklistFile = ".gitignore"

	// Fucking go I swear this is your biggest hindrance
	InvalidUseGitIgnore UseGitIgnoreOptions = 0
	UseGitIgnore        UseGitIgnoreOptions = 1
	DoNotUseGitIgnore   UseGitIgnoreOptions = 2
)

type ConfigOptions struct {
	Whitelist        []string
	UseGitIgnore     UseGitIgnoreOptions
	DefaultOutputDir string
}

// Sets defaults for outdated configuration files and saves them if necessary
func (cfg *ConfigOptions) setDefaults() {
	must_save := false

	if cfg.Whitelist == nil {
		must_save = true
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
		Whitelist:        []string{},
		UseGitIgnore:     UseGitIgnore,
		DefaultOutputDir: ".",
	}

	return defaultConfig
}

func getConfigFilePath() (string, error) {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	var config_dir string
	// GOOS operating systems | GOARCH architectures:
	// - AIX:                 `aix`       | ppc64,
	// - Android:             `android`   | 386, amd64, arm, arm64,
	// - MacOS:               `darwin`    | amd64, arm64,
	// - Dragonfly:           `dragonfly` | amd64,
	// - FreeBSD:             `freebsd`   | 386, amd64, arm,
	// - Illumos:             `illumos`   | amd64,
	// - iOS:                 `iOS`       | arm64,
	// - JSRE:                `js`        | wasm,
	// - Linux:               `linux`     | 386, amd64, arm, arm64, loong64, mips, mipsle, mips64, mips64le, ppc64, ppc64le, riscv64, s390x,
	// - NetBSD:              `netbsd`    | 386, amd64, arm,
	// - OpenBSD:             `openbsd`   | 386, amd64, arm, arm64,
	// - Plan9:               `plan9`     | 386, amd64, arm,
	// - Solaris:             `solaris`   | amd64,
	// - WASI Preview 1:      `wasip1`    | wasm,
	// - Windows:             `windows`   | 386, amd64, arm, arm64
	switch runtime.GOOS {
	case "windows":
		config_dir = filepath.Join(home_dir, "AppData", "Local", "todor")
		break
	case "linux":
		config_dir = filepath.Join(home_dir, ".config", "todor")
		break
	case "darwin": // MacOS
		config_dir = filepath.Join(home_dir, "Library", "Application Support", "todor")
	default:
		return "", fmt.Errorf("operating system <%s> not currently supported", runtime.GOOS)
	}

	if err := os.MkdirAll(config_dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return filepath.Join(config_dir, CONFIG_FILE_NAME), nil
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

func (cfg *ConfigOptions) SetWhiteList(whiteList []string) {
	cfg.Whitelist = whiteList
	err := cfg.saveConfig()
	if err != nil {
		log.Fatal("Unable to save changes to white list: ", err)
	}
}

func (cfg *ConfigOptions) AddIgnore(whitelist_item string) {
	cfg.Whitelist = append(cfg.Whitelist, whitelist_item)
	err := cfg.saveConfig()
	if err != nil {
		log.Fatal("Unable to save changes to ignore list: ", err)
	}
}
