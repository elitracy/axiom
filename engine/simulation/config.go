package simulation

import (
	"fmt"
	"strconv"

	"github.com/elias/axiom/engine/parser"
	"github.com/elias/axiom/engine/subsystems"
	"github.com/elias/axiom/engine/utils"
)

func (ws *WorldState) ValidateConfig(stationConfig parser.ParserConfig) []error {

	tempSubsystems := make(map[string]subsystems.Subsystem)

	var errors []error

	if len(stationConfig.Errors) > 0 {
		for _, err := range stationConfig.Errors {
			errors = append(errors, fmt.Errorf("%s", err))
		}
	}

	for name, systemType := range stationConfig.SubsystemDeclarations {
		system, err := parser.NewSubsystem(name, systemType)
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

func (ws *WorldState) ApplyConfig(stationConfig parser.ParserConfig) {

	for name, systemType := range stationConfig.SubsystemDeclarations {
		if _, exists := ws.subsystems[name]; !exists {
			system, _ := parser.NewSubsystem(name, systemType)
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
