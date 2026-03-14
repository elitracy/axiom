package systems

import (
	"fmt"

	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/systems"
	"github.com/elias/axiom/engine/systems/components"
)

// TODO: create generator input
const (
	minGeneratorTemperature = 0.0
	maxGeneratorTemperature = 500.0
	startingPower           = 1.0

	startingFuel   = 1.0
	ticksTillDeath = 20

	percentFuelUsedPerTick          = 1.0 / ticksTillDeath
	percentTemperatureRaisedPerTick = maxGeneratorTemperature / ticksTillDeath
	percentHealthLostPerTick        = 1.0 / ticksTillDeath

	startingHealth = 1.0
)

type Generator struct {
	*systems.SystemCore
	power       *components.Power
	fuel        *components.Fuel
	temperature *components.Thermal
	health      *components.Health
}

// Creates a Generator system.
// ambientTemperature is the starting temperature of the generator
func NewGenerator(ambientTemperature float64) *Generator {
	system := &Generator{
		SystemCore:  systems.NewSystemCore("Generator"),
		power:       components.NewPowerComponent(startingPower),
		fuel:        components.NewFuelComponent(startingFuel),
		temperature: components.NewThermalComponent(ambientTemperature, minGeneratorTemperature, maxGeneratorTemperature),
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
		logging.Info("temperature change: %v", s.temperature.Value())

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
