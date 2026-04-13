package config

import (
	"fmt"

	"github.com/elias/axiom/engine/subsystems"
)

type setDirective struct {
	System    string
	Component string
	Value     string
}

type connectionDeclaration struct {
	SrcSystem  string
	SrcPort    string
	DestSystem string
	DestPort   string
	Throughput string
}

type StationConfig struct {
	SubsystemDeclarations  map[string]string
	SetDirectives          []setDirective
	ConnectionDeclarations []connectionDeclaration
	Errors                 []parseError
}

func NewStationConfig() StationConfig {
	return StationConfig{
		SubsystemDeclarations: make(map[string]string),
	}
}

func NewSubsystem(name, subsystemType string) (subsystems.Subsystem, error) {
	switch subsystemType {
	case "power":
		return subsystems.NewPower(name, 0.5), nil
	case "cooling":
		return subsystems.NewCooling(name, 0.5), nil
	case "hvac":
		return subsystems.NewHvac(name, 0.2), nil
	default:
		return nil, fmt.Errorf("unknown subsystem type: %s", subsystemType)
	}
}
