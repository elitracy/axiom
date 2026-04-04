package subsystems

import (
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type hvacTickState struct {
	power utils.Unit
	heat  utils.Unit
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
	hvac.AddComponent("temp", components.Temperature, hvac.targetTemp)

	hvac.onInput("heat", func(comp components.Component) { hvac.state.heat = comp.Value() })
	hvac.onInput("power", func(comp components.Component) { hvac.state.power = comp.Value() })

	hvac.profiles["temp-regulation"] = utils.NewThermalResponse(10, 0.01)

	return hvac
}

func (s *Hvac) Effort() utils.Unit { return s.components["temp"].Value() }

func (s *Hvac) Tick(inputs map[string]components.Component) {
	s.dispatchInputs(inputs)

	currentTemp := s.components["temp"].Value()

	net := max(0, s.state.heat-s.state.power)
	regulationDelta := s.profiles["temp-regulation"].Delta(currentTemp, s.targetTemp)

	s.components["temp"].AddValue(net + regulationDelta)
}
