package machines

import (
	"fmt"

	"github.com/elias/axiom/engine/systems"
	"github.com/elias/axiom/engine/systems/components"
	"github.com/elias/axiom/engine/utils"
)

const (
	o2ConsumptionRate  = 0.008
	co2ProductionRatio = 0.8

	targetCO2 = 0.03
	targetO2  = 0.21

	minO2 = 0.1
	maxO2 = 0.3

	minCO2 = 0.0
	maxCO2 = 0.08
)

type ScrubberInput struct {
	PowerAvailable float64
}

type ScrubberOutput struct {
	Status systems.Status
	O2     float64
	CO2    float64
}

type Scrubber struct {
	*systems.SystemCore

	health *components.Health
	o2     components.Component
	co2    components.Component

	scrubberRate  float64
	powerCapacity float64
}

func NewScrubber() *Scrubber {

	targetO2Norm := (targetO2 - minO2) / (maxO2 - minO2)
	targetCO2Norm := (targetCO2 - minCO2) / (maxCO2 - minCO2)

	system := &Scrubber{
		SystemCore:    systems.NewSystemCore("Life Support"),
		health:        components.NewHealthComponent(1.0),
		o2:            components.NewComponent("O2 (%)", targetO2Norm, minO2, maxO2),
		co2:           components.NewComponent("CO2 (%)", targetCO2Norm, minCO2, maxCO2),
		scrubberRate:  o2ConsumptionRate * co2ProductionRatio,
		powerCapacity: 600.0,
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

	effectiveness := input.PowerAvailable / max(s.powerCapacity, 1)
	effectiveness = utils.Clamp(0, effectiveness, 1)

	co2Removed := s.scrubberRate * effectiveness
	s.co2.SetNorm(s.co2.Norm() - co2Removed)
	s.o2.SetNorm(s.o2.Norm() + s.scrubberRate/co2ProductionRatio*effectiveness)

	output := ScrubberOutput{
		Status: s.Status(),
		O2:     s.o2.Value(),
		CO2:    s.co2.Value(),
	}
	return output
}

func (s *Scrubber) String() string {
	output := fmt.Sprintf("%v: %v", s.ID(), s.Name())
	output += fmt.Sprintf("\n%s", s.o2)
	output += fmt.Sprintf("\n%s", s.co2)
	output += fmt.Sprintf("\n%s", s.health)

	return output
}
