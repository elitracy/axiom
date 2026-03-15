package power

import (
	"github.com/elias/axiom/engine/materials"
	"github.com/elias/axiom/engine/systems"
	"github.com/elias/axiom/engine/systems/components"
)

// Creates a new power generator system.
// ambientTemperature is the starting temperature of the generator
func NewGenerator(housingMaterial materials.Metal, ambientTemperature float64) Power {

	system := &PowerCore{
		SystemCore:      systems.NewSystemCore("Generator"),
		power:           components.NewComponent("Power (%)", startingPower, 0.0, 1.0),
		fuel:            components.NewComponent("Fuel (%)", startingFuel, 0.0, 1.0),
		temperature:     components.NewComponent("Temperature (C)", housingMaterial.MinTemperature, ambientTemperature, housingMaterial.MaxTemperature, housingMaterial.TemperatureCurve),
		health:          components.NewHealthComponent(startingHealth),
		housingMaterial: housingMaterial,
	}

	return system
}
