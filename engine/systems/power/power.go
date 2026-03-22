package power

import (
	"fmt"

	"github.com/elias/axiom/engine/materials"
	"github.com/elias/axiom/engine/systems"
	"github.com/elias/axiom/engine/systems/components"
	"github.com/elias/axiom/engine/utils"
)

type PowerInput struct {
	CoolantTemperature float64
	AmbientTemperature float64
}

type PowerOutput struct {
	Power       float64
	Temperature float64
}

// A system that creates power
type Power interface {
	systems.System
	// Returns a readonly reference to the temperature component of the power source.
	Temperature() components.ComponentReader
	// Returns a readonly reference to the power component of the power source.
	Power() components.ComponentReader
	// Returns a readonly reference to the fuel component of the power source.
	Fuel() components.ComponentReader
	// Returns status of the power source.
	Status() systems.Status
	// Updates the power system.
	// Returns relavant info about the system.
	// input is the power inputs for the update.
	Tick(input PowerInput) PowerOutput
}

// TODO: create generator input
const (
	startingPowerOutput = 0.0
	startingFuel        = 1.0
	startingHealth      = 1.0

	fuelLostPerTick   = 0.001
	healthLostPerTick = 0.001
)

type PowerCore struct {
	*systems.SystemCore

	housingMaterial materials.Metal

	power       components.Component
	fuel        components.Component
	temperature components.Component
	health      *components.Health

	heatGenerationRate float64
	powerGrowthRate    float64
}

func (s *PowerCore) Temperature() components.ComponentReader { return s.temperature }
func (s *PowerCore) Power() components.ComponentReader       { return s.power }
func (s *PowerCore) Fuel() components.ComponentReader        { return s.fuel }
func (s *PowerCore) Status() systems.Status                  { return s.health.Status() }

// Updates the component values for the generator
func (s *PowerCore) Tick(input PowerInput) PowerOutput {
	if s.health.Status() == systems.Offline {
		s.power.SetNorm(0.0)
	}

	if s.fuel.Norm() <= 0.0 {
		s.power.SetNorm(s.power.Norm() - s.powerGrowthRate)
	} else {
		s.power.SetNorm(s.power.Norm() + s.powerGrowthRate)
	}

	if input.CoolantTemperature < s.temperature.Value() {
		coolingNorm := (input.CoolantTemperature - s.temperature.Min()) / (s.temperature.Max() - s.temperature.Min())

		temperatureDelta := (s.temperature.Norm() - coolingNorm) * s.housingMaterial.HeatAbsorptionRate
		temperatureDelta = utils.Clamp(0, temperatureDelta, s.housingMaterial.MaxTemperatureDelta)

		s.temperature.SetNorm(s.temperature.Norm() - temperatureDelta)
	}

	if s.power.Norm() > 0.0 {
		s.fuel.SetNorm(s.fuel.Norm() - fuelLostPerTick)
		s.temperature.SetNorm(s.temperature.Norm() + s.heatGenerationRate)
	}

	if s.power.Norm() <= 0.0 && s.temperature.Value() > input.AmbientTemperature {
		s.temperature.SetNorm(s.temperature.Norm() - s.housingMaterial.HeatAbsorptionRate)
	}

	if s.temperature.Norm() >= 1.0 {
		s.health.SetNorm(s.health.Norm() - healthLostPerTick)
	}

	output := PowerOutput{
		Power:       s.power.Value(),
		Temperature: s.temperature.Value(),
	}

	return output
}

// Returns the stringified information for the generator
func (s *PowerCore) String() string {

	output := fmt.Sprintf("%v: %v", s.ID(), s.Name())
	output += fmt.Sprintf("\n%s", s.power)
	output += fmt.Sprintf("\n%s", s.fuel)
	output += fmt.Sprintf("\n%s", s.temperature)
	output += fmt.Sprintf("\n%s", s.health)

	return output
}
