package simulation

import (
	"github.com/elias/axiom/engine/subsystems"
	"github.com/elias/axiom/engine/subsystems/connections"
	"github.com/elias/axiom/engine/utils"
)

type WorldState struct {
	subsystems  map[string]subsystems.Subsystem
	connections map[subsystems.SubsystemID][]*connections.Connection
}

func (ws *WorldState) addSubsystem(subsystem subsystems.Subsystem) {
	ws.subsystems[subsystem.Name()] = subsystem
	ws.connections[subsystem.ID()] = []*connections.Connection{}

}

func (ws *WorldState) addConnection(src *subsystems.OutputPort, dest *subsystems.InputPort, throughput utils.Unit) {
	connection := connections.NewConnection(src, dest, throughput)
	ws.connections[dest.Subsystem().ID()] = append(ws.connections[dest.Subsystem().ID()], connection)
}

func (ws *WorldState) Init() {
	ws.subsystems = make(map[string]subsystems.Subsystem)
	ws.connections = make(map[subsystems.SubsystemID][]*connections.Connection)
}
