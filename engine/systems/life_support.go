package systems

type LifeSupportSystem struct {
	*SubsystemCore
	o2           float64
	co2          float64
	scrubberRate float64

	o2ConsumptionRate float64
}

func (s *LifeSupportSystem) Tick(powerAvailable float64) {
	s.o2 -= s.o2ConsumptionRate
	s.o2 += s.scrubberRate * powerAvailable

	s.co2 += s.o2ConsumptionRate
	s.co2 -= s.scrubberRate * powerAvailable

	s.health -= s.degradationRate

	s.sensorValues["o2"] = s.o2
	s.sensorValues["co2"] = s.co2
	s.sensorValues["scrubber_status"] = s.scrubberRate
}
