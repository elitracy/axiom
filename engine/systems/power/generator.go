package systems

import (
	"fmt"

	"github.com/elias/axiom/engine/systems"
	"github.com/elias/axiom/engine/systems/components"
	"github.com/elias/axiom/engine/utils"
)

// TODO: create generator input
const (
	ticksTillDeath = 20

	maxGeneratorTemperature = 500.0

	startingPower  = 1.0
	startingFuel   = 1.0
	startingHealth = 1.0

	percentFuelUsedPerTick          = 1.0 / ticksTillDeath
	percentTemperatureRaisedPerTick = maxGeneratorTemperature / ticksTillDeath
	percentHealthLostPerTick        = 1.0 / ticksTillDeath
)

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
		normalized := (x - ambientTemperature) / (maxGeneratorTemperature - ambientTemperature)
		return (utils.Tanh(normalized, 3.1, 0)+0.015)*(maxGeneratorTemperature-ambientTemperature) + ambientTemperature
	}

	system := &Generator{
		SystemCore:  systems.NewSystemCore("Generator"),
		power:       components.NewComponent("Power", startingPower, 0.0, 1.0),
		fuel:        components.NewComponent("Fuel", startingFuel, 0.0, 1.0),
		temperature: components.NewComponent("Temperature", ambientTemperature, ambientTemperature, maxGeneratorTemperature, temperatureCurve),
		health:      components.NewHealthComponent(startingHealth),
	}

	return system

}

func (s *Generator) Tick() {
	if s.fuel.ApplyValueCurve() <= s.fuel.Min() {
		s.power.SetValue(s.power.Min())
	}

	if s.temperature.ApplyValueCurve() >= s.temperature.Max() {
		s.power.SetValue(s.power.Min())
	}

	if s.health.Status() == systems.Offline {
		s.power.SetValue(s.power.Min())
	}

	if s.power.ApplyValueCurve() != s.power.Min() {
		s.fuel.SetValue(s.fuel.Value() - percentFuelUsedPerTick)
		s.temperature.SetValue(s.temperature.Value() + percentTemperatureRaisedPerTick)
	} else {
		s.temperature.SetValue(s.temperature.Value() - percentTemperatureRaisedPerTick)
	}

	s.health.SetValue(s.health.Value() - percentHealthLostPerTick)
}

func (s *Generator) String() string {

	output := fmt.Sprintf("%v: %v", s.ID(), s.Name())
	output += fmt.Sprintf("\n%s", s.power)
	output += fmt.Sprintf("\n%s", s.fuel)
	output += fmt.Sprintf("\n%s", s.temperature)
	output += fmt.Sprintf("\n%s", s.health)

	return output
}
