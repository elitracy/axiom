package state

import (
	"fmt"
	"slices"
	"sync"

	"github.com/elias/axiom/engine/subsystems"
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/subsystems/connections"
	"github.com/elias/axiom/engine/telemetry"
	"github.com/elias/axiom/engine/utils"
)

type Subsystem interface {
	ID() subsystems.SubsystemID
	Name() utils.SubsystemName
	Type() utils.SubsystemType
	Components() map[string]*components.Component
	OutputPorts() map[string]*subsystems.OutputPort
	InputPorts() map[string]*subsystems.InputPort
	String() string
	ExportFields() *telemetry.Export
	Tick()
	Status() utils.Status
}

type State struct {
	currentSubsystemID  subsystems.SubsystemID
	currentConnectionID connections.ConnectionID

	subsystems  map[utils.SubsystemName]Subsystem
	connections map[utils.SubsystemName]map[utils.PortType][]*connections.Connection

	mu sync.RWMutex
}

func NewState() *State {
	return &State{
		subsystems:  make(map[utils.SubsystemName]Subsystem),
		connections: make(map[utils.SubsystemName]map[utils.PortType][]*connections.Connection),
	}
}

func (s *State) newSubsystemID() subsystems.SubsystemID {
	id := s.currentSubsystemID
	s.currentSubsystemID++
	return id
}

func (s *State) newConnectionID() connections.ConnectionID {
	id := s.currentConnectionID
	s.currentConnectionID++
	return id
}

func (s *State) newSubsystem(id subsystems.SubsystemID, name utils.SubsystemName, subsystemType utils.SubsystemType) (Subsystem, error) {
	switch subsystemType {
	case utils.Power:
		return subsystems.NewPower(id, name, 0.5), nil
	case utils.Cooling:
		return subsystems.NewCooling(id, name, 0.5), nil
	case utils.Hvac:
		return subsystems.NewHvac(id, name, 0.2), nil
	default:
		return nil, fmt.Errorf("unknown subsystem type: %s", subsystemType)
	}
}

func (s *State) addSubsystem(name utils.SubsystemName, subsystemType utils.SubsystemType) error {
	id := s.newSubsystemID()

	subsystem, err := s.newSubsystem(id, name, subsystemType)
	if err != nil {
		return err
	}

	s.subsystems[subsystem.Name()] = subsystem
	s.connections[subsystem.Name()] = make(map[utils.PortType][]*connections.Connection)

	return nil
}

func (s *State) addConnection(src *subsystems.OutputPort, dest *subsystems.InputPort, srcSystem, destSystem utils.SubsystemName, throughput utils.Unit) {
	id := s.newConnectionID()

	connection := connections.NewConnection(id, src, dest, srcSystem, destSystem, throughput)

	s.connections[srcSystem][utils.PortOutput] = append(s.connections[srcSystem][utils.PortOutput], connection)
	s.connections[destSystem][utils.PortInput] = append(s.connections[destSystem][utils.PortInput], connection)
}

func (s *State) GetSubsystem(name utils.SubsystemName) (Subsystem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if subsystem, exists := s.subsystems[name]; exists {
		return subsystem, nil
	}
	return nil, fmt.Errorf("Subsystem not found %s", name)
}

func (s *State) Subsystems() []Subsystem {

	keys := make([]utils.SubsystemName, 0, len(s.subsystems))
	for k := range s.subsystems {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	sortedSubsystems := make([]Subsystem, 0, len(s.subsystems))
	for _, key := range keys {
		sortedSubsystems = append(sortedSubsystems, s.subsystems[key])
	}

	return sortedSubsystems
}

func (s *State) Connections() map[utils.SubsystemName]map[utils.PortType][]*connections.Connection {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.connections
}
