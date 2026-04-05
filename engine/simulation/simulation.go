package simulation

import (
	"fmt"

	"github.com/elias/axiom/engine"
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/subsystems"
	"github.com/elias/axiom/engine/subsystems/connections"
	"github.com/elias/axiom/engine/utils"
)

type WorldState struct {
	subsystems  map[subsystems.SubsystemID]subsystems.Subsystem
	connections map[subsystems.SubsystemID][]*connections.Connection
}

func (ws *WorldState) addSubsystem(subsystem subsystems.Subsystem) {
	ws.subsystems[subsystem.ID()] = subsystem
	ws.connections[subsystem.ID()] = []*connections.Connection{}

}

func (ws *WorldState) addConnection(src subsystems.Subsystem, srcPortName string, dest subsystems.Subsystem, destPortName string, throughput utils.Unit) error {
	srcPort, exists := src.Ports()[srcPortName]
	if !exists {
		return fmt.Errorf("Port %s doesn't exist on subsystem %s", srcPortName, src.Name())
	}

	destPort, exists := dest.Ports()[destPortName]
	if !exists {
		return fmt.Errorf("Port %s doesn't exist on subsystem %s", destPortName, dest.Name())
	}

	connection := connections.NewConnection(srcPort, destPort, throughput)

	ws.connections[dest.ID()] = append(ws.connections[dest.ID()], connection)

	return nil
}

func (ws *WorldState) Init() {

	ws.subsystems = make(map[subsystems.SubsystemID]subsystems.Subsystem)
	ws.connections = make(map[subsystems.SubsystemID][]*connections.Connection)

	reactor := subsystems.NewPower(.5)
	cooler := subsystems.NewCooling(.5)
	acUnit := subsystems.NewHvac()

	ws.addSubsystem(reactor)
	ws.addSubsystem(cooler)
	ws.addSubsystem(acUnit)

	reactor.AddPort("socket-1", "power-out")
	reactor.AddPort("valve-1", "temp-in")
	reactor.AddPort("valve-2", "temp-out")

	acUnit.AddPort("socket-1", "power-in")
	acUnit.AddPort("valve-1", "temp-in")

	cooler.AddPort("valve-1", "temp-out")

	err := ws.addConnection(reactor, "socket-1", acUnit, "socket-1", 0.5)
	if err != nil {
		logging.Error(err.Error())
	}

	err = ws.addConnection(reactor, "valve-1", acUnit, "valve-1", 0.5)
	if err != nil {
		logging.Error(err.Error())
	}

	err = ws.addConnection(cooler, "valve-1", reactor, "valve-2", 1)
	if err != nil {
		logging.Error(err.Error())
	}
}

func (ws *WorldState) Update(tick *engine.Tick) {

	ws.updateSubsystems()
	for systemID := range len(ws.subsystems) {
		logging.Info(ws.subsystems[subsystems.SubsystemID(systemID)].String())
	}
}

func (ws *WorldState) updateSubsystems() {
	visited := make(map[subsystems.SubsystemID]struct{})

	depStack := utils.NewStack[subsystems.Subsystem]()

	// DFS
	for _, system := range ws.subsystems {
		depStack.Push(system)
		for depStack.Len() > 0 {
			subsystem := depStack.Pop()
			if _, seen := visited[subsystem.ID()]; seen {
				continue
			}

			visited[subsystem.ID()] = struct{}{}
			if len(ws.connections[subsystem.ID()]) <= 0 {
				subsystem.Tick()
			}

			for _, conn := range ws.connections[subsystem.ID()] {
				src := conn.SrcPort().Subsystem()
				if _, seen := visited[src.ID()]; !seen {
					subsystem := ws.subsystems[src.ID()]
					depStack.Push(subsystem)

				}
			}
		}

		for _, conn := range ws.connections[system.ID()] {
			srcComp := *conn.SrcPort().Component()
			srcComp.SetValue(srcComp.Value() * conn.Throughput())
		}
		system.Tick()

	}

}
