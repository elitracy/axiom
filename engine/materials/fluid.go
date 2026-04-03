package materials

type Fluid struct {
	Name string // the name of the fluid

	MinTemperature   float64               // the minimum temperature of the fluid (C)
	MaxTemperature   float64               // the minimum temperature of the fluid (C)
	TemperatureCurve func(float64) float64 // the curve to calculate the applied temperature from the unit value

	MinViscosity   float64               // the minimum viscosity of the fluid (lower = less flow) (cP)
	MaxViscosity   float64               // the maximum viscosity of the fluid (higher = more flow) (cP)
	ViscosityCurve func(float64) float64 // the curve to calculate the applied viscosity from the unit value

	HeatAbsorptionRate   float64 // how much heat is absorbed by the fluid (C/tick)
	MaxTemperatureDelta  float64 // the maximum temperature delta (C/tick)
	ThermalExpansionRate float64 // the rate at which the fluid changes pressure
}

func NewWater() Fluid {
	fluid := Fluid{
		Name: "Water",

		MinTemperature: 0.0,
		MaxTemperature: 100.0,

		MinViscosity: 0.3,
		MaxViscosity: 1.0,

		HeatAbsorptionRate:   0.03,
		MaxTemperatureDelta:  0.05,
		ThermalExpansionRate: 0.8,
	}

	return fluid

}

func NewPropyleneGlycol() Fluid {
	fluid := Fluid{
		Name:                 "Propylene-Glycol",
		MinTemperature:       -50.0,
		MaxTemperature:       180.0,
		MinViscosity:         1.0,
		MaxViscosity:         50.0,
		HeatAbsorptionRate:   0.05,
		MaxTemperatureDelta:  0.08,
		ThermalExpansionRate: 0.4,
	}

	return fluid

}
