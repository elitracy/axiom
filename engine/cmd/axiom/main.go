package main

import (
	"github.com/elias/axiom/engine"
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/materials"
	"github.com/elias/axiom/engine/simulation"
	"github.com/elias/axiom/engine/systems/cooling"
	"github.com/elias/axiom/engine/systems/machines"
	"github.com/elias/axiom/engine/systems/power"
)

func main() {
	startTick := engine.NewTick()
	logging.Init("logging/logs/debug.log", startTick)

	generator := power.NewGenerator(materials.NewSteel(), 25.0)
	coolant := cooling.NewCoolantLoop(materials.NewWater(), materials.NewSteel())
	hvac := machines.NewHvac(25.0)
	scrubber := machines.NewScrubber()
	lifeSupport := machines.NewLifeSupport()

	world := &simulation.WorldState{
		Power:       generator,
		Coolant:     coolant,
		Hvac:        hvac,
		Scrubber:    scrubber,
		LifeSupport: lifeSupport,
	}

	logging.Ok("STARTING AXIOM")
	engine.RunGame(world, startTick)

}
