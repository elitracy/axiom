package simulation

import (
	"github.com/elias/axiom/engine"
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/systems"
)

type WorldState struct {
	systems []systems.System
}

func (ws *WorldState) AddSystem(system systems.System) {
	ws.systems = append(ws.systems, system)
}

func (ws *WorldState) Update(tick engine.Tick) {
	for _, system := range ws.systems {
		system.Tick()
		logging.Info(system.String())
	}
}
