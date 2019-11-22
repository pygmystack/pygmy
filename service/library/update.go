package library

import "github.com/fubarhouse/pygmy/service/amazee"

func Update(c Config) {
	amazee.AmazeeImagePull()
}
