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

type inputHandler func(comp components.Component)

type Subsystem interface {
	ID() SubsystemID
	Name() string
	Effort() utils.Unit
	Components() map[string]*components.Component
	AddComponent(string, components.ComponentType, utils.Unit)
	String() string

	Tick(inputs map[string]components.Component)
	onInput(name string, handler inputHandler)
	dispatchInputs(inputs map[string]components.Component)
}

type subsystemCore struct {
	Subsystem
	id            SubsystemID
	name          string
	components    map[string]*components.Component
	inputHandlers map[string]inputHandler
}

func newSubsystemCore(name string) *subsystemCore {
	return &subsystemCore{
		id:            newID(),
		name:          name,
		components:    make(map[string]*components.Component),
		inputHandlers: make(map[string]inputHandler),
	}
}

func (s *subsystemCore) ID() SubsystemID { return s.id }
func (s *subsystemCore) Name() string    { return s.name }
func (s *subsystemCore) Components() map[string]*components.Component {
	return s.components
}

func (s *subsystemCore) AddComponent(name string, componentType components.ComponentType, value utils.Unit) {
	component := components.NewComponent(name, componentType, value)
	s.components[component.Name()] = component
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

func (s subsystemCore) dispatchInputs(inputs map[string]components.Component) {
	for name, comp := range inputs {
		if handler, exists := s.inputHandlers[name]; exists {
			handler(comp)
		}
	}
}
