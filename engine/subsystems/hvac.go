package subsystems

import (
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type Hvac struct {
	*subsystemCore
	targetTemp utils.Unit
}

func NewHvac() *Hvac {
	hvac := &Hvac{
		subsystemCore: newSubsystemCore("HVAC"),
		targetTemp:    utils.Unit(.2),
	}
	hvac.AddComponent("temp-in", components.Temperature, hvac.targetTemp)
	hvac.AddComponent("power-in", components.Power, 0)

	hvac.AddComponent("temp", components.Temperature, hvac.targetTemp)

	hvac.onInput("temp-in", func(port *InputPort) { hvac.components["temp-in"].SetValue(port.Input()) })
	hvac.onInput("power-in", func(port *InputPort) { hvac.components["power-in"].SetValue(port.Input()) })

	hvac.profiles["temp-regulation"] = utils.NewThermalResponse(10, 0.01)

	return hvac
}

func (s *Hvac) Effort() utils.Unit { return s.components["temp"].Value() }

func (s *Hvac) Tick() {
	s.dispatchInputs()

	currentTemp := s.components["temp"].Value()
	inputTemp := s.components["temp-in"].Value()
	effort := s.components["power-in"].Value()

	net := max(0, inputTemp-effort)
	regulationDelta := s.profiles["temp-regulation"].Delta(currentTemp, s.targetTemp)

	s.components["temp"].AddValue(net + regulationDelta)
}
