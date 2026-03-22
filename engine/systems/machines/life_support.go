package machines

import (
	"fmt"

	"github.com/elias/axiom/engine/systems"
	"github.com/elias/axiom/engine/systems/components"
)

type LifeSupportInput struct {
	PowerAvailable    float64
	TemperatureStatus systems.Status
	OxygenStatus      systems.Status
}

type LifeSupportOutput struct {
	Status systems.Status
}

type LifeSupport struct {
	*systems.SystemCore
	health        *components.Health
	requiredPower float64
}

func NewLifeSupport() *LifeSupport {
	system := &LifeSupport{
		SystemCore:    systems.NewSystemCore("Life Support"),
		health:        components.NewHealthComponent(1.0),
		requiredPower: 600.0,
	}

	return system
}

func (s *LifeSupport) Status() systems.Status { return s.health.Status() }

func (s *LifeSupport) Tick(input LifeSupportInput) LifeSupportOutput {
	var output LifeSupportOutput

	if input.PowerAvailable == 0.0 {
		output.Status = systems.Offline
		return output
	}

	output.Status = min(input.TemperatureStatus, input.OxygenStatus)
	return output
}

func (s *LifeSupport) String() string {

	output := fmt.Sprintf("%v: %v", s.ID(), s.Name())
	output += fmt.Sprintf("\n%s", s.health)

	return output
}
