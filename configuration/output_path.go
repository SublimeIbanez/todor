package configuration

import "os"

func (cfg *ConfigOptions) SetOutputPath(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// TODO: create the path if it doesn't exist?
		return err
	}

	cfg.OutputDirectory = path
	return cfg.saveConfig()
}
