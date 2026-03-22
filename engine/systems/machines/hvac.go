package machines

import (
	"fmt"
	"math"

	"github.com/elias/axiom/engine/systems"
	"github.com/elias/axiom/engine/systems/components"
	"github.com/elias/axiom/engine/utils"
)

const (
	heatRate    = 0.05
	hvacRate    = 0.85
	maxTempPush = .035

	minLivableTemp = -10.0
	maxLivableTemp = 50.0

	hvacEfficiencyZeroThreshold   = 0.0
	hvacEfficiencyMediumThreshold = 0.5
	hvacEfficiencyHighThreshold   = .9
)

type HvacInput struct {
	PowerSupplied float64
	HeatSources   []float64
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

	targetTemperature float64
	powerCapacity     float64
}

func NewHvac(targetTemperature float64) *Hvac {

	// NOTE: Linear interpolation only works while temp curve is linear!!
	normTargetTemp := (targetTemperature - minLivableTemp) / (maxLivableTemp - minLivableTemp)

	system := &Hvac{
		SystemCore:        systems.NewSystemCore("HVAC"),
		temperature:       components.NewComponent("Bunker Temperature (C)", normTargetTemp, minLivableTemp, maxLivableTemp),
		health:            components.NewHealthComponent(1.0),
		targetTemperature: targetTemperature,
		powerCapacity:     1200,
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
	effectiveness := calculateEffectiveness(input.PowerSupplied, s.powerCapacity)

	heatPush := 0.0
	if len(input.HeatSources) > 0 {
		averageHeat := calculateAverageHeat(input.HeatSources)
		heatPush = calculateHeatPush(averageHeat, s.temperature.Value(), heatRate)
	}

	hvacPush := calculateHeatPush(s.targetTemperature, s.temperature.Value(), effectiveness*hvacRate)

	netPush := (heatPush + hvacPush) / (s.temperature.Max() - s.temperature.Min())
	netPush = utils.Clamp(-maxTempPush, netPush, maxTempPush)

	s.temperature.SetNorm(s.temperature.Norm() + netPush)

	output := HvacOutput{
		Status:            s.Status(),
		TargetTemperature: s.targetTemperature,
		Temperature:       s.temperature.Value(),
	}

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

func calculateEffectiveness(powerSupplied, powerCapacity float64) float64 {
	if powerCapacity == 0 {
		return 0
	}

	ratio := powerSupplied / powerCapacity
	if ratio < 0 {
		return 0
	}

	eff := math.Pow(ratio, 2)
	return utils.Clamp(0, eff, 1)
}

func (s *Hvac) String() string {
	output := fmt.Sprintf("%v: %v", s.ID(), s.Name())
	output += fmt.Sprintf("\n%s", s.temperature)
	output += fmt.Sprintf("\n%s", s.health)

	return output
}
