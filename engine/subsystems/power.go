package subsystems

import (
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type powerTickState struct {
	coolantFlow utils.Unit
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
	power.onInput("cooling-rate", func(comp components.Component) { power.state.coolantFlow = comp.Value() })

	return power
}

func (s *Power) Effort() utils.Unit { return s.components["power"].Value() }

func (s *Power) Tick(inputs map[string]components.Component) {
	s.dispatchInputs(inputs)

	delta := calcPowerTempDelta(
		s.components["temp"].Value(),
		s.state.coolantTemp,
		s.state.coolantFlow,
		s.components["power"].Value(),
	)

	s.components["temp"].AddValue(delta)
}

func calcPowerTempDelta(currentTemp, coolantTemp, coolingRate, heatRate utils.Unit) utils.Unit {
	cooling := (currentTemp + coolantTemp) * -coolingRate

	return min(heatRate+cooling, heatRate)
}
