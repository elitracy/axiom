package simulation

import (
	"fmt"
	"strconv"

	"github.com/elias/axiom/engine"
	"github.com/elias/axiom/engine/config"
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/subsystems"
	"github.com/elias/axiom/engine/subsystems/connections"
	"github.com/elias/axiom/engine/utils"
)

type WorldState struct {
	subsystems  map[string]subsystems.Subsystem
	connections map[subsystems.SubsystemID][]*connections.Connection
}

func (ws *WorldState) ValidateConfig(stationConfig config.StationConfig) []error {

	tempSubsystems := make(map[string]subsystems.Subsystem)

	var errors []error

	if len(stationConfig.Errors) > 0 {
		for _, err := range stationConfig.Errors {
			errors = append(errors, fmt.Errorf("%s", err))
		}
	}

	for name, systemType := range stationConfig.SubsystemDeclarations {
		system, err := config.NewSubsystem(name, systemType)
		if err != nil {
			errors = append(errors, err)
		}
		tempSubsystems[name] = system
	}

	for _, setDir := range stationConfig.SetDirectives {
		system, exists := tempSubsystems[setDir.System]

		if !exists {
			err := fmt.Errorf("subsystem %s does not exist", setDir.System)
			errors = append(errors, err)
		}

		_, exists = system.Components()[setDir.Component]

		if !exists {
			err := fmt.Errorf("Component %s does not exist on system %s", setDir.Component, setDir.System)
			errors = append(errors, err)
		}

		parsedFloat, err := strconv.ParseFloat(setDir.Value, 64)
		if err != nil || parsedFloat < 0 || parsedFloat > 1 {
			err := fmt.Errorf("Not a valid component value %s, must be [0-1]", setDir.Value)
			errors = append(errors, err)
		}
	}

	for _, connection := range stationConfig.ConnectionDeclarations {
		srcSystem, exists := tempSubsystems[connection.SrcSystem]
		if !exists {
			err := fmt.Errorf("Source subsystem %s does not exist", connection.SrcSystem)
			errors = append(errors, err)
		}

		_, exists = srcSystem.OutputPorts()[connection.SrcPort]
		if !exists {
			err := fmt.Errorf("Port %s does not exist on subsystem %s", connection.SrcPort, connection.SrcSystem)
			errors = append(errors, err)
		}

		destSystem, exists := tempSubsystems[connection.DestSystem]
		if !exists {
			err := fmt.Errorf("Destination subsystem %s does not exist", connection.DestSystem)
			errors = append(errors, err)
		}

		_, exists = destSystem.InputPorts()[connection.DestPort]
		if !exists {
			err := fmt.Errorf("Port %s does not exist on subsystem %s", connection.DestPort, connection.DestSystem)
			errors = append(errors, err)
		}

		throughputFloat, err := strconv.ParseFloat(connection.Throughput, 64)
		if err != nil || throughputFloat < 0 || throughputFloat > 1 {
			err := fmt.Errorf("Not a valid throughput value %s, must be [0-1]", connection.Throughput)
			errors = append(errors, err)
		}

	}

	return errors
}

func (ws *WorldState) ApplyConfig(stationConfig config.StationConfig) {

	for name, systemType := range stationConfig.SubsystemDeclarations {
		if _, exists := ws.subsystems[name]; !exists {
			system, _ := config.NewSubsystem(name, systemType)
			ws.addSubsystem(system)
		}
	}

	for _, setDir := range stationConfig.SetDirectives {
		system := ws.subsystems[setDir.System]
		comp := system.Components()[setDir.Component]
		parsedFloat, _ := strconv.ParseFloat(setDir.Value, 64)

		comp.SetValue(utils.Unit(parsedFloat))
	}

	for _, connection := range stationConfig.ConnectionDeclarations {
		srcSystem := ws.subsystems[connection.SrcSystem]
		srcPort := srcSystem.OutputPorts()[connection.SrcPort]

		destSystem := ws.subsystems[connection.DestSystem]
		destPort := destSystem.InputPorts()[connection.DestPort]

		throughputFloat, _ := strconv.ParseFloat(connection.Throughput, 64)

		ws.addConnection(srcPort, destPort, utils.Unit(throughputFloat))
	}

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

func (ws *WorldState) Update(tick *engine.Tick) {
	ws.updateSubsystems()

	for name := range ws.subsystems {
		logging.Info(ws.subsystems[name].String())
		logging.Info("")
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
					subsystem := ws.subsystems[src.Name()]
					depStack.Push(subsystem)

				}
			}
		}

		for _, conn := range ws.connections[system.ID()] {
			srcComp := *conn.Src().Component()

			unit := new(utils.Unit)
			*unit = srcComp.Value() * conn.Throughput()

			conn.Dest().SetInput(unit)
		}
		system.Tick()

	}

}
