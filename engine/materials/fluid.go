package materials

import "github.com/elias/axiom/engine/utils"

type Fluid struct {
	Name                    string                // the name of the fluid
	MinTemperature          float64               // the minimum temperature of the fluid (C)
	MaxTemperature          float64               // the minimum temperature of the fluid (C)
	MinViscosity            float64               // the minimum viscosity of the fluid (lower = less flow) (cP)
	MaxViscosity            float64               // the maximum viscosity of the fluid (higher = more flow) (cP)
	HeatAbsorptionRate      float64               // how much heat is absorbed by the fluid (C/tick)
	ThermalConductivityRate float64               // the maximum temperature delta (C/tick)
	ThermalExpansionRate    float64               // the rate at which the fluid changes pressure
	TemperatureCurve        func(float64) float64 // the curve to calculate the applied temperature from the normalized value
	ViscosityCurve          func(float64) float64 // the curve to calculate the applied viscosity from the normalized value
}

func NewWater() Fluid {
	fluid := Fluid{
		Name:                    "Water",
		MinTemperature:          0.0,
		MaxTemperature:          100.0,
		MinViscosity:            0.3,
		MaxViscosity:            1.0,
		HeatAbsorptionRate:      0.02,
		ThermalConductivityRate: 0.02,
		ThermalExpansionRate:    0.8,
	}

	fluid.TemperatureCurve = func(x float64) float64 {
		temp := (utils.Tanh(x, 3.1, 0)+0.015)*(fluid.MaxTemperature-fluid.MinTemperature) + fluid.MinTemperature
		return utils.Clamp(fluid.MinTemperature, temp, fluid.MaxTemperature)
	}

	fluid.ViscosityCurve = func(x float64) float64 {
		viscosity := (1-utils.Tanh(x, 4.2, 0))*(fluid.MaxViscosity-fluid.MinViscosity) + fluid.MinViscosity
		return utils.Clamp(fluid.MinViscosity, viscosity, fluid.MaxViscosity)
	}

	return fluid

}

func NewPropyleneGlycol() Fluid {
	fluid := Fluid{
		Name:                    "Propylene-Glycol",
		MinTemperature:          -50.0,
		MaxTemperature:          180.0,
		MinViscosity:            1.0,
		MaxViscosity:            50.0,
		HeatAbsorptionRate:      0.01,
		ThermalConductivityRate: 0.02,
		ThermalExpansionRate:    0.4,
	}

	fluid.TemperatureCurve = func(x float64) float64 {
		temp := (utils.Tanh(x, 3.1, 0)+0.015)*(fluid.MaxTemperature-fluid.MinTemperature) + fluid.MinTemperature
		return utils.Clamp(fluid.MinTemperature, temp, fluid.MaxTemperature)
	}

	fluid.ViscosityCurve = func(x float64) float64 {
		viscosity := (1-utils.Tanh(x, 10.0, 0))*(fluid.MaxViscosity-fluid.MinViscosity) + fluid.MinViscosity
		return utils.Clamp(fluid.MinViscosity, viscosity, fluid.MaxViscosity)
	}

	return fluid

}
