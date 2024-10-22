package runtimes

// TODO.

// ServiceRuntime is ...
type ServiceRuntime interface {
	Setup() error
	Start() error
	Create() error
	Status() (bool, error)
	Labels() (map[string]string, error)
	ID() (string, error)
	Clean() error
	Stop() error
	StopAndRemove() error
	Remove() error
}
