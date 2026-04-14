package simulation

import (
	"fmt"

	"github.com/elias/axiom/engine/subsystems"
	"github.com/elias/axiom/engine/subsystems/connections"
	"github.com/elias/axiom/engine/utils"
)

type WorldState struct {
	currentSubsystemID subsystems.SubsystemID

	subsystems  map[string]subsystems.Subsystem
	connections map[subsystems.SubsystemID][]*connections.Connection
}

func (ws *WorldState) newSubsystemID() subsystems.SubsystemID {
	id := ws.currentSubsystemID
	ws.currentSubsystemID++
	return id
}

func (ws *WorldState) newSubsystem(id subsystems.SubsystemID, name, subsystemType string) (subsystems.Subsystem, error) {
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
	ws.connections[subsystem.ID()] = []*connections.Connection{}

	return nil
}

func (ws *WorldState) addConnection(src *subsystems.OutputPort, dest *subsystems.InputPort, throughput utils.Unit) {
	connection := connections.NewConnection(src, dest, throughput)
	ws.connections[dest.Subsystem().ID()] = append(ws.connections[dest.Subsystem().ID()], connection)
}

func (ws *WorldState) Init() {
	ws.subsystems = make(map[string]subsystems.Subsystem)
	ws.connections = make(map[subsystems.SubsystemID][]*connections.Connection)
}
