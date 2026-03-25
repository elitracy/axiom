package subsystems

import (
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

const (
	coolingCoef utils.Norm = -0.5
)

var currentSubsystemID = 0

func newID() SubsystemID {
	id := currentSubsystemID
	currentSubsystemID++
	return SubsystemID(id)
}

type SubsystemID int64

type Subsystem interface {
	ID() SubsystemID
	Name() string
	Effort() utils.Norm
	Components() map[components.ComponentType]*components.Component
	AddComponent(components.ComponentType, utils.Norm)

	Tick(inputs map[components.ComponentType]*components.Component)
}

type subsystemCore struct {
	id         SubsystemID
	name       string
	components map[components.ComponentType]*components.Component
}

func newSubsystemCore(name string) *subsystemCore {
	return &subsystemCore{
		id:         newID(),
		name:       name,
		components: make(map[components.ComponentType]*components.Component),
	}
}

func (s *subsystemCore) ID() SubsystemID { return s.id }
func (s *subsystemCore) Name() string    { return s.name }
func (s *subsystemCore) Components() map[components.ComponentType]*components.Component {
	return s.components
}

func (s *subsystemCore) AddComponent(componentType components.ComponentType, value utils.Norm) {
	component := components.NewComponent(componentType, value)
	s.components[componentType] = component
}
