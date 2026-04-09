package simulation

import (
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

func (ws *WorldState) addConnection(src *subsystems.OutputPort, dest *subsystems.InputPort, throughput utils.Unit) error {
	connection := connections.NewConnection(src, dest, throughput)

	ws.connections[dest.Subsystem().ID()] = append(ws.connections[dest.Subsystem().ID()], connection)

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

	reactor.AddPort("socket-1", "power-out", subsystems.Output)
	reactor.AddPort("valve-1", "temp-in", subsystems.Input)
	reactor.AddPort("valve-2", "temp-out", subsystems.Output)

	acUnit.AddPort("socket-1", "power-in", subsystems.Input)
	acUnit.AddPort("valve-1", "temp-in", subsystems.Input)

	cooler.AddPort("valve-1", "temp-out", subsystems.Output)

	err := ws.addConnection(reactor.OutputPorts()["socket-1"], acUnit.InputPorts()["socket-1"], 0.5)
	if err != nil {
		logging.Error(err.Error())
	}

	err = ws.addConnection(reactor.OutputPorts()["valve-2"], acUnit.InputPorts()["valve-1"], 0.5)
	if err != nil {
		logging.Error(err.Error())
	}

	err = ws.addConnection(cooler.OutputPorts()["valve-1"], reactor.InputPorts()["valve-1"], 1)
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
				src := conn.Src().Subsystem()
				if _, seen := visited[src.ID()]; !seen {
					subsystem := ws.subsystems[src.ID()]
					depStack.Push(subsystem)

				}
			}
		}

		for _, conn := range ws.connections[system.ID()] {
			srcComp := *conn.Src().Component()
			conn.Dest().SetInput(srcComp.Value() * conn.Throughput())
		}
		system.Tick()

	}

}
