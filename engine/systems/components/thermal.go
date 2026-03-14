package components

import "github.com/elias/axiom/engine/utils"

type Thermal struct {
	*ComponentCore
}

// Creates a new thermal component
// initial is the initial temperature of the thermal component, min is the minimum temperature, and max is the maximum temperature.
func NewThermalComponent(initial, min, max float64) *Thermal {
	return &Thermal{
		ComponentCore: NewComponentCore("Thermal",
			initial,
			min,
			max,
			func(x float64) float64 {
				normalized := (x - min) / (max - min)
				return (utils.Tanh(normalized, 3.1, 0))*(max-min) + min
			},
		),
	}
}
