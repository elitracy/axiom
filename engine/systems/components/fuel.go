package components

import "github.com/elias/axiom/engine/utils"

const (
	min_Fuel = 0.0
	max_Fuel = 1.0
)

type Fuel struct {
	*ComponentCore
}

// Creates a new Fuel component. Fuel ranges from 0.0 to 1.0
// initial is the starting fuel level.
func NewFuelComponent(initial float64) *Fuel {
	initial = utils.Clamp(min_Fuel, initial, max_Fuel)

	return &Fuel{
		ComponentCore: NewComponentCore(
			"Fuel",
			initial,
			min_Fuel,
			max_Fuel,
		),
	}
}
