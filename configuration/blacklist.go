package configuration

import "slices"

func (cfg *ConfigOptions) AddToBlacklist(item string) error {
	if slices.Contains(cfg.Blacklist, item) {
		return nil
	}

	cfg.Blacklist = append(cfg.Blacklist, item)

	return cfg.saveConfig()
}

func (cfg *ConfigOptions) RemoveFromBlacklist(item string) error {
	if !slices.Contains(cfg.Blacklist, item) {
		return nil
	}

	index := slices.Index(cfg.Blacklist, item)
	cfg.Blacklist = append(cfg.Blacklist[:index], cfg.Blacklist[index+1:]...)

	return cfg.saveConfig()
}

func (cfg *ConfigOptions) ResetBlacklist() error {
	cfg.Blacklist = []string{}

	return cfg.saveConfig()
}
