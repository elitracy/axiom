package subsystems

import (
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type Hvac struct {
	subsystem
	targetTemp utils.Unit
}

func NewHvac(id SubsystemID, name string, targetTemp utils.Unit) *Hvac {
	hvac := &Hvac{
		subsystem:  newSubsystem(id, name, utils.Machine),
		targetTemp: targetTemp,
	}

	hvac.AddComponent("temp", components.Temperature, targetTemp)
	hvac.AddComponent("temp-in", components.Temperature, targetTemp)
	hvac.AddComponent("power-in", components.Power, 0)

	hvac.thermalResponses["temp-regulation"] = utils.NewThermalResponse(10, 0.01)

	hvac.AddPorts("socket", 5, "power-in", PortInput)
	hvac.AddPorts("valve", 5, "temp-in", PortInput)

	return hvac
}

func (s *Hvac) Tick() {
	currentTemp := s.components["temp"].Value()

	tempIn := s.PortValue("temp-in")
	powerIn := s.PortValue("power-in")

	net := max(0, tempIn-powerIn)

	regulationDelta := s.thermalResponses["temp-regulation"].Delta(currentTemp, s.targetTemp)

	s.components["temp"].AddValue(net + regulationDelta)

	for _, port := range s.InputPorts() {
		port.Clear()
	}
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
