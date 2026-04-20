package main

import (
	"os"
	"time"

	"github.com/elias/axiom/engine"
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/parser"
	"github.com/elias/axiom/engine/state"
)

func main() {
	startTick := engine.NewTick()
	logging.Init("logging/logs/debug.log", startTick)
	args := os.Args[1:]
	path := args[0]

	stationConfig := parser.NewParserConfig()
	parser := parser.NewParser(stationConfig)
	file, err := os.ReadFile(path)
	if err != nil {
		logging.Error("Could not read file: %s", path)
		logging.Flush()
		return
	}

	parser.Parse(file)

	world := &state.WorldState{}
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

	go func() {
		for {
			system, err := world.GetSubsystem("ac")
			if err != nil {
				logging.Error(err.Error())
			} else {
				logging.Info(system.String())
			}
			time.Sleep(2 * time.Second)
		}
	}()

	engine.RunGame(world, startTick)
}
