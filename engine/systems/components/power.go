package components

import "github.com/elias/axiom/engine/utils"

const (
	min_power = 0.0
	max_power = 1.0
)

type Power struct {
	*ComponentCore
}

// Creates a new Power component. Power can range from 0.0 (offline) to 1.0 (fully operational).
// initial is the starting power value
func NewPowerComponent(initial float64) *Power {
	initial = utils.Clamp(min_power, initial, max_power)

	return &Power{
		ComponentCore: NewComponentCore(
			"Power",
			initial,
			min_power,
			max_power,
		),
	}
}
