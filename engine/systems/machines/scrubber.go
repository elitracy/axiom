package machines

import (
	"fmt"

	"github.com/elias/axiom/engine/systems"
	"github.com/elias/axiom/engine/systems/components"
)

const (
	o2ConsumptionRate  = 0.002
	co2ProductionRatio = 0.8

	targetCo2 = 0.03
	targetO2  = 0.21
)

type ScrubberInput struct {
	PowerAvailable float64
}

type ScrubberOutput struct {
	Status systems.Status
}

type Scrubber struct {
	*systems.SystemCore

	health *components.Health
	o2     components.Component
	co2    components.Component

	scrubberRate  float64
	requiredPower float64
}

func NewScrubber() *Scrubber {
	system := &Scrubber{
		SystemCore:    systems.NewSystemCore("Life Support"),
		health:        components.NewHealthComponent(1.0),
		o2:            components.NewComponent("O2 (%)", targetO2, 0.1, 0.3),
		co2:           components.NewComponent("CO2 (%)", targetCo2, 0.0, 0.08),
		scrubberRate:  o2ConsumptionRate * co2ProductionRatio,
		requiredPower: 600.0,
	}

	return system
}

func (s *Scrubber) Status() systems.Status {
	switch {
	case s.o2.Value() <= s.o2.Min() || s.co2.Value() >= s.co2.Max():
		return systems.Offline
	case s.o2.Value() <= .14 || s.co2.Value() >= .06:
		return systems.Critical
	case s.o2.Value() <= .18 || s.co2.Value() >= .04:
		return systems.Degraded
	default:
		return systems.Online
	}

}

func (s *Scrubber) Tick(input ScrubberInput) ScrubberOutput {

	s.o2.SetNorm(s.o2.Norm() - o2ConsumptionRate)
	s.co2.SetNorm(s.co2.Norm() + o2ConsumptionRate*co2ProductionRatio)

	output := ScrubberOutput{
		Status: s.Status(),
	}

	effectiveness := 1.0

	switch {
	case input.PowerAvailable <= 0.0:
		effectiveness = 0.0
	case input.PowerAvailable <= s.requiredPower*.8:
		effectiveness = 0.7
	default:
		effectiveness = 1.0
	}

	co2Removed := s.scrubberRate * effectiveness
	s.co2.SetNorm(s.co2.Norm() - co2Removed)
	s.o2.SetNorm(s.o2.Norm() + s.scrubberRate/co2ProductionRatio*effectiveness)

	return output
}

func (s *Scrubber) String() string {

	output := fmt.Sprintf("%v: %v", s.ID(), s.Name())
	output += fmt.Sprintf("\n%s", s.o2)
	output += fmt.Sprintf("\n%s", s.co2)
	output += fmt.Sprintf("\n%s", s.health)

	return output
}
