package runtime

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

	SetField(name string, value interface{}) error
	GetFieldString(field string) (string, error)
	GetFieldInt(field string) (int, error)
	GetFieldBool(field string) (bool, error)
}
