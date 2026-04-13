package subsystems

import (
	"fmt"

	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type Hvac struct {
	*subsystemCore
}

func NewHvac(name string, targetTemp utils.Unit) *Hvac {
	hvac := &Hvac{
		subsystemCore: newSubsystemCore(name),
	}

	hvac.AddComponent("temp", components.Temperature, targetTemp)
	hvac.AddComponent("target-temp", components.Temperature, targetTemp)

	hvac.onInput("temp-in", hvac.accumulateInput("temp-in", components.Temperature))
	hvac.onInput("power-in", hvac.accumulateInput("power-in", components.Power))

	hvac.profiles["temp-regulation"] = utils.NewThermalResponse(10, 0.01)

	for i := range 5 {
		hvac.AddPort(fmt.Sprintf("socket-%d", i), "power-in", PortInput)
	}

	for i := range 5 {
		hvac.AddPort(fmt.Sprintf("valve-%d", i), "temp-in", PortInput)
	}

	return hvac
}

func (s *Hvac) Effort() utils.Unit { return s.components["temp"].Value() }

func (s *Hvac) Tick() {
	s.dispatchInputs()

	currentTemp := s.components["temp"].Value()

	var net, tempVal, powerVal utils.Unit
	if temp, exists := s.inputComponents["temp-in"]; exists {
		tempVal = temp.Value()
	}
	if power, exists := s.inputComponents["power-in"]; exists {
		powerVal = power.Value()
	}

	net = max(0, tempVal-powerVal)

	regulationDelta := s.profiles["temp-regulation"].Delta(currentTemp, s.Components()["target-temp"].Value())

	s.components["temp"].AddValue(net + regulationDelta)

	for key := range s.inputComponents {
		delete(s.inputComponents, key)
	}
}
