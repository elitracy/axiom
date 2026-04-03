package subsystems

import (
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type powerTickState struct {
	coolantFlow utils.Norm
	coolantTemp utils.Norm
}

type Power struct {
	*subsystemCore
	state powerTickState
}

func NewPower(initPower utils.Norm) *Power {
	power := &Power{
		subsystemCore: newSubsystemCore("Power"),
	}

	power.AddComponent("power", components.Power, initPower)
	power.AddComponent("temp", components.Temperature, 0)

	power.onInput("temp-out", func(comp components.Component) {
		power.state.coolantTemp += comp.Value()
		power.state.coolantTemp = power.state.coolantTemp.Clamp()
	})

	power.onInput("flow-out", func(comp components.Component) {
		power.state.coolantFlow += comp.Value()
		power.state.coolantFlow = power.state.coolantTemp.Clamp()
	})

	return power
}

func (s *Power) Effort() utils.Norm { return s.components["power"].Value() }

func (s *Power) Tick(inputs map[string]components.Component) {
	s.dispatchInputs(inputs)

	delta := calcPowerTempDelta(
		s.components["temp"].Value(),
		s.state.coolantTemp,
		s.state.coolantFlow,
		s.components["power"].Value(),
	)
	logging.Info("DELTA: %.2f", delta)

	s.components["temp"].AddValue(delta)
}

func calcPowerTempDelta(currentTemp, coolantTemp, coolingRate, heatRate utils.Norm) utils.Norm {
	cooling := (currentTemp + coolantTemp) * -coolingRate

	return min(heatRate+cooling, heatRate)
}
