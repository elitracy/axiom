package systems

type LifeSupportSystem struct {
	*SubsystemCore
	o2           float64 // % of atmo
	co2          float64
	scrubberRate float64

	o2ConsumptionRate float64
}

func NewLifeSupportSystem() *LifeSupportSystem {
	return &LifeSupportSystem{
		SubsystemCore:     NewSubsystem("Life Support"),
		o2:                .18,
		co2:               .04,
		scrubberRate:      .005,
		o2ConsumptionRate: .02, // 100 ticks before problem
	}
}

func (s *LifeSupportSystem) Tick(powerAvailable float64) {
	s.o2 -= s.o2ConsumptionRate
	s.o2 += s.scrubberRate * powerAvailable

	s.co2 += s.o2ConsumptionRate * 1 / 3 // co2 created at a third the rate of o2
	s.co2 -= s.scrubberRate * powerAvailable * 1 / 3

	s.health -= s.degradationRate

	s.sensorValues["o2"] = s.o2
	s.sensorValues["co2"] = s.co2
	s.sensorValues["scrubber_status"] = s.scrubberRate
}
