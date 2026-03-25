package subsystems

import (
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type Power struct {
	*subsystemCore
}

func NewPower(initPower utils.Norm) *Power {
	power := &Power{
		subsystemCore: newSubsystemCore("Power"),
	}

	power.AddComponent(components.Effort, initPower)

	return power
}

func (s *Power) Effort() utils.Norm { return s.components[components.Effort].Value() }

func (s *Power) Tick(inputs map[components.ComponentType]*components.Component) {
	cooling := utils.Norm(0.0)
	if coolant, ok := inputs[components.Temperature]; ok {
		cooling = coolant.Value()
	}

	delta := calcPowerTempDelta(s.components[components.Temperature].Value(), cooling, coolingCoef, s.components[components.Effort].Value())

	s.components[components.Temperature].AddValue(delta)
}

func calcPowerTempDelta(currentTemp, coolantTemp, coolingRate, heatRate utils.Norm) utils.Norm {
	cooling := (currentTemp - coolantTemp) * coolingRate

	return min(heatRate+cooling, heatRate)
}
