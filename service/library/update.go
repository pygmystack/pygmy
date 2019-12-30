package library

import "github.com/fubarhouse/pygmy-go/service/amazee"

// Update provides business logic for the `update` command.
func Update(c Config) {
	amazee.AmazeeImagePull()
}
