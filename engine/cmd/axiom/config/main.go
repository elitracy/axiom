package main

import (
	"os"

	"github.com/elias/axiom/engine"
	"github.com/elias/axiom/engine/config"
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/simulation"
)

func main() {
	startTick := engine.NewTick()
	logging.Init("logging/logs/debug.log", startTick)
	args := os.Args[1:]
	path := args[0]

	stationConfig := config.NewStationConfig()
	parser := config.NewParser(stationConfig)
	parser.ReadFile(path)
	logging.Info("LEN: %v", len(parser.Config.SubsystemDeclarations))

	world := &simulation.WorldState{}
	world.Init()

	logging.Ok("STARTING AXIOM")

	errs := world.ValidateConfig(parser.Config)
	if len(errs) > 0 {
		for _, err := range errs {
			logging.Error(err.Error())
		}
		logging.Flush()
		return
	}
	logging.Ok("VALID CONFIG")

	world.ApplyConfig(parser.Config)

	logging.Ok("APPLIED CONFIG")
	engine.RunGame(world, startTick)
}
