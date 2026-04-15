package subsystems

import (
	"fmt"

	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type SubsystemID int64

type subsystem struct {
	id                 SubsystemID
	name               string
	components         map[string]*components.Component
	thermalResponses   map[string]utils.ThermalResponse
	inputPorts         map[string]*InputPort
	outputPorts        map[string]*OutputPort
	currentPortID      PortID
	currentComponentID components.ComponentID
}

func newSubsystem(id SubsystemID, name string) subsystem {
	return subsystem{
		id:               id,
		name:             name,
		components:       make(map[string]*components.Component),
		thermalResponses: make(map[string]utils.ThermalResponse),
		inputPorts:       make(map[string]*InputPort),
		outputPorts:      make(map[string]*OutputPort),
	}
}

func (s *subsystem) newPortID() PortID {
	id := s.currentPortID
	s.currentPortID++
	return PortID(id)
}

func (s *subsystem) newComponentID() components.ComponentID {
	id := s.currentComponentID
	s.currentComponentID++

	return id
}

func (s subsystem) ID() SubsystemID { return s.id }
func (s subsystem) Name() string    { return s.name }
func (s subsystem) Components() map[string]*components.Component {
	return s.components
}
func (s subsystem) InputPorts() map[string]*InputPort   { return s.inputPorts }
func (s subsystem) OutputPorts() map[string]*OutputPort { return s.outputPorts }

func (s *subsystem) AddComponent(name string, componentType components.ComponentType, value utils.Unit) error {
	if _, exists := s.components[name]; exists {
		return fmt.Errorf("Could not add component, component %v already exists on %v", name, s.Name())
	}

	id := s.newComponentID()

	component := components.NewComponent(id, name, componentType, value)
	s.components[component.Name()] = component

	return nil
}

func (s *subsystem) AddPort(name string, component string, kind PortType) error {

	id := s.newPortID()

	switch kind {
	case PortInput:
		name = "in." + name

		if _, exists := s.inputPorts[name]; exists {
			return fmt.Errorf("Could not add port, input port %v already exists on %v", name, s.Name())
		}

		port := newInputPort(id, name, component)
		s.inputPorts[name] = port
	case PortOutput:
		name = "out." + name

		if _, exists := s.components[component]; !exists {
			return fmt.Errorf("Could not add output port, component %v doesn't exist on %v", component, s.Name())
		}
		if _, exists := s.outputPorts[name]; exists {
			return fmt.Errorf("Could not add port, output port %v already exists on %v", name, s.Name())
		}

		port := newOutputPort(id, name, s.Components()[component])
		s.outputPorts[name] = port
	}

	return nil
}

func (s *subsystem) AddPorts(prefix string, count int, component string, kind PortType) {
	for i := range count {
		name := fmt.Sprintf("%s-%d", prefix, i)
		s.AddPort(name, component, kind)
	}
}

func (s subsystem) String() string {

	output := fmt.Sprintf("%s[%d]", s.name, s.id)

	for _, comp := range s.components {
		output += fmt.Sprintf("\n%s[%s]: %.2f", comp.Name(), comp.Type(), comp.Value())
	}

	return output
}

func (s *subsystem) InputSum(channel string) (utils.Unit, bool) {
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
