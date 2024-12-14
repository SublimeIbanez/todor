package configuration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (cfg *ConfigOptions) SetOutputDirectory(directory string) error {
	directory_status, err := os.Stat(directory)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// TODO: create the path if it doesn't exist?
		return err
	}
	if !directory_status.IsDir() {
		return fmt.Errorf("must provide a directory; file provided: %s", directory)
	}

	full_path, err := filepath.Abs(directory)
	if err != nil {
		return err
	}

	cfg.OutputDirectory = full_path
	return cfg.saveConfig()
}

func (cfg *ConfigOptions) SetOutputFilename(filename string) error {
	full_filename := filename
	if !strings.Contains(full_filename, ".") {
		full_filename += ".md"
	}

	cfg.OutputFilename = full_filename
	return cfg.saveConfig()
}

func (cfg *ConfigOptions) SetOutputPath(path string) error {

	return cfg.saveConfig()
}
