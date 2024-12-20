package resolv

// Resolv is a struct of properties which are translates to a local resolv for
// dnsmasq to redirect a given domain suffix to the local docker daemon.
// Windows has a custom solution, however this will be used on both Mac
// and Linux.
type Resolv struct {
	Data    string `yaml:"contents"`
	Enabled bool   `yaml:"enable"`
	File    string `yaml:"file"`
	Folder  string `yaml:"folder"`
	Name    string `yaml:"name"`
}
