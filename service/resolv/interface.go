package resolv

type resolv interface {
	Clean()
	Configure()
	New() Resolv
	Status() bool
}
