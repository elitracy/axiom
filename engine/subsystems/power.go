package subsystems

import (
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type Power struct {
	*subsystemCore
}

func NewPower(initPower utils.Unit) *Power {
	power := &Power{
		subsystemCore: newSubsystemCore("Power"),
	}

	power.AddComponent("temp-in", components.Temperature, 0)

	power.AddComponent("power-out", components.Power, initPower)
	power.AddComponent("temp-out", components.Temperature, 0)

	power.onInput("temp-in", func(comp components.Component) {

		power.components["temp-in"].SetValue(comp.Value())

	})

	power.profiles["cooling"] = utils.NewThermalResponse(10, .05)
	power.profiles["heating"] = utils.NewThermalResponse(10, .05)

	return power
}

func (s *Power) Effort() utils.Unit { return s.components["power"].Value() }

func (s *Power) Tick() {
	s.dispatchInputs()

	currentTemp := s.components["temp-out"].Value()

	heatingDelta := s.profiles["heating"].Delta(currentTemp, s.components["power-out"].Value())
	coolingDelta := s.profiles["cooling"].Delta(currentTemp, s.components["temp-in"].Value())

	s.components["temp-out"].AddValue(heatingDelta + coolingDelta)
}
