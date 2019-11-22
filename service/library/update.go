package library

import "github.com/fubarhouse/pygmy-go/service/amazee"

func Update(c Config) {
	amazee.AmazeeImagePull()
}
