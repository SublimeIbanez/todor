package configuration

import "slices"

func (cfg *ConfigOptions) AddToBlacklist(item string) error {
	if slices.Contains(cfg.Blacklist, item) {
		return nil
	}

	cfg.Blacklist = append(cfg.Blacklist, item)
	if err := cfg.saveConfig(); err != nil {
		return err
	}

	return nil
}

func (cfg *ConfigOptions) RemoveFromBlacklist(item string) error {
	if !slices.Contains(cfg.Blacklist, item) {
		return nil
	}

	index := slices.Index(cfg.Blacklist, item)
	cfg.Blacklist = append(cfg.Blacklist[:index], cfg.Blacklist[index+1:]...)

	if err := cfg.saveConfig(); err != nil {
		return err
	}

	return nil
}
