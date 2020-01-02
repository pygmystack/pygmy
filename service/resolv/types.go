package resolv

type Resolv struct {
	Data     string `yaml:"contents"`
	Disabled bool   `yaml:"disabled"`
	File     string `yaml:"file"`
	Folder   string `yaml:"folder"`
	Name     string `yaml:"name"`
}
