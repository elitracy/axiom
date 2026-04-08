package subsystems

import (
	"fmt"

	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

var currentSubsystemID = 0

func newID() SubsystemID {
	id := currentSubsystemID
	currentSubsystemID++
	return SubsystemID(id)
}

type SubsystemID int64

type inputHandler func(port *InputPort)

type Subsystem interface {
	ID() SubsystemID
	Name() string
	Effort() utils.Unit
	Components() map[string]*components.Component
	InputPorts() map[string]*InputPort
	OutputPorts() map[string]*OutputPort
	AddComponent(string, components.ComponentType, utils.Unit) error
	AddPort(string, string, PortType) error
	String() string

	Tick()
	onInput(name string, handler inputHandler)
	dispatchInputs()
}

type subsystemCore struct {
	Subsystem
	id            SubsystemID
	name          string
	components    map[string]*components.Component
	inputHandlers map[string]inputHandler
	profiles      map[string]utils.ThermalResponse
	inputPorts    map[string]*InputPort
	outputPorts   map[string]*OutputPort
}

func newSubsystemCore(name string) *subsystemCore {
	return &subsystemCore{
		id:            newID(),
		name:          name,
		components:    make(map[string]*components.Component),
		inputHandlers: make(map[string]inputHandler),
		profiles:      make(map[string]utils.ThermalResponse),
		inputPorts:    make(map[string]*InputPort),
		outputPorts:   make(map[string]*OutputPort),
	}
}

func (s *subsystemCore) ID() SubsystemID { return s.id }
func (s *subsystemCore) Name() string    { return s.name }
func (s *subsystemCore) Components() map[string]*components.Component {
	return s.components
}
func (s *subsystemCore) InputPorts() map[string]*InputPort   { return s.inputPorts }
func (s *subsystemCore) OutputPorts() map[string]*OutputPort { return s.outputPorts }

func (s *subsystemCore) AddComponent(name string, componentType components.ComponentType, value utils.Unit) error {
	if _, exists := s.components[name]; exists {
		return fmt.Errorf("Could not add component, component %v already exists on %v", name, s.Name())
	}

	component := components.NewComponent(name, componentType, value)
	s.components[component.Name()] = component

	return nil
}

func (s *subsystemCore) AddPort(name string, componentName string, portType PortType) error {
	if _, exists := s.components[componentName]; !exists {
		return fmt.Errorf("Could not add port, component %v doesn't exist on %v", componentName, s.Name())
	}
	switch portType {
	case Input:
		if _, exists := s.inputPorts[name]; exists {
			return fmt.Errorf("Could not add port, input port %v already exists on %v", name, s.Name())
		}

		port := NewInputPort(name, s, s.Components()[componentName])
		s.inputPorts[name] = port

	case Output:
		if _, exists := s.outputPorts[name]; exists {
			return fmt.Errorf("Could not add port, output port %v already exists on %v", name, s.Name())
		}

		port := NewOutputPort(name, s, s.Components()[componentName])
		s.outputPorts[name] = port
	}

	return nil
}

func (s subsystemCore) String() string {

	output := fmt.Sprintf("%s[%d]", s.name, s.id)

	for _, comp := range s.components {
		output += fmt.Sprintf("\n%v[%s]: %.2f", comp.Type(), comp.Name(), comp.Value())
	}

	return output
}

func (s *subsystemCore) onInput(name string, handler inputHandler) {
	s.inputHandlers[name] = handler
}

func (s subsystemCore) dispatchInputs() {
	for _, port := range s.inputPorts {
		if handler, exists := s.inputHandlers[port.Component().Name()]; exists {
			handler(port)
		}
	}
}
