package systems

import (
	"math"

	"github.com/elias/axiom/engine/utils"
)

const (
	startingO2  = 0.21
	startingCO2 = 0.01
	co2Ratio    = 0.333
)

type LifeSupportSystem struct {
	*SystemCore
	o2                 float64 // % of atmo
	co2                float64
	scrubberEfficiency float64

	o2ConsumptionRate float64
}

func NewLifeSupportSystem() *LifeSupportSystem {
	system := &LifeSupportSystem{
		SystemCore:      NewSystem("Life Support"),
		o2:                 startingO2,
		co2:                startingCO2,
		scrubberEfficiency: 1.0,
		o2ConsumptionRate:  .005, // 100 ticks before problem
	}
	system.Sensors()["o2"] = startingO2
	system.Sensors()["co2"] = startingCO2

	return system
}

func (s *LifeSupportSystem) Tick(powerAvailable float64) {
	s.scrubberEfficiency = math.Tanh(s.Health() * 3.8)

	s.o2 -= s.o2ConsumptionRate
	s.co2 += s.o2ConsumptionRate * co2Ratio

	s.o2 += s.o2ConsumptionRate * s.scrubberEfficiency * powerAvailable
	s.co2 -= s.o2ConsumptionRate * co2Ratio * s.scrubberEfficiency * powerAvailable

	s.health -= s.degradationRate

	s.health = utils.Clamp(s.health, 0, 1)
	s.o2 = utils.Clamp(s.o2, 0, 1)
	s.co2 = utils.Clamp(s.co2, 0, 1)

	s.sensorValues["o2"] = s.o2
	s.sensorValues["co2"] = s.co2
	s.sensorValues["scrubber_efficiency"] = s.scrubberEfficiency
}
