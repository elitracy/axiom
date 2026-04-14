package subsystems

import (
	"fmt"

	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type Hvac struct {
	*subsystemCore
}

func NewHvac(id SubsystemID, name string, targetTemp utils.Unit) *Hvac {
	hvac := &Hvac{
		subsystemCore: newSubsystemCore(id, name),
	}

	hvac.AddComponent("temp", components.Temperature, targetTemp)
	hvac.AddComponent("target-temp", components.Temperature, targetTemp)

	hvac.thermalResponses["temp-regulation"] = utils.NewThermalResponse(10, 0.01)

	for i := range 5 {
		hvac.AddPort(fmt.Sprintf("socket-%d", i), "power-in", PortInput)
	}

	for i := range 5 {
		hvac.AddPort(fmt.Sprintf("valve-%d", i), "temp-in", PortInput)
	}

	return hvac
}

func (s *Hvac) Tick() {
	currentTemp := s.components["temp"].Value()

	tempIn, _ := s.InputSum("temp-in")
	powerIn, _ := s.InputSum("power-in")

	net := max(0, tempIn-powerIn)
	logging.Debug("TEMP IN: %v", tempIn)
	logging.Debug("POWER IN: %v", powerIn)
	logging.Debug("NET: %v", net)

	regulationDelta := s.thermalResponses["temp-regulation"].Delta(currentTemp, s.Components()["target-temp"].Value())

	s.components["temp"].AddValue(net + regulationDelta)

	for key := range s.inputComponents {
		delete(s.inputComponents, key)
	}
}
