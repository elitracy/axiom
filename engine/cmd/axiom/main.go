package main

import (
	"github.com/elias/axiom/engine"
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/simulation"
	"github.com/elias/axiom/engine/systems/power"
)

func main() {
	startTick := engine.NewTick()
	logging.Init("logging/logs/debug.log", startTick)

	world := &simulation.WorldState{}
	world.AddSystem(systems.NewGenerator(25.0))

	engine.RunGame(world, startTick)

}
