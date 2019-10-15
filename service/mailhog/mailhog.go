package mailhog

import (
	model "github.com/fubarhouse/pygmy/service/interface"
)

func New() model.Service {
	return model.Service{
		Name:          "Mailhog",
		Address:       "127.0.0.1",
		ContainerName: "mailhog.docker.amazee.io",
		Domain:        "docker.amazee.io",
		ImageName:     "mailhog/mailhog",
		RunCmd: []string{
			"run",
			"--restart=always",
			"-d",
			"-p",
			"1025:1025",
			"--expose",
			"80",
			"-u",
			"0",
			"--name=mailhog.docker.amazee.io",
			"-e \"MH_UI_BIND_ADDR=0.0.0.0:80\"",
			"-e \"MH_API_BIND_ADDR=0.0.0.0:80\"",
			"-e \"AMAZEEIO=AMAZEEIO\"",
			"mailhog/mailhog",
		},
	}
}
