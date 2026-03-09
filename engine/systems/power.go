package systems

import (
	"github.com/elias/axiom/engine/utils"
)

type PowerSystem struct {
	*SubsystemCore
	fuelLevel           float64
	fuelConsumptionRate float64
	outputLevel         float64 // 0-1
	heatRate            float64
	coolRate            float64
	temp                float64
	maxTemp             float64
}

func NewPowerSystem(ambientTemp float64) *PowerSystem {
	return &PowerSystem{
		SubsystemCore:       NewSubsystem("Power"),
		fuelLevel:           100,
		fuelConsumptionRate: 1, // run out of fuel at 100 ticks
		outputLevel:         1,
		coolRate:            5, // neutral temp rate
		heatRate:            5, // overheat at 100 ticks
		temp:                ambientTemp,
		maxTemp:             500,
	}
}

func (s *PowerSystem) SetOutputLevel(level float64) { s.outputLevel = level }

func (s *PowerSystem) Tick(coolantFlow, ambientTemp float64) {

	heatGenerated := s.outputLevel * s.heatRate
	heatDissipated := coolantFlow * s.coolRate
	s.temp += heatGenerated - heatDissipated
	s.temp = utils.Clamp(s.temp, 0, s.maxTemp)

	s.health -= s.degradationRate
	s.fuelLevel -= s.fuelConsumptionRate * s.outputLevel

	s.sensorValues["temp"] = s.temp
	s.sensorValues["output_level"] = s.outputLevel
	s.sensorValues["fuel_level"] = s.fuelLevel
}
