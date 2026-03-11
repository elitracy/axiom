package systems

import "fmt"

type SystemStatus int8

const (
	Offline SystemStatus = iota
	Critical
	Degraded
	Online
)

type Subsystem interface {
	ID() string
	Name() string
	Health() float64
	Status() SystemStatus
	Sensors() map[string]float64
	DegradationRate() float64
	String() string
}

type SubsystemCore struct {
	id              string
	name            string
	sensorValues    map[string]float64
	health          float64
	degradationRate float64
}

func NewSubsystem(name string) *SubsystemCore {
	return &SubsystemCore{
		id:              "",
		name:            name,
		sensorValues:    make(map[string]float64),
		health:          1.0,
		degradationRate: 0.01,
	}
}

func (s *SubsystemCore) ID() string                  { return s.id }
func (s *SubsystemCore) Name() string                { return s.name }
func (s *SubsystemCore) Health() float64             { return s.health }
func (s *SubsystemCore) Sensors() map[string]float64 { return s.sensorValues }
func (s *SubsystemCore) DegradationRate() float64    { return s.degradationRate }
func (s *SubsystemCore) Status() SystemStatus {
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

func (s *SubsystemCore) String() string {
	output := ""

	output += fmt.Sprintf("%v: %.2f%%\n", s.name, s.health*100)
	output += fmt.Sprintf("%v", s.Sensors())

	return output
}
