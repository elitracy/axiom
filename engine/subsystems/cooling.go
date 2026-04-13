package subsystems

import (
	"fmt"

	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type Cooling struct {
	*subsystemCore
}

func NewCooling(name string, initEffort utils.Unit) *Cooling {

	cooling := &Cooling{
		subsystemCore: newSubsystemCore(name),
	}

	cooling.AddComponent("temp-out", components.Temperature, 0.5)

	for i := range 5 {
		cooling.AddPort(fmt.Sprintf("valve-%d", i), "temp-out", PortOutput)
	}

	return cooling
}

func (s *Cooling) Effort() utils.Unit { return 1 }

func (s *Cooling) Tick() {}
