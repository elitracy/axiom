package subsystems

import (
	"fmt"

	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type SubsystemID int64

type Subsystem interface {
	AddPort(string, string, PortType) error
	String() string

	Tick()
}

type subsystemCore struct {
	id               SubsystemID
	name             string
	components       map[string]*components.Component
	inputComponents  map[string]*components.Component
	thermalResponses map[string]utils.ThermalResponse
	inputPorts       map[string]*InputPort
	outputPorts      map[string]*OutputPort
	currentPortID    PortID
}

func newSubsystemCore(id SubsystemID, name string) *subsystemCore {
	return &subsystemCore{
		id:               id,
		name:             name,
		components:       make(map[string]*components.Component),
		inputComponents:  make(map[string]*components.Component),
		thermalResponses: make(map[string]utils.ThermalResponse),
		inputPorts:       make(map[string]*InputPort),
		outputPorts:      make(map[string]*OutputPort),
	}
}

func (s *subsystemCore) newPortID() PortID {
	id := s.currentPortID
	s.currentPortID++
	return PortID(id)
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

func (s *subsystemCore) AddInputComponent(name string, componentType components.ComponentType, value utils.Unit) error {
	if _, exists := s.inputComponents[name]; exists {
		return fmt.Errorf("Could not add input component, input component %v already exists on %v", name, s.Name())
	}

	component := components.NewComponent(name, componentType, value)
	s.inputComponents[component.Name()] = component

	return nil
}

func (s *subsystemCore) AddPort(name string, component string, portType PortType) error {

	id := s.newPortID()

	switch portType {
	case PortInput:
		if _, exists := s.inputPorts[name]; exists {
			return fmt.Errorf("Could not add port, input port %v already exists on %v", name, s.Name())
		}

		port := newInputPort(id, name, s, component)
		s.inputPorts[name] = port
	case PortOutput:

		if _, exists := s.components[component]; !exists {
			return fmt.Errorf("Could not add output port, component %v doesn't exist on %v", component, s.Name())
		}
		if _, exists := s.outputPorts[name]; exists {
			return fmt.Errorf("Could not add port, output port %v already exists on %v", name, s.Name())
		}

		port := newOutputPort(id, name, s, s.Components()[component])
		s.outputPorts[name] = port
	}

	return nil
}

func (s subsystemCore) String() string {

	output := fmt.Sprintf("%s[%d]", s.name, s.id)

	for _, comp := range s.components {
		output += fmt.Sprintf("\n%s[%s]: %.2f", comp.Name(), comp.Type(), comp.Value())
	}

	return output
}

func (s *subsystemCore) InputSum(channel string) (utils.Unit, bool) {
	var sum utils.Unit
	got := false

	for _, p := range s.inputPorts {
		if p.channel == channel && p.received {
			got = true
			sum += p.value
			p.Clear()
		}
	}

	return sum, got
}
