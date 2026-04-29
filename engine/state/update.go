package state

import (
	"github.com/elias/axiom/engine"
	"github.com/elias/axiom/engine/subsystems"
	"github.com/elias/axiom/engine/utils"
)

func (ws *State) Update(tick *engine.Tick) {
	ws.tickSubsystems()
}

func (ws *State) tickSubsystems() {
	visited := make(map[subsystems.SubsystemID]struct{})

	depStack := utils.NewStack[Subsystem]()

	for _, system := range ws.subsystems {
		depStack.Push(system)

		for depStack.Len() > 0 {
			subsystem := depStack.Pop()
			if _, seen := visited[subsystem.ID()]; seen {
				continue
			}

			visited[subsystem.ID()] = struct{}{}

			for _, conn := range ws.connections[subsystem.Name()] {
				src := ws.subsystems[conn.SrcSystem()]
				if _, seen := visited[src.ID()]; !seen {
					subsystem := ws.subsystems[src.Name()]
					depStack.Push(subsystem)

				}
			}
		}

		for _, conn := range ws.connections[system.Name()] {
			srcComp := *conn.Src().Component()
			conn.Dest().SetValue(srcComp.Value() * conn.Throughput())
		}
		system.Tick()

	}

}
