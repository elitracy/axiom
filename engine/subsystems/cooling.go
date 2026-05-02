package subsystems

import (
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type Cooling struct {
	subsystem
}

func NewCooling(id SubsystemID, name utils.SubsystemName, initTemp utils.Unit) *Cooling {

	cooling := &Cooling{
		subsystem: newSubsystem(id, name, utils.Cooling),
	}

	cooling.AddComponent("temp-out", components.Temperature, initTemp)
	cooling.AddPorts("valve", 5, "temp-out", utils.PortOutput)

	return cooling
}

func (s *Cooling) Tick() {}

func (s *Cooling) Status() utils.Status {
	temp := s.components["temp-out"].Value()
	switch {
	case temp < .4:
		return utils.Healthy
	case temp < .7:
		return utils.Warning
	case temp < 1:
		return utils.Critical
	default:
		return utils.Offline
	}
}
