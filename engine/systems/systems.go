package systems

import "fmt"

type SystemStatus int8

const (
	Offline SystemStatus = iota
	Critical
	Degraded
	Online
)

type System interface {
	ID() string
	Name() string
	Health() float64
	Status() SystemStatus
	Sensors() map[string]float64
	DegradationRate() float64
	String() string
}

type SystemCore struct {
	id              string
	name            string
	sensorValues    map[string]float64
	health          float64
	degradationRate float64
}

func NewSystem(name string) *SystemCore {
	return &SystemCore{
		id:              "",
		name:            name,
		sensorValues:    make(map[string]float64),
		health:          1.0,
		degradationRate: 0.01,
	}
}

func (s *SystemCore) ID() string                  { return s.id }
func (s *SystemCore) Name() string                { return s.name }
func (s *SystemCore) Health() float64             { return s.health }
func (s *SystemCore) Sensors() map[string]float64 { return s.sensorValues }
func (s *SystemCore) DegradationRate() float64    { return s.degradationRate }
func (s *SystemCore) Status() SystemStatus {
	switch {
	case s.health <= 0:
		return Offline
	case s.health <= .30:
		return Critical
	case s.health <= .70:
		return Degraded
	default:
		return Online
	}
}

func (s *SystemCore) String() string {
	output := ""

	output += fmt.Sprintf("%v: %.2f%%\n", s.name, s.health*100)
	output += fmt.Sprintf("%v", s.Sensors())

	return output
}
