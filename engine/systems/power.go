package systems

import (
	"github.com/elias/axiom/engine/utils"
)

type PowerSystem struct {
	*SubsystemCore
	fuelLevel           float64 // % of tank
	fuelConsumptionRate float64
	outputLevel         float64
	heatRate            float64 // celsius/tick
	coolRate            float64
	temp                float64
	maxTemp             float64
}

func NewPowerSystem(ambientTemp float64) *PowerSystem {
	return &PowerSystem{
		SubsystemCore:       NewSubsystem("Power"),
		fuelLevel:           1.0,
		fuelConsumptionRate: 0.01,
		outputLevel:         1.0,
		coolRate:            5,
		heatRate:            5,
		temp:                ambientTemp,
		maxTemp:             500,
	}
}

func (s *PowerSystem) SetOutputLevel(level float64) { s.outputLevel = level }

func (s *PowerSystem) Tick(coolantFlowRate, ambientTemp float64) {

	heatGenerated := s.outputLevel * s.heatRate
	heatDissipated := coolantFlowRate * s.coolRate
	s.temp += heatGenerated - heatDissipated
	s.temp = utils.Clamp(s.temp, 0, s.maxTemp)

	s.health -= s.degradationRate
	s.fuelLevel -= s.fuelConsumptionRate * s.outputLevel

	s.sensorValues["temp"] = s.temp
	s.sensorValues["output_level"] = s.outputLevel
	s.sensorValues["fuel_level"] = s.fuelLevel
}
