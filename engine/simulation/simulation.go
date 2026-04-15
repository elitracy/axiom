package simulation

import (
	"fmt"

	"github.com/elias/axiom/engine/subsystems"
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/subsystems/connections"
	"github.com/elias/axiom/engine/utils"
)

type subsystem interface {
	ID() subsystems.SubsystemID
	Name() string
	Components() map[string]*components.Component
	OutputPorts() map[string]*subsystems.OutputPort
	InputPorts() map[string]*subsystems.InputPort
	String() string
	Tick()
}

type WorldState struct {
	currentSubsystemID  subsystems.SubsystemID
	currentConnectionID connections.ConnectionID

	subsystems  map[string]subsystem
	connections map[string][]*connections.Connection
}

func (ws *WorldState) newSubsystemID() subsystems.SubsystemID {
	id := ws.currentSubsystemID
	ws.currentSubsystemID++
	return id
}

func (ws *WorldState) newConnectionID() connections.ConnectionID {
	id := ws.currentConnectionID
	ws.currentConnectionID++
	return id
}

func (ws *WorldState) newSubsystem(id subsystems.SubsystemID, name, subsystemType string) (subsystem, error) {
	switch subsystemType {
	case "power":
		return subsystems.NewPower(id, name, 0.5), nil
	case "cooling":
		return subsystems.NewCooling(id, name, 0.5), nil
	case "hvac":
		return subsystems.NewHvac(id, name, 0.2), nil
	default:
		return nil, fmt.Errorf("unknown subsystem type: %s", subsystemType)
	}
}

func (ws *WorldState) addSubsystem(name, subsystemType string) error {
	id := ws.newSubsystemID()

	subsystem, err := ws.newSubsystem(id, name, subsystemType)
	if err != nil {
		return err
	}

	ws.subsystems[subsystem.Name()] = subsystem
	ws.connections[subsystem.Name()] = []*connections.Connection{}

	return nil
}

func (ws *WorldState) addConnection(src *subsystems.OutputPort, dest *subsystems.InputPort, srcSystem string, destSystem string, throughput utils.Unit) {
	id := ws.newConnectionID()

	connection := connections.NewConnection(id, src, dest, srcSystem, destSystem, throughput)

	ws.connections[destSystem] = append(ws.connections[destSystem], connection)
}

func (ws *WorldState) Init() {
	ws.subsystems = make(map[string]subsystem)
	ws.connections = make(map[string][]*connections.Connection)
}
