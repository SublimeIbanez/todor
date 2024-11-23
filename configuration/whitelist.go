package configuration

import "slices"

func (cfg *ConfigOptions) AddToWhitelist(group string) error {
	if slices.Contains(cfg.Whitelist, group) {
		return nil
	}

	cfg.Whitelist = append(cfg.Whitelist, group)
	if err := cfg.saveConfig(); err != nil {
		return err
	}

	return nil
}

func (cfg *ConfigOptions) RemoveFromWhitelist(group string) error {
	return nil
}
