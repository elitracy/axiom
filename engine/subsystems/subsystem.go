package subsystems

import (
	"fmt"
	"slices"
	"strings"

	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/telemetry"
	"github.com/elias/axiom/engine/utils"
)

type SubsystemID int64

type subsystem struct {
	id                 SubsystemID
	name               utils.SubsystemName
	subsystemType      utils.SubsystemType
	components         map[string]*components.Component
	thermalResponses   map[string]utils.ThermalResponse
	inputPorts         map[string]*InputPort
	outputPorts        map[string]*OutputPort
	currentPortID      PortID
	currentComponentID components.ComponentID
}

func newSubsystem(id SubsystemID, name utils.SubsystemName, subsystemType utils.SubsystemType) subsystem {
	return subsystem{
		id:               id,
		name:             name,
		subsystemType:    subsystemType,
		components:       make(map[string]*components.Component),
		thermalResponses: make(map[string]utils.ThermalResponse),
		inputPorts:       make(map[string]*InputPort),
		outputPorts:      make(map[string]*OutputPort),
	}
}

func (s *subsystem) newPortID() PortID {
	id := s.currentPortID
	s.currentPortID++
	return id
}

func (s *subsystem) newComponentID() components.ComponentID {
	id := s.currentComponentID
	s.currentComponentID++
	return id
}

func (s subsystem) ID() SubsystemID           { return s.id }
func (s subsystem) Name() utils.SubsystemName { return s.name }
func (s subsystem) Type() utils.SubsystemType { return s.subsystemType }

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

func (s *subsystem) addPort(name string, component string, kind utils.PortType) error {

	id := s.newPortID()

	switch kind {
	case utils.PortInput:
		name = "in." + name

		if _, exists := s.inputPorts[name]; exists {
			return fmt.Errorf("Could not add port, input port %v already exists on %v", name, s.Name())
		}

		port := newInputPort(id, name, s.Components()[component])
		s.inputPorts[name] = port
	case utils.PortOutput:
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

func (s *subsystem) AddPorts(prefix string, count int, component string, kind utils.PortType) {
	for i := range count {
		name := fmt.Sprintf("%s-%d", prefix, i)
		s.addPort(name, component, kind)
	}
}

func (s *subsystem) String() string {

	output := fmt.Sprintf("%s[%d]", s.name, s.id)

	for _, comp := range s.components {
		output += fmt.Sprintf("\n%s[%s]: %.2f", comp.Name(), comp.Type(), comp.Value())
	}

	return output
}

func (s *subsystem) ExportFields() *telemetry.Export {
	export := telemetry.NewExport()
	export.Add("name", string(s.Name()))

	var comps []*components.Component
	for _, comp := range s.components {
		comps = append(comps, comp)
	}

	slices.SortFunc(comps, func(a, b *components.Component) int {
		return strings.Compare(a.Name(), b.Name())
	})

	for _, comp := range comps {
		export.Add(comp.Name(), fmt.Sprintf("%.2f", comp.Value()))
	}

	return export

}
