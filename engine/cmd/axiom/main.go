package main

import (
	"github.com/elias/axiom/engine"
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/materials"
	"github.com/elias/axiom/engine/simulation"
	"github.com/elias/axiom/engine/systems/cooling"
	"github.com/elias/axiom/engine/systems/power"
)

func main() {
	startTick := engine.NewTick()
	logging.Init("logging/logs/debug.log", startTick)

	generator := power.NewGenerator(materials.NewSteel(), 20.0)
	coolant := cooling.NewCoolantLoop(materials.NewPropyleneGlycol(), materials.NewSteel())

	generator.UpdateCoolingSource(coolant.Temperature())
	coolant.UpdateHeatSource(generator.Temperature())

	world := &simulation.WorldState{
		Power:   generator,
		Coolant: coolant,
	}

	engine.RunGame(world, startTick)

}
