package main

import (
	"github.com/elias/axiom/engine"
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/simulation"
	"github.com/elias/axiom/engine/systems/power"
)

func main() {
	world := &simulation.WorldState{}
	world.AddSystem(systems.NewGenerator(25.0))
	logging.NewLogger("logging/logs/debug.log")

	tick := engine.NewTick(0)
	logging.SetTick(tick)
	engine.RunGame(world, tick)

}
