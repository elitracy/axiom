package state

import (
	"github.com/elias/axiom/engine"
	"github.com/elias/axiom/engine/subsystems"
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

func (s *State) Update(tick *engine.Tick) {
	s.tickSubsystems()
	s.writeSubsystems()
}

func (s *State) writeSubsystems() {
	for _, w := range s.writers {
		for _, system := range s.Subsystems() {
			w.Write(system.ExportFields())
		}
	}
}

func (s *State) tickSubsystems() {

	for _, system := range s.subsystems {
		for _, port := range system.InputPorts() {
			port.Component().Clear()
		}
	}

	sums := make(map[*components.Component]utils.Unit)

	for _, conns := range s.connections {
		for _, conn := range conns[utils.PortInput] {
			dest := conn.DestPort().Component()
			sums[dest] += conn.SrcPort().Component().Value() * conn.Throughput()
		}
	}

	for conn, value := range sums {
		conn.SetValue(value)
	}

	for _, system := range s.topoSort() {
		system.Tick()
	}
}

func (s *State) topoSort() []Subsystem {
	visited := make(map[subsystems.SubsystemID]struct{})
	var sorted []Subsystem

	var visit func(subsystem Subsystem)
	visit = func(subsystem Subsystem) {

		if _, seen := visited[subsystem.ID()]; seen {
			return
		}

		visited[subsystem.ID()] = struct{}{}

		for _, conn := range s.connections[subsystem.Name()][utils.PortInput] {
			src := s.subsystems[conn.SrcSystem()]
			visit(src)
		}

		sorted = append(sorted, subsystem)
	}

	for _, system := range s.subsystems {
		visit(system)
	}

	return sorted
}
