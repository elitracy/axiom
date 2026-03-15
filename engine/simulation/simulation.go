package simulation

import (
	"github.com/elias/axiom/engine"
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/systems/cooling"
	"github.com/elias/axiom/engine/systems/power"
)

type WorldState struct {
	Power   power.Power
	Coolant cooling.Coolant
}

func (ws *WorldState) Update(tick *engine.Tick) {
	ws.Power.Tick()
	ws.Coolant.Tick()

	logging.Info("%s", ws.Power.String())
	logging.Info("%s", ws.Coolant.String())
}
