package subsystems

import (
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type Cooling struct {
	*subsystemCore
}

func NewCooling(initEffort utils.Norm) *Cooling {
	cooling := &Cooling{
		subsystemCore: newSubsystemCore("Cooling"),
	}

	cooling.AddComponent(components.Effort, initEffort)

	return cooling
}

func (s *Cooling) Effort() utils.Norm { return s.components[components.Effort].Value() }

func (s *Cooling) Tick(inputs map[components.ComponentType]*components.Component) {}
