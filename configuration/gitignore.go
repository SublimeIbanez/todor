package configuration

func (cfg *ConfigOptions) SetGitIgnore(value *bool) error {
	cfg.Gitignore = value
	return cfg.saveConfig()
}
