package main

import (
	"github.com/elias/axiom/engine"
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/simulation"
)

func main() {
	startTick := engine.NewTick()
	logging.Init("logging/logs/debug.log", startTick)

	world := &simulation.WorldState{}

	logging.Ok("STARTING AXIOM")
	engine.RunGame(world, startTick)

}
