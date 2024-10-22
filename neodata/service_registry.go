package neodata

import "sync"

type ServiceRegistry struct {
	services sync.Map
}

func (sr *ServiceRegistry) Register(name string, service interface{}) {
	sr.services.Store(name, service)
}

func (sr *ServiceRegistry) Get(name string) (interface{}, bool) {
	return sr.services.Load(name)
}

// A typed helper for casting services
func GetService[T any](sr *ServiceRegistry, name string) (*T, bool) {
	service, ok := sr.Get(name)
	if !ok {
		return nil, false
	}
	casted, ok := service.(*T)
	return casted, ok
}
