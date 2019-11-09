package resolv

type resolv interface {
	Clean()
	Configure()
	New(Resolv) Resolv
	Status() bool
}
