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
	startingPower  = 1.0
	startingFuel   = 1.0
	startingHealth = 1.0

	percentFuelLostPerTick   = 1.0 / systems.TICKS_TILL_DEATH_DEBUG
	percentHealthLostPerTick = 1.0 / systems.TICKS_TILL_DEATH_DEBUG
	percentPowerUsedPerTick  = 1.0 / systems.TICKS_TILL_DEATH_DEBUG // NOTE: will be replaced by power output like cooling source

)

type PowerCore struct {
	*systems.SystemCore

	housingMaterial materials.Metal

	power       components.Component
	fuel        components.Component
	temperature components.Component
	health      *components.Health
}

func (s *PowerCore) Temperature() components.ComponentReader { return s.temperature }
func (s *PowerCore) Power() components.ComponentReader       { return s.power }
func (s *PowerCore) Fuel() components.ComponentReader        { return s.fuel }
func (s *PowerCore) Status() systems.Status                  { return s.health.Status() }

// Updates the component values for the generator
func (s *PowerCore) Tick(input PowerInput) PowerOutput {
	if s.fuel.Value() <= s.fuel.Min() {
		s.power.SetNorm(s.power.Norm() - percentPowerUsedPerTick)
	}

	if input.CoolantTemperature < s.temperature.Value() {
		temperatureDelta := (s.temperature.Norm() - input.CoolantTemperature/s.temperature.Max())
		temperatureDelta = utils.Clamp(-s.housingMaterial.MaxTemperatureDelta, temperatureDelta, s.housingMaterial.MaxTemperatureDelta)

		s.temperature.SetNorm(s.temperature.Norm() + temperatureDelta)
	}

	if s.health.Status() == systems.Offline {
		s.power.SetNorm(0.0)
	}

	if s.power.Value() > s.power.Min() {
		s.fuel.SetNorm(s.fuel.Norm() - percentFuelLostPerTick)
		s.temperature.SetNorm(s.temperature.Norm() + s.housingMaterial.HeatAbsorptionRate)
	}

	if s.power.Value() <= s.power.Min() && s.temperature.Value() > input.AmbientTemperature {
		s.temperature.SetNorm(s.temperature.Norm() - s.housingMaterial.HeatAbsorptionRate)
	}

	if s.temperature.Value() >= s.temperature.Max() {
		s.health.SetNorm(s.health.Norm() - percentHealthLostPerTick)
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
