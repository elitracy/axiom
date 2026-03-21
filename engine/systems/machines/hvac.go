package machines

import (
	"fmt"
	"log"

	"github.com/elias/axiom/engine/systems"
	"github.com/elias/axiom/engine/systems/components"
	"github.com/elias/axiom/engine/utils"
)

const (
	heatRate = 0.005
	hvacRate = 0.03

	minLivableTemp = -10.0
	maxLivableTemp = 50.0

	hvacEfficiencyZeroThreshold   = 0.0
	hvacEfficiencyMediumThreshold = 0.5
	hvacEfficiencyHighThreshold   = .9
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

	maxTempPush       float64
	targetTemperature float64
	requiredPower     float64
}

func NewHvac(targetTemperature float64) *Hvac {

	// NOTE: Linear interpolation only works while temp curve is linear!!
	normTargetTemp := (targetTemperature - minLivableTemp) / (maxLivableTemp - minLivableTemp)

	system := &Hvac{
		SystemCore:        systems.NewSystemCore("Life Support"),
		temperature:       components.NewComponent("Bunker Temperature (C)", normTargetTemp, minLivableTemp, maxLivableTemp),
		health:            components.NewHealthComponent(1.0),
		maxTempPush:       0.03,
		targetTemperature: targetTemperature,
		requiredPower:     1200,
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
	log.Printf("TEMP BEFORE: %v", s.temperature.Value())

	effectiveness := calculateEffectiveness(input.PowerAvailable, s.requiredPower)

	log.Printf("EFF: %v", effectiveness)
	averageHeat := calculateAverageHeat(input.HeatSources)

	log.Printf("AVG HEAT: %v", averageHeat)

	heatPush := calculateHeatPush(averageHeat, s.temperature.Value(), heatRate)
	hvacPush := calculateHeatPush(s.targetTemperature, s.temperature.Value(), effectiveness*hvacRate)

	log.Printf("heatPush: %.2f", heatPush)
	log.Printf("hvacPush: %.2f", hvacPush)

	netPush := (heatPush + hvacPush) / (s.temperature.Max() - s.temperature.Min())
	log.Printf("net: %.2f", netPush)
	netPush = utils.Clamp(-s.maxTempPush, netPush, s.maxTempPush)

	log.Printf("net clamped: %.2f", netPush)
	s.temperature.SetNorm(s.temperature.Norm() + netPush)

	output := HvacOutput{
		Status:            s.Status(),
		TargetTemperature: s.targetTemperature,
		Temperature:       s.temperature.Value(),
	}

	log.Printf("TEMP: %v", output.Temperature)
	return output
}

func calculateAverageHeat(heatSources []float64) float64 {

	averageHeat := 0.0
	for _, heat := range heatSources {
		averageHeat += heat
	}

	averageHeat /= max(float64(len(heatSources)), 1.0)

	return averageHeat
}

func calculateHeatPush(targetTemp float64, currentTemp, heatingRate float64) float64 {
	return (targetTemp - currentTemp) * heatingRate
}

func calculateEffectiveness(powerAvailable, requiredPower float64) float64 {
	if powerAvailable >= .25*requiredPower {
		return min(powerAvailable/requiredPower, 1.0)
	}

	return 0.0
}

func (s *Hvac) String() string {

	output := fmt.Sprintf("%v: %v", s.ID(), s.Name())
	output += fmt.Sprintf("\n%s", s.temperature)
	output += fmt.Sprintf("\n%s", s.health)

	return output
}
