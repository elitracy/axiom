package materials

import "github.com/elias/axiom/engine/utils"

type Metal struct {
	Name                string                // the name of the metal
	MinTemperature      float64               // the minimum temperature of the metal (C)
	MaxTemperature      float64               // the minimum temperature of the metal (C)
	MinPressure         float64               // the minimum pressure of the metal (kPa)
	MaxPressure         float64               // the maximum pressure of the metal (kPa)
	HeatAbsorptionRate  float64               // how much heat is absorbed by the metal (C/tick)
	MaxTemperatureDelta float64               // the maximum temperature delta per tick (C/tick)
	TemperatureCurve    func(float64) float64 // the curve to calculate the applied temperature from the normalized value
	PressureCurve       func(float64) float64 // the curve to calculate the applied pressure from the normalized value
}

func NewSteel() Metal {
	metal := Metal{
		Name:                "Steel",
		MinTemperature:      -29.0,
		MaxTemperature:      427.0,
		MinPressure:         1.0,
		MaxPressure:         800.0,
		MaxTemperatureDelta: 0.05,
		HeatAbsorptionRate:  0.05,
	}

	metal.TemperatureCurve = func(x float64) float64 {
		return utils.Tanh(x, 3.1, 0)
	}

	metal.PressureCurve = func(x float64) float64 {
		return utils.Tanh(x, 1.1, 0.3)
	}

	return metal
}
