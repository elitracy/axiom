package subsystems

import (
	"math"

	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

const (
	hvacHeatingRate = 0.5
)

type hvacTickState struct {
	power utils.Unit
}

type Hvac struct {
	*subsystemCore
	targetTemp utils.Unit
	state      hvacTickState
}

func NewHvac() *Hvac {
	hvac := &Hvac{
		subsystemCore: newSubsystemCore("HVAC"),
		targetTemp:    utils.Unit(.2),
	}
	hvac.AddComponent("ambient-temp", components.Temperature, hvac.targetTemp)

	hvac.onInput("heat", func(comp components.Component) {
		delta := calcHvacHeatDelta(hvac.targetTemp, comp.Value(), hvacHeatingRate)
		hvac.components["ambient-temp"].AddValue(delta)
	})

	hvac.onInput("power", func(comp components.Component) {
		hvac.state.power = comp.Value()
	})

	return hvac
}

func (s *Hvac) Effort() utils.Unit { return s.components["ambient-temp"].Value() }

func (s *Hvac) Tick(inputs map[string]components.Component) {
	s.dispatchInputs(inputs)

	currentTemp := s.components["ambient-temp"].Value()

	diff := utils.Unit(math.Abs(float64(s.targetTemp - currentTemp)))
	coolDelta := min(s.state.power, diff)

	switch {
	case currentTemp < s.targetTemp:
		s.components["ambient-temp"].AddValue(coolDelta)
	case currentTemp > s.targetTemp:
		s.components["ambient-temp"].AddValue(-coolDelta)
	}
}

func calcHvacHeatDelta(targetTemp, heat, rate utils.Unit) utils.Unit {
	heating := (heat - targetTemp) * rate
	return heating
}
