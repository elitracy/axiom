package cooling

import (
	"github.com/elias/axiom/engine/systems"
	"github.com/elias/axiom/engine/systems/components"
	"github.com/elias/axiom/engine/utils"
)

const (
	startingCoolantFluid   = 1.0
	minWaterTemp           = 5.0
	maxWaterTemp           = 90.0
	minPropyleneGlycolTemp = -50.0
	maxPropyleneGlycolTemp = 180.0
)

// A cooling loop wraps the heat source and offsets the heat using the fluid contained
type CoolantLoop struct {
	*systems.SystemCore
	fluid       components.Component
	temperature components.Component
	viscosity   components.Component
	health      *components.Health
}

func NewWaterCoolant() *CoolantLoop {
	temperatureCurve := func(x float64) float64 {
		normalized := (x - minWaterTemp) / (maxWaterTemp - minWaterTemp)
		return (utils.Tanh(normalized, 3.1, 0)+0.015)*(maxWaterTemp-minWaterTemp) + minWaterTemp
	}

	system := &CoolantLoop{
		SystemCore:  systems.NewSystemCore("Cooling Loop"),
		fluid:       components.NewComponent("Fluid", 1.0),
		temperature: components.NewComponent("Temperature", 0.0, temperatureCurve),
	}

	return system
}

func NewCoolingLoop(coolant CoolantType) *CoolantLoop {
}
