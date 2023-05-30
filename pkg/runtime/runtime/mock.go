// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package runtime

import "storj.io/storj-up/pkg/recipe"

// MockService is a service implementation for MockRuntime.
type MockService struct {
	Identifier  ServiceInstance
	Image       string
	Persisted   []string
	Config      map[string]string
	Environment map[string]string
	Flag        []string
	Ports       map[int]int
	Label       []string
}

// GetENV implements runtime.Service.
func (m *MockService) GetENV() map[string]*string {
	// TODO implement me
	panic("implement me")
}

// GetVolumes implements runtime.Service.
func (m *MockService) GetVolumes() []VolumeMount {
	// TODO implement me
	panic("implement me")
}

// UseFolder implements runtime.Service.
func (m *MockService) UseFolder(path string, name string) error {
	return nil
}

// UseFile implements runtime.Service.
func (m *MockService) UseFile(path string, name string, data string) error {
	return nil
}

// NewMockService can create a new mock service.
func NewMockService(name string) *MockService {
	return &MockService{
		Identifier: ServiceInstance{
			Name: name,
		},
		Config:      map[string]string{},
		Environment: map[string]string{},
		Flag:        []string{},
		Persisted:   []string{},
	}
}

// Labels implements runtime.Service.
func (m *MockService) Labels() []string {
	return m.Label
}

// ID implements runtime.Service.
func (m *MockService) ID() ServiceInstance {
	return m.Identifier
}

// ChangeImage implements runtime.Service.
func (m *MockService) ChangeImage(f func(string) string) error {
	m.Image = f(m.Image)
	return nil
}

// AddConfig implements runtime.Service.
func (m *MockService) AddConfig(key string, value string) error {
	m.Config[key] = value
	return nil
}

// AddFlag implements runtime.Service.
func (m *MockService) AddFlag(flag string) error {
	m.Flag = append(m.Flag, flag)
	return nil
}

// RemoveFlag implements runtime.Service.
func (m *MockService) RemoveFlag(flag string) error {
	panic("implement me")
}

// AddEnvironment implements runtime.Service.
func (m *MockService) AddEnvironment(key string, value string) error {
	m.Environment[key] = value
	return nil
}

// AddPortForward implements runtime.Service.
func (m *MockService) AddPortForward(portMap PortMap) error {
	m.Ports[portMap.External] = portMap.Internal
	return nil
}

// RemovePortForward implements runtime.Service.
func (m *MockService) RemovePortForward(portMap PortMap) error {
	delete(m.Ports, portMap.External)
	return nil
}

// Persist  implements runtime.Service.
func (m *MockService) Persist(dir string) error {
	m.Persisted = append(m.Persisted, dir)
	return nil
}

var _ Service = &MockService{}

// MockRuntime is a runtime `implementation` for testing.
type MockRuntime struct {
	Services []Service
}

var _ Runtime = &MockRuntime{}

// NewMockRuntime creates a runtime for testing.
func NewMockRuntime() *MockRuntime {
	return &MockRuntime{
		Services: []Service{},
	}
}

// GetHost implements runtime.Runtime.
func (m *MockRuntime) GetHost(serviceInstance ServiceInstance, hostType string) string {
	panic("implement me")
}

// GetPort implements runtime.Runtime.
func (m *MockRuntime) GetPort(serviceInstance ServiceInstance, portType string) PortMap {
	panic("implement me")
}

// Get implements runtime.Runtime.
func (m *MockRuntime) Get(serviceInstance ServiceInstance, name string) string {
	panic("implement me")
}

// AddService implements runtime.Runtime.
func (m *MockRuntime) AddService(service recipe.Service) (Service, error) {
	s := NewMockService(service.Name)
	err := InitFromRecipe(s, service)
	if err != nil {
		return s, err
	}
	m.Services = append(m.Services, s)
	return s, nil
}

// Write implements runtime.Runtime.
func (m *MockRuntime) Write() error {
	return nil
}

// GetServices implements runtime.Runtime.
func (m *MockRuntime) GetServices() []Service {
	return m.Services
}

// Reload implements runtime.Runtime.
func (m *MockRuntime) Reload(stack recipe.Stack) error {
	return nil
}
