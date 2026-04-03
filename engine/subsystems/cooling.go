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

	cooling.AddComponent("flow-out", components.Flow, initEffort)
	cooling.AddComponent("temp-out", components.Temperature, 0.5)

	return cooling
}

func (s *Cooling) Effort() utils.Norm { return s.components["flow-out"].Value() }

func (s *Cooling) Tick(inputs map[string]components.Component) {}
