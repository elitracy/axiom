package simulation

import (
	"github.com/elias/axiom/engine"
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/subsystems"
	"github.com/elias/axiom/engine/utils"
)

func (ws *WorldState) Update(tick *engine.Tick) {
	ws.updateSubsystems()

	for name := range ws.subsystems {
		logging.Info(ws.subsystems[name].String())
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
