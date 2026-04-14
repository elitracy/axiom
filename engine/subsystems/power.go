package subsystems

import (
	"fmt"

	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type Power struct {
	*subsystemCore
}

func NewPower(name string, initPower utils.Unit) *Power {
	power := &Power{
		subsystemCore: newSubsystemCore(name),
	}

	power.AddComponent("power-out", components.Power, initPower)
	power.AddComponent("temp-out", components.Temperature, 0)

	power.profiles["cooling"] = utils.NewThermalResponse(10, .05)
	power.profiles["heating"] = utils.NewThermalResponse(10, .05)

	for i := range 5 {
		power.AddPort(fmt.Sprintf("socket-%d", i), "power-out", PortOutput)
	}

	for i := range 5 {
		power.AddPort(fmt.Sprintf("valve-%d", i), "temp-out", PortOutput)
	}

	for i := range 5 {
		power.AddPort(fmt.Sprintf("valve-%d", i+5), "temp-in", PortInput)
	}

	return power
}

func (s *Power) Effort() utils.Unit { return s.components["power"].Value() }

func (s *Power) Tick() {

	currentTemp := s.components["temp-out"].Value()

	heatingDelta := s.profiles["heating"].Delta(currentTemp, s.components["power-out"].Value())

	var coolingDelta utils.Unit
	if tempIn, ok := s.InputSum("temp-in"); ok {
		coolingDelta = s.profiles["cooling"].Delta(currentTemp, tempIn)
	}

	s.components["temp-out"].AddValue(heatingDelta + coolingDelta)

	for key := range s.inputComponents {
		delete(s.inputComponents, key)
	}
}
