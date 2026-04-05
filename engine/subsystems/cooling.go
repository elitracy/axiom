package subsystems

import (
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type Cooling struct {
	*subsystemCore
}

func NewCooling(initEffort utils.Unit) *Cooling {
	cooling := &Cooling{
		subsystemCore: newSubsystemCore("Cooling"),
	}

	cooling.AddComponent("temp-out", components.Temperature, 0.5)

	return cooling
}

func (s *Cooling) Effort() utils.Unit { return 1 }

func (s *Cooling) Tick() {}
