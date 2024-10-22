package service

import "github.com/pygmystack/pygmy/internal/runtime"

type Service interface {
	New(c *runtime.Params) runtime.Service
	NewDefaultPorts() runtime.Service
}
