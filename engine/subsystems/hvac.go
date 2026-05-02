package subsystems

import (
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type Hvac struct {
	subsystem
	targetTemp utils.Unit
}

func NewHvac(id SubsystemID, name utils.SubsystemName, targetTemp utils.Unit) *Hvac {
	hvac := &Hvac{
		subsystem:  newSubsystem(id, name, utils.Hvac),
		targetTemp: targetTemp,
	}

	hvac.thermalResponses["temp-regulation"] = utils.NewThermalResponse(10, 0.01)

	hvac.AddComponent("temp", components.Temperature, targetTemp)
	hvac.AddComponent("temp-in", components.Temperature, targetTemp)
	hvac.AddComponent("power-in", components.Power, 0)

	hvac.AddPorts("socket", 5, "power-in", utils.PortInput)
	hvac.AddPorts("valve", 5, "temp-in", utils.PortInput)

	return hvac
}

func (s *Hvac) Tick() {
	currentTemp := s.components["temp"].Value()

	tempIn := s.components["temp-in"].Value()
	powerIn := s.components["power-in"].Value()

	net := max(0, tempIn-powerIn)

	regulationDelta := s.thermalResponses["temp-regulation"].Delta(currentTemp, s.targetTemp)

	s.components["temp"].AddValue(net + regulationDelta)
}

func (s *Hvac) Status() utils.Status {
	temp := s.components["temp"].Value()
	switch {
	case temp < .15:
		return utils.Warning
	case temp < .25:
		return utils.Healthy
	case temp < .35:
		return utils.Warning
	case temp < .4:
		return utils.Critical
	default:
		return utils.Offline
	}
}
