package simulation

import (
	"fmt"

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
	ports       map[connections.PortID]*connections.Port
}

func (ws *WorldState) addSubsystem(subsystem subsystems.Subsystem) {
	ws.subsystems[subsystem.ID()] = subsystem
	ws.connections[subsystem.ID()] = []*connections.Connection{}

}

func (ws *WorldState) AddPort(name string, subsystem subsystems.Subsystem, component string) *connections.Port {
	comp, exists := subsystem.Components()[component]
	if !exists {
		return nil
	}

	port := connections.NewPort(name, comp, subsystem)
	ws.ports[port.ID()] = port

	return port
}

func (ws *WorldState) addConnection(src *connections.Port, dest *connections.Port, throughput utils.Unit) error {
	connection := connections.NewConnection(src, dest, throughput)

	_, srcExists := ws.ports[src.ID()]
	if !srcExists {
		return fmt.Errorf("src port doesn't exist: %v", src)
	}

	_, destExists := ws.ports[dest.ID()]
	if !destExists {
		return fmt.Errorf("dest port doesn't exist: %v", dest)
	}

	ws.connections[dest.Subsystem().ID()] = append(ws.connections[dest.Subsystem().ID()], connection)

	return nil
}

func (ws *WorldState) Init() {

	ws.subsystems = make(map[subsystems.SubsystemID]subsystems.Subsystem)
	ws.connections = make(map[subsystems.SubsystemID][]*connections.Connection)
	ws.ports = make(map[connections.PortID]*connections.Port)

	reactor := subsystems.NewPower(.5)
	cooler := subsystems.NewCooling(.5)
	acUnit := subsystems.NewHvac()

	ws.addSubsystem(reactor)
	ws.addSubsystem(cooler)
	ws.addSubsystem(acUnit)

	acPowerPort := ws.AddPort("socket-1", acUnit, "power-in")
	acTempPort := ws.AddPort("valve-1", acUnit, "temp-in")
	reactorPowerPort := ws.AddPort("socket-1", reactor, "power-out")

	reactorTempInPort := ws.AddPort("valve-1", reactor, "temp-in")
	reactorTempOutPort := ws.AddPort("valve-1", reactor, "temp-out")
	coolerTempPort := ws.AddPort("valve-1", cooler, "temp-out")

	err := ws.addConnection(reactorPowerPort, acPowerPort, 0.5)
	if err != nil {
		logging.Error(err.Error())
	}

	err = ws.addConnection(reactorTempOutPort, acTempPort, 0.5)
	if err != nil {
		logging.Error(err.Error())
	}

	err = ws.addConnection(coolerTempPort, reactorTempInPort, 1)
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
				subsystem.Tick(nil)
			}

			for _, conn := range ws.connections[subsystem.ID()] {
				port := ws.ports[conn.Src().ID()]
				src := port.Subsystem()
				if _, seen := visited[src.ID()]; !seen {
					subsystem := ws.subsystems[src.ID()]
					depStack.Push(subsystem)

				}
			}
		}

		inputs := make(map[string]components.Component, 0)
		for _, conn := range ws.connections[system.ID()] {
			srcComp := *conn.Src().Component()
			destComp := *conn.Dest().Component()
			srcComp.SetValue(srcComp.Value() * conn.Throughput())

			inputs[destComp.Name()] = srcComp
		}
		system.Tick(inputs)

	}

}
