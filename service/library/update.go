package library

import "github.com/fubarhouse/pygmy/v1/service/amazee"

func Update(c Config) {
	amazee.AmazeeImagePull()
}
