package cooling

import (
	"fmt"
	"log"

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
		temperature:  components.NewComponent("Temperature (C)", 0.0, coolantFluid.MinTemperature, coolantFluid.MaxTemperature, coolantFluid.TemperatureCurve),
		viscosity:    components.NewComponent("Viscosity (cP)", 0.5, coolantFluid.MinViscosity, coolantFluid.MaxViscosity, coolantFluid.ViscosityCurve),
		pressure:     components.NewComponent("Pressure (kPa)", 0.5, pipeMetal.MinPressure, pipeMetal.MaxPressure, pipeMetal.PressureCurve),
		health:       components.NewHealthComponent(1.0),
		coolantFluid: coolantFluid,
		pipeMetal:    pipeMetal,
	}

	system.viscosity.SetNorm(1 - system.temperature.Norm())
	system.pressure.SetNorm(system.volume.Norm() * system.temperature.Norm() * system.coolantFluid.ThermalExpansionRate)

	return system
}

func (s *CoolantCore) Status() systems.Status { return s.health.Status() }

func (s *CoolantCore) Tick(input CoolantInput) CoolantOutput {
	s.health.SetNorm(s.health.Norm() - healthDecayPerTick)

	flow := 0.0
	if s.viscosity.Norm() > 0 {
		flow = s.pressure.Norm() * (1 - s.viscosity.Norm()) * s.volume.Norm()
	}
	log.Printf("FLOW: %.2f", flow)

	cooling := flow * s.coolantFluid.HeatAbsorptionRate
	log.Printf("COOLING: %.2f", cooling)

	normalizeLoad := input.LoadTemperature / s.temperature.Max()
	log.Printf("NORM LOAD: %.2f", normalizeLoad)

	temperatureDelta := normalizeLoad - cooling
	temperatureDelta = utils.Clamp(-s.coolantFluid.MaxTemperatureDelta, temperatureDelta, s.coolantFluid.MaxTemperatureDelta)

	log.Printf("DELTA: %.2f", temperatureDelta)

	s.temperature.SetNorm(s.temperature.Norm() + temperatureDelta)
	log.Printf("TEMP: %.2f", s.temperature.Value())

	if s.temperature.Value() >= s.temperature.Max() {
		s.volume.SetNorm(s.volume.Norm() - volumeLossPerTick)
	}

	s.viscosity.SetNorm(s.calculateViscosity())
	s.pressure.SetNorm(s.calculatePressure())

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

func (s *CoolantCore) calculateViscosity() float64 {
	return 1 - s.temperature.Norm()
}

func (s *CoolantCore) calculatePressure() float64 {
	pressure := s.volume.Norm() * s.temperature.Norm() * s.coolantFluid.ThermalExpansionRate
	return pressure
}
