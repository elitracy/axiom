package machines

import (
	"fmt"

	"github.com/elias/axiom/engine/systems"
	"github.com/elias/axiom/engine/systems/components"
)

const (
	heatBleedRate = 0.01

	minLivableTemp = -10.0
	maxLivableTemp = 50.0
)

type HvacInput struct {
	PowerAvailable float64
	HeatSources    []float64
}

type HvacOutput struct {
	Status            systems.Status
	TargetTemperature float64
	Temperature       float64
}

type Hvac struct {
	*systems.SystemCore

	temperature components.Component
	health      *components.Health

	maxTemperatureDelta float64
	targetTemperature   float64
	requiredPower       float64
}

func NewHvac(targetTemperature float64) *Hvac {

	// NOTE: Linear interpolation only works while temp curve is linear!!
	normTargetTemp := (targetTemperature - minLivableTemp) / (maxLivableTemp - minLivableTemp)

	system := &Hvac{
		SystemCore:          systems.NewSystemCore("Life Support"),
		temperature:         components.NewComponent("Bunker Temperature (C)", normTargetTemp, minLivableTemp, maxLivableTemp),
		health:              components.NewHealthComponent(1.0),
		maxTemperatureDelta: 0.02,
		targetTemperature:   targetTemperature,
		requiredPower:       1200,
	}

	return system
}

func (s *Hvac) Status() systems.Status {
	switch {
	case s.temperature.Norm() >= 1.0 || s.temperature.Norm() <= 0.0:
		return systems.Offline
	case s.temperature.Value() >= 40.0 || s.temperature.Value() <= 0.0:
		return systems.Critical
	case s.temperature.Value() >= 30.0 || s.temperature.Value() <= 20.0:
		return systems.Degraded
	default:
		return systems.Online
	}
}

func (s *Hvac) Tick(input HvacInput) HvacOutput {
	effectiveness := 1.0
	switch {
	case input.PowerAvailable <= 0.0:
		effectiveness = 0.0
	case input.PowerAvailable <= s.requiredPower*.8:
		effectiveness = 0.7
	default:
		effectiveness = 1.0
	}

	totalBleed := 0.0
	for _, heat := range input.HeatSources {
		delta := (heat - s.temperature.Value()) / s.temperature.Max()
		if delta > 0 {
			totalBleed += delta * heatBleedRate
		}
	}

	s.temperature.SetNorm(s.temperature.Norm() + totalBleed)

	temperatureDelta := s.maxTemperatureDelta * effectiveness

	if s.temperature.Value() > s.targetTemperature {
		s.temperature.SetNorm(s.temperature.Norm() - temperatureDelta)
	}
	if s.temperature.Value() < s.targetTemperature {
		s.temperature.SetNorm(s.temperature.Norm() + temperatureDelta)
	}

	output := HvacOutput{
		Status:            s.Status(),
		TargetTemperature: s.targetTemperature,
		Temperature:       s.temperature.Value(),
	}

	return output
}

func (s *Hvac) String() string {

	output := fmt.Sprintf("%v: %v", s.ID(), s.Name())
	output += fmt.Sprintf("\n%s", s.temperature)
	output += fmt.Sprintf("\n%s", s.health)

	return output
}
