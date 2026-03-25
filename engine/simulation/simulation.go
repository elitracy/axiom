package simulation

import (
	"github.com/elias/axiom/engine"
	"github.com/elias/axiom/engine/subsystems"
	"github.com/elias/axiom/engine/subsystems/components"
)

type subsystemConnection struct {
	componentID components.ComponentID
	subsystemID subsystems.SubsystemID
}

type WorldState struct {
	subsystems   map[subsystems.SubsystemID]subsystems.Subsystem
	dependencies map[subsystems.SubsystemID][]subsystemConnection
}

func (ws *WorldState) Update(tick *engine.Tick) {

	power := subsystems.NewPower(.5)
	cooling := subsystems.NewCooling(.5)
	hvac := subsystems.NewHvac()

	ws.subsystems[power.ID()] = power
	ws.subsystems[cooling.ID()] = cooling
	ws.subsystems[hvac.ID()] = hvac
}
