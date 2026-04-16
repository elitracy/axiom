package subsystems

import (
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type Power struct {
	subsystem
}

func NewPower(id SubsystemID, name string, initPower utils.Unit) *Power {
	power := &Power{
		subsystem: newSubsystem(id, name, utils.Power),
	}

	power.AddComponent("power-out", components.Power, initPower)
	power.AddComponent("temp-out", components.Temperature, 0)

	power.thermalResponses["cooling"] = utils.NewThermalResponse(10, .05)
	power.thermalResponses["heating"] = utils.NewThermalResponse(10, .05)

	power.AddPorts("socket", 5, "power-out", PortOutput)
	power.AddPorts("valve", 5, "temp-out", PortOutput)
	power.AddPorts("valve", 5, "temp-in", PortInput)

	return power
}

func (s *Power) Tick() {

	currentTemp := s.components["temp-out"].Value()

	heatingDelta := s.thermalResponses["heating"].Delta(currentTemp, s.components["power-out"].Value())

	var coolingDelta utils.Unit
	if tempIn, ok := s.InputSum("temp-in"); ok {
		coolingDelta = s.thermalResponses["cooling"].Delta(currentTemp, tempIn)
	}

	s.components["temp-out"].AddValue(heatingDelta + coolingDelta)
}
