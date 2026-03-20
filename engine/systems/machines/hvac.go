package machines

import (
	"fmt"
	"log"

	"github.com/elias/axiom/engine/systems"
	"github.com/elias/axiom/engine/systems/components"
	"github.com/elias/axiom/engine/utils"
)

const (
	heatRate = 0.05

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

	maxTemperatureNormDelta float64
	targetTemperature       float64
	requiredPower           float64
}

func NewHvac(targetTemperature float64) *Hvac {

	// NOTE: Linear interpolation only works while temp curve is linear!!
	normTargetTemp := (targetTemperature - minLivableTemp) / (maxLivableTemp - minLivableTemp)

	system := &Hvac{
		SystemCore:              systems.NewSystemCore("Life Support"),
		temperature:             components.NewComponent("Bunker Temperature (C)", normTargetTemp, minLivableTemp, maxLivableTemp),
		health:                  components.NewHealthComponent(1.0),
		maxTemperatureNormDelta: 0.02,
		targetTemperature:       targetTemperature,
		requiredPower:           1200,
	}

	return system
}

func (s *Hvac) Status() systems.Status {
	switch {
	case s.temperature.Value() >= maxLivableTemp || s.temperature.Value() <= minLivableTemp:
		return systems.Offline
	case s.temperature.Value() >= 40.0 || s.temperature.Value() <= 10.0:
		return systems.Critical
	case s.temperature.Value() >= 30.0 || s.temperature.Value() < 20.0:
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

	averageHeat := 0.0
	for _, heat := range input.HeatSources {
		averageHeat += heat
	}
	if len(input.HeatSources) > 0 {
		averageHeat /= float64(len(input.HeatSources))
	}

	log.Printf("AVG HEAT: %v", averageHeat)

	temperatureRange := s.temperature.Max() - s.temperature.Min()

	heatPushNorm := (averageHeat - s.temperature.Value()) / temperatureRange * heatRate
	log.Printf("heatPush: %.2f", heatPushNorm)
	heatPushNorm = utils.Clamp(-s.maxTemperatureNormDelta, heatPushNorm, s.maxTemperatureNormDelta)

	hvacPushNorm := ((s.targetTemperature - s.temperature.Value()) / temperatureRange) * effectiveness
	log.Printf("hvacPush: %.2f", hvacPushNorm)
	hvacPushNorm = utils.Clamp(-s.maxTemperatureNormDelta, hvacPushNorm, s.maxTemperatureNormDelta)

	log.Printf("heatPush: %.2f", heatPushNorm)
	log.Printf("hvacPush: %.2f", hvacPushNorm)

	netPushNorm := heatPushNorm + hvacPushNorm
	log.Printf("net: %.2f", netPushNorm)
	s.temperature.SetNorm(s.temperature.Norm() + netPushNorm)

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
