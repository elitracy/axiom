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

	power.AddComponent(components.Power, initPower)
	power.AddComponent(components.Temperature, 0)
	return power
}

func (s *Power) Effort() utils.Norm { return s.components[components.Power].Value() }

func (s *Power) Tick(inputs map[components.ComponentType][]*components.Component) {
	coolantTemp := utils.Norm(0.0)
	if temps, ok := inputs[components.Temperature]; ok {
		for _, t := range temps {
			coolantTemp += t.Value()
		}
	}

	coolantFlow := utils.Norm(0.0)
	if flows, ok := inputs[components.Flow]; ok {
		for _, f := range flows {
			coolantFlow += f.Value()
		}
	}

	delta := calcPowerTempDelta(s.components[components.Temperature].Value(), coolantTemp, coolantFlow, s.components[components.Power].Value())

	s.components[components.Temperature].AddValue(delta)
}

func calcPowerTempDelta(currentTemp, coolantTemp, coolingRate, heatRate utils.Norm) utils.Norm {
	cooling := (currentTemp + coolantTemp) * -coolingRate

	return min(heatRate+cooling, heatRate)
}
