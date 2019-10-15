// +build windows

package ssh_addkey

import (
	"fmt"
	model "github.com/fubarhouse/pygmy/service/interface"
	"github.com/mitchellh/go-homedir"
	"os"
)

var (
	image_name = "amazeeio/ssh-agent"
	container_name = "amazeeio-ssh-agent-add-key"
)

func NewAdder(key string) model.Service {
	if key == "" {
		homedir, _ := homedir.Dir()
		key = fmt.Sprintf("%v%v.ssh%vid_rsa", homedir, string(os.PathSeparator), string(os.PathSeparator))
	}

	return model.Service{
		Name: "amazeeio-ssh-agent-add-key",
		Address: "",
		ContainerName: "amazeeio-ssh-agent-add-key",
		Domain: "",
		ImageName: "amazeeio/ssh-agent",
		RunCmd: []string{
			"run",
			"--rm",
			fmt.Sprintf("--volume=%v:/key", key),
			"--volumes-from=amazeeio-ssh-agent",
			"--name="+container_name,
			image_name,
			"windows-key-add",
			"/key",
		},
	}
}

func NewShower() model.Service {
	return model.Service{
		Address: "",
		ContainerName: "amazeeio-ssh-agent-add-key",
		Domain: "",
		ImageName: "amazeeio/ssh-agent",
		RunCmd: []string{
			"run",
			"--rm",
			"--volumes-from=amazeeio-ssh-agent",
			"--name="+container_name,
			image_name,
			"ssh-add",
			"-l",
		},
	}
}
