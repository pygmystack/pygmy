package library

import "github.com/fubarhouse/pygmy-go/service/amazee"

// Update will update the Amazee images
func Update(c Config) {
	amazee.AmazeeImagePull()
}
