package subsystems

import (
	"math"

	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

const (
	hvacHeatingRate = 0.5
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
	hvac.AddComponent("temp-ambient", components.Temperature, hvac.targetTemp)

	return hvac
}

func (s *Hvac) Effort() utils.Norm { return s.components["temp-ambient"].Value() }

func (s *Hvac) Tick(inputs map[components.ComponentType][]components.Component) {
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
	heatDelta := calcHvacHeatDelta(s.targetTemp, heat, hvacHeatingRate)

	s.components["temp-ambient"].AddValue(heatDelta)

	currentTemp := s.components["temp-ambient"].Value()
	diff := utils.Norm(math.Abs(float64(s.targetTemp - currentTemp)))
	coolDelta := min(power, diff)

	switch {
	case currentTemp < s.targetTemp:
		s.components["temp-ambient"].AddValue(coolDelta)
	case currentTemp > s.targetTemp:
		s.components["temp-ambient"].AddValue(-coolDelta)
	}

}

func calcHvacHeatDelta(targetTemp, heat, rate utils.Norm) utils.Norm {
	heating := (heat - targetTemp) * rate
	return heating
}
