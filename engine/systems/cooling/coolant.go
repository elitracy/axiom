package cooling

import (
	"fmt"

	"github.com/elias/axiom/engine/materials"
	"github.com/elias/axiom/engine/systems"
	"github.com/elias/axiom/engine/systems/components"
	"github.com/elias/axiom/engine/utils"
)

const (
	healthDecayPerTick = 0.001
	volumeLossPerTick  = 0.001

	basePressure = 0.05
	minFlow      = 0.05
)

// The input to the coolant tick function
type CoolantInput struct {
	LoadTemperature float64
}

// The output of the coolant tick function
type CoolantOutput struct {
	Temperature float64
}

type Coolant interface {
	systems.System
	// Updates the coolant system.
	// Returns relavant info about the system.
	// input is the coolant inputs for the update.
	Tick(input CoolantInput) CoolantOutput
}

// A cooling loop wraps the heat source and offsets the heat using the fluid contained
type CoolantCore struct {
	*systems.SystemCore

	coolantFluid materials.Fluid
	pipeMetal    materials.Metal

	volume      components.Component
	temperature components.Component
	viscosity   components.Component
	pressure    components.Component
	health      *components.Health

	basePressure float64
}

// Creates a new coolant loop
func NewCoolantLoop(coolantFluid materials.Fluid, pipeMetal materials.Metal) *CoolantCore {

	system := &CoolantCore{
		SystemCore:   systems.NewSystemCore("Cooling Loop"),
		volume:       components.NewComponent("Volume (%)", 1.0, 0.0, 1.0),
		temperature:  components.NewComponent("Temperature (C)", 0.01, coolantFluid.MinTemperature, coolantFluid.MaxTemperature, coolantFluid.TemperatureCurve),
		viscosity:    components.NewComponent("Viscosity (cP)", 0.0, coolantFluid.MinViscosity, coolantFluid.MaxViscosity, coolantFluid.ViscosityCurve),
		pressure:     components.NewComponent("Pressure (kPa)", 0.0, pipeMetal.MinPressure, pipeMetal.MaxPressure, pipeMetal.PressureCurve),
		health:       components.NewHealthComponent(1.0),
		coolantFluid: coolantFluid,
		pipeMetal:    pipeMetal,

		basePressure: basePressure,
	}

	system.viscosity.SetNorm(system.calculateViscosityNorm())
	system.pressure.SetNorm(system.calculatePressureNorm())

	return system
}

func (s *CoolantCore) Status() systems.Status { return s.health.Status() }

func (s *CoolantCore) Tick(input CoolantInput) CoolantOutput {
	s.health.SetNorm(s.health.Norm() - healthDecayPerTick)

	flow := max(minFlow, s.pressure.Norm()*(1-s.viscosity.Norm())*s.volume.Norm())

	percentHeatAbsorbed := flow * s.coolantFluid.HeatAbsorptionRate

	normalizeLoad := (input.LoadTemperature - s.temperature.Min()) / (s.temperature.Max() - s.temperature.Min())

	temperatureDelta := normalizeLoad * percentHeatAbsorbed
	temperatureDelta = utils.Clamp(-s.coolantFluid.MaxTemperatureDelta, temperatureDelta, s.coolantFluid.MaxTemperatureDelta)

	if temperatureDelta > 0 {
		s.temperature.SetNorm(s.temperature.Norm() + temperatureDelta)
	}

	if s.temperature.Norm() >= 1.0 {
		s.volume.SetNorm(s.volume.Norm() - volumeLossPerTick)
	}

	s.viscosity.SetNorm(s.calculateViscosityNorm())
	s.pressure.SetNorm(s.calculatePressureNorm())

	if s.pressure.Norm() >= 1.0 {
		s.health.SetNorm(s.health.Norm() - healthDecayPerTick)
	}

	output := CoolantOutput{
		Temperature: s.temperature.Value(),
	}

	return output
}

func (s *CoolantCore) String() string {
	output := fmt.Sprintf("%v: %v", s.ID(), s.Name())
	output += fmt.Sprintf("\n%s", s.volume)
	output += fmt.Sprintf("\n%s", s.temperature)
	output += fmt.Sprintf("\n%s", s.viscosity)
	output += fmt.Sprintf("\n%s", s.pressure)
	output += fmt.Sprintf("\n%s", s.health)

	return output
}

func (s *CoolantCore) calculateViscosityNorm() float64 {
	viscosity := 1 - s.temperature.Norm()
	viscosity = utils.Clamp(0.0, viscosity, 1.0)

	return viscosity
}

func (s *CoolantCore) calculatePressureNorm() float64 {
	pressure := s.basePressure + s.volume.Norm()*s.temperature.Norm()*s.coolantFluid.ThermalExpansionRate
	pressure = utils.Clamp(0.0, pressure, 1.0)

	return pressure
}
