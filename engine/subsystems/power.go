package subsystems

import (
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type powerTickState struct {
	coolantRate utils.Unit
	coolantTemp utils.Unit
}

type Power struct {
	*subsystemCore
	state powerTickState
}

func NewPower(initPower utils.Unit) *Power {
	power := &Power{
		subsystemCore: newSubsystemCore("Power"),
	}

	power.AddComponent("power", components.Power, initPower)
	power.AddComponent("temp", components.Temperature, 0)

	power.onInput("cooling", func(comp components.Component) { power.state.coolantTemp = comp.Value() })
	power.onInput("cooling-rate", func(comp components.Component) { power.state.coolantRate = comp.Value() })

	power.profiles["cooling"] = utils.NewThermalResponse(10, .05)
	power.profiles["heating"] = utils.NewThermalResponse(10, .05)

	return power
}

func (s *Power) Effort() utils.Unit { return s.components["power"].Value() }

func (s *Power) Tick(inputs map[string]components.Component) {
	s.dispatchInputs(inputs)

	currentTemp := s.components["temp"].Value()

	heatingDelta := s.profiles["heating"].Delta(currentTemp, s.components["power"].Value())
	coolingDelta := s.profiles["cooling"].Delta(currentTemp, s.state.coolantTemp) * s.state.coolantRate

	s.components["temp"].AddValue(heatingDelta + coolingDelta)
}
