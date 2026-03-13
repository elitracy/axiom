package systems

import (
	"github.com/elias/axiom/engine/utils"
)

type PowerSystem struct {
	*SystemCore
	fuelLevel           float64 // % of tank
	fuelConsumptionRate float64
	outputLevel         float64
	heatRate            float64 // celsius/tick
	coolRate            float64
	temp                float64
	maxTemp             float64
}

func NewPowerSystem(ambientTemp float64) *PowerSystem {
	system := &PowerSystem{
		SystemCore:       NewSystem("Power"),
		fuelLevel:           1.0,
		fuelConsumptionRate: 0.01,
		outputLevel:         1.0,
		heatRate:            1.1,
		temp:                ambientTemp,
		maxTemp:             300,
	}

	system.Sensors()["output_level"] = 1.0
	system.Sensors()["fuel_level"] = 1.0
	system.Sensors()["temp"] = ambientTemp

	return system
}

func (s *PowerSystem) SetOutputLevel(level float64) { s.outputLevel = level }

func (s *PowerSystem) Tick(coolRate, ambientTemp float64) {

	heatGenerated := s.outputLevel * s.heatRate
	heatDissipated := s.outputLevel * coolRate
	s.temp += heatGenerated - heatDissipated
	s.temp = utils.Clamp(s.temp, 0, s.maxTemp)

	s.health -= s.degradationRate
	s.outputLevel = utils.Sigmoid(s.health, 0.3, 9.5)
	s.fuelLevel -= s.fuelConsumptionRate * s.outputLevel

	s.health = utils.Clamp(s.health, 0, 1)
	s.outputLevel = max(0, s.outputLevel)
	s.fuelLevel = utils.Clamp(s.fuelLevel, 0, 1)

	s.sensorValues["temp"] = s.temp
	s.sensorValues["output_level"] = s.outputLevel
	s.sensorValues["fuel_level"] = s.fuelLevel
}
