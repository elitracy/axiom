package state

import (
	"fmt"
	"slices"

	"github.com/elias/axiom/engine/subsystems"
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/subsystems/connections"
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
	Tick()
	Status() utils.Status
}

type State struct {
	currentSubsystemID  subsystems.SubsystemID
	currentConnectionID connections.ConnectionID

	subsystems  map[utils.SubsystemName]Subsystem
	connections map[utils.SubsystemName]map[utils.PortType][]*connections.Connection
}

func NewState() *State {
	return &State{
		subsystems:  make(map[utils.SubsystemName]Subsystem),
		connections: make(map[utils.SubsystemName]map[utils.PortType][]*connections.Connection),
	}
}

func (ws *State) newSubsystemID() subsystems.SubsystemID {
	id := ws.currentSubsystemID
	ws.currentSubsystemID++
	return id
}

func (ws *State) newConnectionID() connections.ConnectionID {
	id := ws.currentConnectionID
	ws.currentConnectionID++
	return id
}

func (ws *State) newSubsystem(id subsystems.SubsystemID, name utils.SubsystemName, subsystemType utils.SubsystemType) (Subsystem, error) {
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

func (ws *State) addSubsystem(name utils.SubsystemName, subsystemType utils.SubsystemType) error {
	id := ws.newSubsystemID()

	subsystem, err := ws.newSubsystem(id, name, subsystemType)
	if err != nil {
		return err
	}

	ws.subsystems[subsystem.Name()] = subsystem
	ws.connections[subsystem.Name()] = make(map[utils.PortType][]*connections.Connection)

	return nil
}

func (ws *State) addConnection(src *subsystems.OutputPort, dest *subsystems.InputPort, srcSystem, destSystem utils.SubsystemName, throughput utils.Unit) {
	id := ws.newConnectionID()

	connection := connections.NewConnection(id, src, dest, srcSystem, destSystem, throughput)

	ws.connections[srcSystem][utils.PortOutput] = append(ws.connections[srcSystem][utils.PortOutput], connection)
	ws.connections[destSystem][utils.PortInput] = append(ws.connections[destSystem][utils.PortInput], connection)
}

func (ws State) GetSubsystem(name utils.SubsystemName) (Subsystem, error) {
	if subsystem, exists := ws.subsystems[name]; exists {
		return subsystem, nil
	}
	return nil, fmt.Errorf("Subsystem not found %s", name)
}

func (ws *State) Subsystems() []Subsystem {
	keys := make([]utils.SubsystemName, 0, len(ws.subsystems))
	for k := range ws.subsystems {
		keys = append(keys, k)
	}

	slices.Sort(keys)

	sortedSubsystems := make([]Subsystem, 0, len(ws.subsystems))
	for _, key := range keys {
		sortedSubsystems = append(sortedSubsystems, ws.subsystems[key])
	}

	return sortedSubsystems
}

func (ws *State) Connections() map[utils.SubsystemName]map[utils.PortType][]*connections.Connection {
	return ws.connections
}
