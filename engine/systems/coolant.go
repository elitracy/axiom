package systems

import "github.com/elias/axiom/engine/utils"

type CoolantSystem struct {
	*SubsystemCore
	pressure float64 // psi
	flowRate float64 // liters/tick
	temp     float64 // celsius

	maxPressure       float64
	pressureDeltaRate float64
	maxTempDelta      float64
	coolingRate       float64
}

func NewCoolantSystem(ambientTemp float64) *CoolantSystem {
	return &CoolantSystem{
		SubsystemCore: NewSubsystem("Coolant"),
		temp:          ambientTemp,
		pressure:      30,

		maxPressure:       100,
		pressureDeltaRate: -0.05,

		maxTempDelta: 300, // matches powers max temp
		coolingRate:  8,
	}
}

func (s *CoolantSystem) Tick(heatLoad float64) {
	s.health -= s.degradationRate
	s.flowRate = utils.Sigmoid(s.pressure, 12.5, .12)
	s.pressure += (s.flowRate * s.pressureDeltaRate)

	tempDelta := heatLoad - s.temp
	transferRate := utils.Sigmoid(s.flowRate, 0.5, 5) * s.coolingRate * (tempDelta / s.maxTempDelta)
	s.temp += transferRate

	s.pressure = utils.Clamp(s.pressure, 0, s.maxPressure)
	s.temp = min(s.temp, heatLoad)
	s.flowRate = max(0, s.flowRate)

	s.sensorValues["flow_rate"] = s.flowRate
	s.sensorValues["pressure"] = s.pressure
	s.sensorValues["temp"] = s.temp

}
