package ioc

type Bean interface {
	Name() string
}

type IOC interface {
	Apply(bean any)
	Set(...any)
	Get(any) any
	ApplyAll()
}
