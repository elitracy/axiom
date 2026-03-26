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

	cooling.AddComponent(components.Flow, initEffort)
	cooling.AddComponent(components.Temperature, 0.5)

	return cooling
}

func (s *Cooling) Effort() utils.Norm { return s.components[components.Temperature].Value() }

func (s *Cooling) Tick(inputs map[components.ComponentType][]*components.Component) {}
