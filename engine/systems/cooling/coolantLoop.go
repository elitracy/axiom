package cooling

import (
	"github.com/elias/axiom/engine/systems"
	"github.com/elias/axiom/engine/systems/components"
	"github.com/elias/axiom/engine/utils"
)

const (
	startingFluid       = 1.0
	startingTemperature = 0.0
	startingViscosity   = 0.0
	startingHealth      = 1.0
	startingPressure    = 0.0

	minWaterTemp           = 5.0
	maxWaterTemp           = 90.0
	minPropyleneGlycolTemp = -50.0
	maxPropyleneGlycolTemp = 180.0

	maxWaterViscosity           = 1.0 // cP (centipoise)
	minWaterViscosity           = 0.3
	maxPropyleneGlycolViscosity = 50.0
	minPropyleneGlycolViscosity = 1.0

	waterCoolingCoefficient           = 2.0
	propyleneGlycolCoolingCoefficient = 1.0

	minPressure = 1.0
	maxPressure = 800.0

	percentHealthLostPerTick = 1.0 / systems.TICKS_TILL_DEATH_DEBUG
)

// A cooling loop wraps the heat source and offsets the heat using the fluid contained
type CoolantLoop struct {
	*systems.SystemCore
	fluid       components.Component
	temperature components.Component
	viscosity   components.Component
	pressure    components.Component
	health      *components.Health

	heatInput          float64
	coolingCoefficient float64
}

func (s *CoolantLoop) Tick() {
	s.health.SetValue(s.health.Value() - percentHealthLostPerTick)

	pressure := s.pressure.ApplyValueCurve()
	viscosity := s.viscosity.ApplyValueCurve()
	temperature := s.temperature.ApplyValueCurve()

	flowRate := s.pressure.Value()
	if s.viscosity.Value() > 0 {
		flowRate = utils.Clamp(0.0, utils.Tanh(s.pressure.Value()/s.viscosity.Value(), 2.0, 0.0), 1.0)
	}
	dissipation := flowRate * s.coolingCoefficient
	tempDelta := s.heatInput - dissipation

}

func (s *CoolantLoop) SetHeatInput(temp float64) { s.heatInput = temp }

func NewWaterCoolantLoop() *CoolantLoop {
	temperatureCurve := func(x float64) float64 {
		temp := (utils.Tanh(x, 3.1, 0)+0.015)*(maxWaterTemp-minWaterTemp) + minWaterTemp
		return utils.Clamp(minWaterTemp, temp, maxWaterTemp)
	}

	viscosityCurve := func(x float64) float64 {
		viscosity := (1-utils.Tanh(x, 4.2, 0))*(maxWaterViscosity-minWaterViscosity) + minWaterViscosity
		return utils.Clamp(minWaterViscosity, viscosity, maxWaterViscosity)
	}

	pressureCurve := func(x float64) float64 {
		pressure := utils.Tanh(x, 4.2, 0)*(maxPressure-minPressure) + minPressure
		return utils.Clamp(float64(minPressure), pressure, float64(maxPressure))
	}

	system := &CoolantLoop{
		SystemCore:         systems.NewSystemCore("Cooling Loop"),
		fluid:              components.NewComponent("Fluid (%%)", startingFluid),
		temperature:        components.NewComponent("Temperature (C)", startingTemperature, temperatureCurve),
		viscosity:          components.NewComponent("Viscosity (cP)", startingViscosity, viscosityCurve),
		pressure:           components.NewComponent("Pressure (kPa)", startingPressure, pressureCurve),
		health:             components.NewHealthComponent(startingHealth),
		coolingCoefficient: waterCoolingCoefficient,
	}

	return system
}

func NewPropyleneGlycolCoolantLoop() *CoolantLoop {
	temperatureCurve := func(x float64) float64 {
		temp := (utils.Tanh(x, 3.1, 0)+0.015)*(maxPropyleneGlycolTemp-minPropyleneGlycolTemp) + minPropyleneGlycolTemp
		return utils.Clamp(minPropyleneGlycolTemp, temp, maxPropyleneGlycolTemp)
	}

	viscosityCurve := func(x float64) float64 {
		viscosity := (1-utils.Tanh(x, 10.0, 0))*(maxPropyleneGlycolViscosity-minPropyleneGlycolViscosity) + minPropyleneGlycolViscosity
		return utils.Clamp(minPropyleneGlycolViscosity, viscosity, maxPropyleneGlycolViscosity)
	}

	pressureCurve := func(x float64) float64 {
		pressure := utils.Tanh(x, 4.2, 0)*(maxPressure-minPressure) + minPressure
		return utils.Clamp(float64(minPressure), pressure, float64(maxPressure))
	}

	system := &CoolantLoop{
		SystemCore:         systems.NewSystemCore("Cooling Loop"),
		fluid:              components.NewComponent("Fluid", startingFluid),
		temperature:        components.NewComponent("Temperature", startingTemperature, temperatureCurve),
		viscosity:          components.NewComponent("Viscosity", startingViscosity, viscosityCurve),
		pressure:           components.NewComponent("Pressure (kPa)", startingPressure, pressureCurve),
		health:             components.NewHealthComponent(startingHealth),
		coolingCoefficient: propyleneGlycolCoolingCoefficient,
	}

	return system
}
