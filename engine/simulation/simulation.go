package simulation

import (
	"github.com/elias/axiom/engine"
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/subsystems"
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/subsystems/connections"
	"github.com/elias/axiom/engine/utils"
)

type WorldState struct {
	subsystems  map[subsystems.SubsystemID]subsystems.Subsystem
	connections map[subsystems.SubsystemID][]*connections.Connection
	ports       map[connections.PortID]subsystems.SubsystemID
}

func (ws *WorldState) addSubsystem(subsystem subsystems.Subsystem) {
	ws.subsystems[subsystem.ID()] = subsystem
	ws.connections[subsystem.ID()] = []*connections.Connection{}

}

func (ws *WorldState) AddPort(name string, subsystemID subsystems.SubsystemID, component *components.Component) *connections.Port {
	subsystem, exists := ws.subsystems[subsystemID]
	if !exists {
		return nil
	}

	if _, exists := subsystem.Components()[component.Name()]; !exists {
		return nil
	}

	port := connections.NewPort(name, component)

	ws.ports[port.ID()] = subsystemID

	return port
}

func (ws *WorldState) addConnection(src *connections.Port, subsystemID subsystems.SubsystemID, throughput utils.Norm) {
	connection := connections.NewConnection(src, subsystemID, throughput)

	_, exists := ws.ports[src.ID()]
	if !exists {
		return
	}

	ws.connections[subsystemID] = append(ws.connections[subsystemID], connection)
}

func (ws *WorldState) Init() {

	ws.subsystems = make(map[subsystems.SubsystemID]subsystems.Subsystem)
	ws.connections = make(map[subsystems.SubsystemID][]*connections.Connection)
	ws.ports = make(map[connections.PortID]subsystems.SubsystemID)

	power := subsystems.NewPower(.5)
	cooling := subsystems.NewCooling(.5)
	hvac := subsystems.NewHvac()

	ws.addSubsystem(power)
	ws.addSubsystem(cooling)
	ws.addSubsystem(hvac)

	powerPort := ws.AddPort("socket-1", power.ID(), power.Components()["power"])
	powerTemp := ws.AddPort("valve-1", power.ID(), power.Components()["temp"])
	coolingTempPort := ws.AddPort("valve-1", cooling.ID(), cooling.Components()["temp-out"])
	coolingFlowPort := ws.AddPort("valve-2", cooling.ID(), cooling.Components()["flow-out"])

	ws.addConnection(powerPort, hvac.ID(), .5)
	ws.addConnection(powerTemp, hvac.ID(), 1)
	ws.addConnection(coolingFlowPort, power.ID(), 1)
	ws.addConnection(coolingTempPort, power.ID(), 1)
}

// updates the world state
func (ws *WorldState) Update(tick *engine.Tick) {

	ws.updateSubsystems()
	for systemID := range len(ws.subsystems) {
		logging.Info(ws.subsystems[subsystems.SubsystemID(systemID)].String())
	}
}

// iterates through the connection dependency tree for subsystems using DFS
func (ws *WorldState) updateSubsystems() {
	visited := make(map[subsystems.SubsystemID]struct{})

	depStack := utils.NewStack[subsystems.Subsystem]()

	for _, system := range ws.subsystems {
		depStack.Push(system)
		for depStack.Len() > 0 {
			subsystem := depStack.Pop()
			if _, seen := visited[subsystem.ID()]; seen {
				continue
			}

			visited[subsystem.ID()] = struct{}{}
			if len(ws.connections[subsystem.ID()]) <= 0 {
				subsystem.Tick(nil)
			}

			for _, conn := range ws.connections[subsystem.ID()] {
				srcID := ws.ports[conn.Src().ID()]
				if _, seen := visited[srcID]; !seen {
					subsystem := ws.subsystems[srcID]
					depStack.Push(subsystem)

				}
			}
		}

		inputs := make(map[string]components.Component, 0)
		for _, conn := range ws.connections[system.ID()] {
			srcComp := *conn.Src().Component()
			srcComp.SetValue(srcComp.Value() * conn.Throughput())
			inputs[srcComp.Name()] = srcComp
		}
		system.Tick(inputs)

	}

}
