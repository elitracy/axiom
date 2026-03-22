package simulation

import (
	"github.com/elias/axiom/engine"
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/systems/cooling"
	"github.com/elias/axiom/engine/systems/machines"
	"github.com/elias/axiom/engine/systems/power"
)

type WorldState struct {
	Power       power.Power
	Coolant     cooling.Coolant
	LifeSupport *machines.LifeSupport
	Hvac        *machines.Hvac
	Scrubber    *machines.Scrubber
}

var lastCoolantOut cooling.CoolantOutput
var lastHvacOut machines.HvacOutput

func (ws *WorldState) Update(tick *engine.Tick) {

	var heatSources []float64

	powerOut := ws.Power.Tick(power.PowerInput{CoolantTemperature: lastCoolantOut.Temperature})
	coolantOut := ws.Coolant.Tick(cooling.CoolantInput{LoadTemperature: powerOut.Temperature})

	heatSources = append(heatSources, powerOut.Temperature)

	hvacOut := ws.Hvac.Tick(machines.HvacInput{PowerSupplied: powerOut.Power * .5, HeatSources: heatSources})
	scrubberOut := ws.Scrubber.Tick(machines.ScrubberInput{PowerAvailable: powerOut.Power * .25})
	ws.LifeSupport.Tick(machines.LifeSupportInput{PowerAvailable: powerOut.Power * .25, TemperatureStatus: hvacOut.Status, OxygenStatus: scrubberOut.Status})

	lastCoolantOut = coolantOut
	lastHvacOut = hvacOut

	logging.Info("%s", ws.Power.String())
	logging.Info("%s", ws.Coolant.String())
	// logging.Info("%s", ws.Hvac.String())
	// logging.Info("%s", ws.Scrubber.String())
	// logging.Info("%s", ws.LifeSupport.String())
}
