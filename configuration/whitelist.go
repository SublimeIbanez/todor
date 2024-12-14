package configuration

import "slices"

func (cfg *ConfigOptions) AddItemToWhitelist(item string) error {
	if slices.Contains(cfg.Whitelist, item) {
		return nil
	}

	cfg.Whitelist = append(cfg.Whitelist, item)

	return cfg.saveConfig()
}

func (cfg *ConfigOptions) RemoveFromWhitelist(item string) error {
	if !slices.Contains(cfg.Whitelist, item) {
		return nil
	}

	index := slices.Index(cfg.Whitelist, item)
	cfg.Whitelist = append(cfg.Whitelist[:index], cfg.Whitelist[index+1:]...)

	return cfg.saveConfig()
}

func (cfg *ConfigOptions) ResetWhitelist() error {
	cfg.Whitelist = []string{}

	return cfg.saveConfig()
}
