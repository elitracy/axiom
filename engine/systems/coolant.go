package systems

import "github.com/elias/axiom/engine/utils"

type CoolantSystem struct {
	*SubsystemCore
	pressure float64
	flowRate float64
	temp     float64

	maxPressure      float64 //these could be later adjusted by the quality of pump or pipe
	pressureLeakRate float64
	pressureGainRate float64
	maxFlowRate      float64
	maxTempDelta     float64
	coolingRate      float64
}

func NewCoolantSystem(ambientTemp float64) *CoolantSystem {
	return &CoolantSystem{
		SubsystemCore: NewSubsystem("Coolant"),
		temp:          ambientTemp,
		pressure:      30,

		maxPressure:      100,
		pressureLeakRate: 2.5,
		pressureGainRate: 2,

		maxFlowRate: 1,

		maxTempDelta: 500, // matches powers max temp
		coolingRate:  6,
	}
}

func (s *CoolantSystem) Tick(heatLoad float64) {
	s.health -= s.degradationRate
	s.flowRate = s.maxFlowRate * (s.health / 100) * (1 - s.pressure/s.maxPressure)
	s.pressure += (s.flowRate * s.pressureGainRate) - s.pressureLeakRate
	s.pressure = utils.Clamp(s.pressure, 0, s.maxPressure)

	tempDelta := heatLoad - s.temp
	transferRate := s.flowRate * s.coolingRate * (tempDelta / s.maxTempDelta)
	s.temp += transferRate
	s.temp = min(s.temp, heatLoad)

	s.sensorValues["flow_rate"] = s.flowRate
	s.sensorValues["pressure"] = s.pressure
	s.sensorValues["temp"] = s.temp

}
