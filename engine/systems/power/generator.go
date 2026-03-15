package systems

import (
	"fmt"

	"github.com/elias/axiom/engine/systems"
	"github.com/elias/axiom/engine/systems/components"
	"github.com/elias/axiom/engine/utils"
)

// TODO: create generator input
const (
	startingPower       = 1.0
	startingFuel        = 1.0
	startingHealth      = 1.0
	startingTemperature = 0.0

	maxTemperature = 500.0

	percentFuelLostPerTick          = 1.0 / systems.TICKS_TILL_DEATH_DEBUG
	percentTemperatureGainedPerTick = 1.0 / systems.TICKS_TILL_DEATH_DEBUG
	percentHealthLostPerTick        = 1.0 / systems.TICKS_TILL_DEATH_DEBUG
)

// A system that creates power from fuel but generates a lot of heat
type Generator struct {
	*systems.SystemCore
	power       components.Component
	fuel        components.Component
	temperature components.Component
	health      *components.Health
}

// Creates a Generator system.
// ambientTemperature is the starting temperature of the generator
func NewGenerator(ambientTemperature float64) *Generator {

	temperatureCurve := func(x float64) float64 {
		temp := (utils.Tanh(x, 3.1, 0)+0.015)*(maxTemperature-ambientTemperature) + ambientTemperature
		return utils.Clamp(ambientTemperature, temp, maxTemperature)
	}

	system := &Generator{
		SystemCore:  systems.NewSystemCore("Generator"),
		power:       components.NewComponent("Power (%%)", startingPower),
		fuel:        components.NewComponent("Fuel (%%)", startingFuel),
		temperature: components.NewComponent("Temperature (C)", startingTemperature, temperatureCurve),
		health:      components.NewHealthComponent(startingHealth),
	}

	return system
}

// Updates the component values for the generator
func (s *Generator) Tick() {
	if s.fuel.Value() <= s.fuel.Min() {
		s.power.SetValue(s.power.Min())
	}

	if s.temperature.Value() >= s.temperature.Max() {
		s.power.SetValue(s.power.Min())
	}

	if s.health.Status() == systems.Offline {
		s.power.SetValue(s.power.Min())
	}

	if s.power.Value() != s.power.Min() {
		s.fuel.SetValue(s.fuel.Value() - percentFuelLostPerTick)
		s.temperature.SetValue(s.temperature.Value() + percentTemperatureGainedPerTick)
	} else {
		s.temperature.SetValue(s.temperature.Value() - percentTemperatureGainedPerTick)
	}

	s.health.SetValue(s.health.Value() - percentHealthLostPerTick)
}

// Returns the stringified information for the generator
func (s *Generator) String() string {

	output := fmt.Sprintf("%v: %v", s.ID(), s.Name())
	output += fmt.Sprintf("\n%s", s.power)
	output += fmt.Sprintf("\n%s", s.fuel)
	output += fmt.Sprintf("\n%s", s.temperature)
	output += fmt.Sprintf("\n%s", s.health)

	return output
}
