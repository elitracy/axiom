package systems

import "github.com/elias/axiom/engine/utils"

const (
	startingO2         = 0.21
	startingCO2        = 0.01
	scrubberEffeciency = 0.9
	co2Ratio           = 0.333
)

type LifeSupportSystem struct {
	*SubsystemCore
	o2           float64 // % of atmo
	co2          float64
	scrubberRate float64

	o2ConsumptionRate float64
}

func NewLifeSupportSystem() *LifeSupportSystem {
	system := &LifeSupportSystem{
		SubsystemCore:     NewSubsystem("Life Support"),
		o2:                startingO2,
		co2:               startingCO2,
		scrubberRate:      1.0,
		o2ConsumptionRate: .005, // 100 ticks before problem
	}
	system.Sensors()["o2"] = startingO2
	system.Sensors()["co2"] = startingCO2

	return system
}

func (s *LifeSupportSystem) Tick(powerAvailable float64) {
	s.scrubberRate = s.Health()

	s.o2 -= s.o2ConsumptionRate
	s.co2 += s.o2ConsumptionRate * co2Ratio

	s.o2 += s.o2ConsumptionRate * scrubberEffeciency * s.scrubberRate * powerAvailable
	s.co2 -= s.o2ConsumptionRate * scrubberEffeciency * co2Ratio * s.scrubberRate * powerAvailable

	s.health -= s.degradationRate

	s.health = utils.Clamp(s.health, 0, 1)
	s.o2 = utils.Clamp(s.o2, 0, 1)
	s.co2 = utils.Clamp(s.co2, 0, 1)

	s.sensorValues["o2"] = s.o2
	s.sensorValues["co2"] = s.co2
	s.sensorValues["scrubber_status"] = s.scrubberRate
}
