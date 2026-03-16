package cooling

import (
	"fmt"

	"github.com/elias/axiom/engine/materials"
	"github.com/elias/axiom/engine/systems"
	"github.com/elias/axiom/engine/systems/components"
	"github.com/elias/axiom/engine/utils"
)

const (
	healthDecayPerTick = 1.0 / systems.TICKS_TILL_DEATH_DEBUG
	volumeLossPerTick  = 1.0 / systems.TICKS_TILL_DEATH_DEBUG
)

type CoolantInput struct {
	LoadTemperature float64
}

type CoolantOutput struct {
	Temperature float64
}

type Coolant interface {
	systems.System
	// Returns a readonly reference to the temperature component of the coolant
	Temperature() components.ComponentReader
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
}

func NewCoolantLoop(coolantFluid materials.Fluid, pipeMetal materials.Metal) *CoolantCore {

	system := &CoolantCore{
		SystemCore:   systems.NewSystemCore("Cooling Loop"),
		volume:       components.NewComponent("Volume (%)", 1.0, 0.0, 1.0),
		temperature:  components.NewComponent("Temperature (C)", coolantFluid.MinTemperature, coolantFluid.MinTemperature, coolantFluid.MaxTemperature, coolantFluid.TemperatureCurve),
		viscosity:    components.NewComponent("Viscosity (cP)", coolantFluid.MinViscosity, coolantFluid.MinViscosity, coolantFluid.MaxViscosity, coolantFluid.ViscosityCurve),
		pressure:     components.NewComponent("Pressure (kPa)", pipeMetal.MinPressure, pipeMetal.MinPressure, pipeMetal.MaxPressure, pipeMetal.PressureCurve),
		health:       components.NewHealthComponent(1.0),
		coolantFluid: coolantFluid,
		pipeMetal:    pipeMetal,
	}

	return system
}

func (s *CoolantCore) Temperature() components.ComponentReader { return s.temperature }
func (s *CoolantCore) Status() systems.Status                  { return s.health.Status() }

func (s *CoolantCore) Tick(input CoolantInput) CoolantOutput {
	s.health.SetNorm(s.health.Norm() - healthDecayPerTick)

	flowRate := 0.0
	if s.viscosity.Norm() > 0 {
		flowRate = s.pressure.Norm() * (1 - s.viscosity.Norm()) * s.volume.Norm()
	}

	cooling := flowRate * s.coolantFluid.HeatAbsorptionRate
	temperatureDelta := input.LoadTemperature/s.temperature.Max() - cooling
	temperatureDelta = utils.Clamp(-s.coolantFluid.ThermalConductivityRate, temperatureDelta, s.coolantFluid.ThermalConductivityRate)

	s.temperature.SetNorm(s.temperature.Norm() + temperatureDelta)

	if s.temperature.Value() >= s.temperature.Max() {
		s.volume.SetNorm(s.volume.Norm() - volumeLossPerTick)
	}

	s.viscosity.SetNorm(1 - s.temperature.Norm())
	s.pressure.SetNorm(s.volume.Norm() * s.temperature.Norm() * s.coolantFluid.ThermalExpansionRate)

	if s.pressure.Value() >= s.pressure.Max() {
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
