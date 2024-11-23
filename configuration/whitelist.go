package configuration

import "slices"

func (cfg *ConfigOptions) AddToWhitelist(item string) error {
	if slices.Contains(cfg.Whitelist, item) {
		return nil
	}

	cfg.Whitelist = append(cfg.Whitelist, item)
	if err := cfg.saveConfig(); err != nil {
		return err
	}

	return nil
}

func (cfg *ConfigOptions) RemoveFromWhitelist(item string) error {
	if !slices.Contains(cfg.Whitelist, item) {
		return nil
	}

	index := slices.Index(cfg.Whitelist, item)
	cfg.Whitelist = append(cfg.Whitelist[:index], cfg.Whitelist[index+1:]...)
	if err := cfg.saveConfig(); err != nil {
		return err
	}

	return nil
}
