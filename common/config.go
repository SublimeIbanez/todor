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

var (
	CurrentVersion         string   = "0.0.1"
	DefaultWhitelist       []string = []string{}
	DefaultBlacklist       []string = []string{}
	DefaultUseGitignore    bool     = true
	DefaultOutputDirectory string   = "."
)

type ConfigOptions struct {
	// Zero-width to disallow positional initialization
	_ struct{}
	// Semantic versioning
	Version string `json:"version"`
	// List of all file/directory names and filetypes to parse
	Whitelist []string `json:"whitelist"`
	// List of all file/directory names and filetypes to not parse
	Blacklist []string `json:"blacklist"`
	// Uses the root directory's .gitignore as well as loaded blacklist
	Gitignore *bool `json:"gitignore"`
	// Output directory for the markdown file
	OutputDirectory string `json:"output_directory"`
}

// Generates a default configuration struct
func DefaultConfig() ConfigOptions {
	defaultConfig := ConfigOptions{
		Version:         CurrentVersion,
		Whitelist:       DefaultWhitelist,
		Blacklist:       DefaultBlacklist,
		Gitignore:       &DefaultUseGitignore,
		OutputDirectory: DefaultOutputDirectory,
	}

	return defaultConfig
}

// Sets defaults for outdated configuration files and saves them if necessary
func (cfg *ConfigOptions) validate() error {
	must_save := false

	if len(cfg.Version) == 0 || cfg.Version != CurrentVersion {
		must_save = true
		cfg.Version = CurrentVersion
	}
	if cfg.Whitelist == nil {
		must_save = true
		cfg.Whitelist = DefaultWhitelist
	}
	if cfg.Blacklist == nil {
		must_save = true
		cfg.Blacklist = DefaultBlacklist
	}
	if cfg.OutputDirectory == "" {
		must_save = true
		cfg.OutputDirectory = DefaultOutputDirectory
	}
	if cfg.Gitignore == nil {
		must_save = true
		cfg.Gitignore = &DefaultUseGitignore
	}

	if must_save {
		fmt.Println("Updating config file to new version")
		if err := cfg.saveConfig(); err != nil {
			return err
		}
	}

	return nil
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
		return ConfigOptions{}, fmt.Errorf("could not get configuration file path: %w", err)
	}

	data, err := os.ReadFile(config_file_path)
	if err != nil {
		// If the config doesn't exist
		if os.IsNotExist(err) {
			default_config := DefaultConfig()
			fmt.Printf("No configuration file found at %s, generating defaults", config_file_path)
			if err := default_config.saveConfig(); err != nil {
				log.Fatalf("Could not save default_config: %s", err.Error())
			}
			return default_config, nil
		}

		// All other options
		return ConfigOptions{}, fmt.Errorf("could not read configuration file: %w", err)
	}

	var cfg ConfigOptions
	if err = json.Unmarshal(data, &cfg); err != nil {
		return ConfigOptions{}, fmt.Errorf("could not unmarshal configuration file: %w", err)
	}

	if err = cfg.validate(); err != nil {
		return ConfigOptions{}, fmt.Errorf("unable to validate configuration file: %w", err)
	}

	return cfg, nil
}

func (config *ConfigOptions) saveConfig() error {
	config_file_path, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("unable to obtain configuration file path: %w", err)
	}

	file, err := os.OpenFile(config_file_path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fs.FileMode(DEFAULT_FILE_PERMISSIONS))
	if err != nil {
		return fmt.Errorf("unable to open the file: %w", err)
	}
	defer file.Close()

	json_data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return fmt.Errorf("unable to marshal data: %w", err)
	}

	_, err = file.Write(json_data)
	if err != nil {
		return fmt.Errorf("unable to write to file: %w", err)
	}

	return nil
}
