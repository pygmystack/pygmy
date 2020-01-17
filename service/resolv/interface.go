package resolv

type resolv interface {
	Clean()
	Configure()
	Status() bool
}
