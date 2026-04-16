package subsystems

import (
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type Hvac struct {
	subsystem
}

func NewHvac(id SubsystemID, name string, targetTemp utils.Unit) *Hvac {
	hvac := &Hvac{
		subsystem: newSubsystem(id, name, utils.Machine),
	}

	hvac.AddComponent("temp", components.Temperature, targetTemp)
	hvac.AddComponent("target-temp", components.Temperature, targetTemp)

	hvac.thermalResponses["temp-regulation"] = utils.NewThermalResponse(10, 0.01)

	hvac.AddPorts("socket", 5, "power-in", PortInput)
	hvac.AddPorts("valve", 5, "temp-in", PortInput)

	return hvac
}

func (s *Hvac) Tick() {
	currentTemp := s.components["temp"].Value()

	tempIn, _ := s.InputSum("temp-in")
	powerIn, _ := s.InputSum("power-in")

	net := max(0, tempIn-powerIn)

	regulationDelta := s.thermalResponses["temp-regulation"].Delta(currentTemp, s.Components()["target-temp"].Value())

	s.components["temp"].AddValue(net + regulationDelta)
}
