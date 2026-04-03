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
	power utils.Norm
}

type Hvac struct {
	*subsystemCore
	targetTemp utils.Norm
	state      hvacTickState
}

func NewHvac() *Hvac {
	hvac := &Hvac{
		subsystemCore: newSubsystemCore("HVAC"),
		targetTemp:    utils.Norm(.2),
	}
	hvac.AddComponent("temp-ambient", components.Temperature, hvac.targetTemp)

	hvac.onInput("heat", func(comp components.Component) {
		delta := calcHvacHeatDelta(hvac.targetTemp, comp.Value(), hvacHeatingRate)
		hvac.components["temp-ambient"].AddValue(delta)
	})

	hvac.onInput("power", func(comp components.Component) {
		hvac.state.power = comp.Value()
	})

	return hvac
}

func (s *Hvac) Effort() utils.Norm { return s.components["temp-ambient"].Value() }

func (s *Hvac) Tick(inputs map[string]components.Component) {
	s.dispatchInputs(inputs)

	currentTemp := s.components["temp-ambient"].Value()

	diff := utils.Norm(math.Abs(float64(s.targetTemp - currentTemp)))
	coolDelta := min(s.state.power, diff)

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
