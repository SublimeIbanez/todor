package common

type Group []string

var (
	React   Group = []string{".html", ".css", ".js", ".jsx", ".ts", ".tsx"}
	ReactTs Group = []string{".html", ".css", ".ts", ".tsx"}
	ReactJs Group = []string{".html", ".css", ".js", ".jsx"}
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
