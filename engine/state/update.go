package state

import (
	"github.com/elias/axiom/engine"
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/subsystems"
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

func (ws *State) Update(tick *engine.Tick) {
	ws.tickSubsystems()
}

func (ws *State) tickSubsystems() {

	for _, system := range ws.subsystems {
		for _, port := range system.InputPorts() {
			port.Component().Clear()
		}
	}

	sums := make(map[*components.Component]utils.Unit)

	for _, conns := range ws.connections {
		for _, conn := range conns[utils.PortInput] {
			dest := conn.DestPort().Component()
			sums[dest] += conn.SrcPort().Component().Value() * conn.Throughput()
		}
	}

	for conn, value := range sums {
		conn.SetValue(value)
	}

	for _, system := range ws.topoSort() {
		logging.Debug("system: %v", system.Name())
		system.Tick()
	}
	logging.Debug("")
}

func (ws *State) topoSort() []Subsystem {
	visited := make(map[subsystems.SubsystemID]struct{})
	var sorted []Subsystem

	var visit func(subsystem Subsystem)
	visit = func(subsystem Subsystem) {

		if _, seen := visited[subsystem.ID()]; seen {
			return
		}

		visited[subsystem.ID()] = struct{}{}

		for _, conn := range ws.connections[subsystem.Name()][utils.PortInput] {
			src := ws.subsystems[conn.SrcSystem()]
			visit(src)
		}

		sorted = append(sorted, subsystem)
	}

	for _, system := range ws.subsystems {
		visit(system)
	}

	return sorted
}
