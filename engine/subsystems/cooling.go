package subsystems

import (

	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type Cooling struct {
	subsystem
}

func NewCooling(id SubsystemID, name string, initTemp utils.Unit) *Cooling {

	cooling := &Cooling{
		subsystem: newSubsystem(id, name),
	}

	cooling.AddComponent("temp-out", components.Temperature, initTemp)
		cooling.AddPorts("valve", 5, "temp-out", PortOutput)

	return cooling
}

func (s *Cooling) Tick() {}
