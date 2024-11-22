package configuration

type WhitelistGroup []string

var (
	React_Whitelist   WhitelistGroup = []string{".html", ".css", ".js", ".jsx", ".ts", ".tsx"}
	ReactTs_Whitelist WhitelistGroup = []string{".html", ".css", ".ts", ".tsx"}
	ReactJs_Whitelist WhitelistGroup = []string{".html", ".css", ".js", ".jsx"}
)

func (cfg *ConfigOptions) AddGroup(group string) error {
	return nil
}

func (cfg *ConfigOptions) AddSingle(group string) error {
	return nil
}

func (cfg *ConfigOptions) RemoveGroup(group string) error {
	return nil
}

func (cfg *ConfigOptions) RemoveSingle(group string) error {
	return nil
}
