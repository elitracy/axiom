package subsystems

import (
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type Hvac struct {
	*subsystemCore
	targetTemp utils.Norm
}

func NewHvac() *Hvac {
	hvac := &Hvac{
		subsystemCore: newSubsystemCore("HVAC"),
		targetTemp:    utils.Norm(.2),
	}
	hvac.AddComponent(components.Temperature, hvac.targetTemp)

	return hvac
}

func (s *Hvac) Effort() utils.Norm { return s.components[components.Effort].Value() }

func (s *Hvac) Tick(inputs map[components.ComponentType][]*components.Component) {
	power := utils.Norm(0.0)
	if powers, ok := inputs[components.Power]; ok {
		for _, p := range powers {
			power += p.Value()
		}

	}

	heat := utils.Norm(0.0)
	if temps, ok := inputs[components.Temperature]; ok {
		for _, t := range temps {
			heat += t.Value()
		}
	}

	currentTemp := s.components[components.Temperature].Value()
	delta := calcHvacHeatDelta(s.targetTemp, heat, power)
	logging.Info("DELTA: %2.f", delta)

	if currentTemp != s.targetTemp {
		s.components[components.Temperature].AddValue(delta)
	}

}

func calcHvacHeatDelta(targetTemp, heat, rate utils.Norm) utils.Norm {
	cooling := (heat - targetTemp) * -rate
	return cooling
}
