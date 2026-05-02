package subsystems

import (
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type Power struct {
	subsystem
}

func NewPower(id SubsystemID, name utils.SubsystemName, initPower utils.Unit) *Power {
	power := &Power{
		subsystem: newSubsystem(id, name, utils.Power),
	}

	power.AddComponent("power-out", components.Power, initPower)
	power.AddComponent("temp-out", components.Temperature, 0)
	power.AddComponent("temp-in", components.Temperature, 0)

	power.thermalResponses["cooling"] = utils.NewThermalResponse(10, .05)
	power.thermalResponses["heating"] = utils.NewThermalResponse(10, .05)

	power.AddPorts("socket", 5, "power-out", utils.PortOutput)
	power.AddPorts("valve", 5, "temp-out", utils.PortOutput)
	power.AddPorts("valve", 5, "temp-in", utils.PortInput)

	return power
}

func (s *Power) Tick() {

	currentTemp := s.components["temp-out"]
	tempIn := s.components["temp-in"]

	heatingDelta := s.thermalResponses["heating"].Delta(currentTemp.Value(), s.components["power-out"].Value())

	var coolingDelta utils.Unit
	if tempIn.HasValue() {
		coolingDelta = s.thermalResponses["cooling"].Delta(currentTemp.Value(), tempIn.Value())
	}

	s.components["temp-out"].AddValue(heatingDelta + coolingDelta)
}

func (s *Power) Status() utils.Status {
	temp := s.components["temp-out"].Value()
	switch {
	case temp < .3:
		return utils.Healthy
	case temp < .6:
		return utils.Warning
	case temp < .8:
		return utils.Critical
	default:
		return utils.Offline
	}
}
