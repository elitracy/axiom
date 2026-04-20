package main

import (
	"os"
	"time"

	"github.com/elias/axiom/engine"
	"github.com/elias/axiom/engine/commands"
	"github.com/elias/axiom/engine/filesystem"
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/parser"
	"github.com/elias/axiom/engine/state"
)

func main() {
	startTick := engine.NewTick()
	logging.Init("logging/logs/debug.log", startTick)

	args := os.Args[1:]
	if len(args) != 1 {
		logging.Error("No path provided")
		logging.Flush()
		return
	}

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

	errs := commands.Reload(world, parser.Config)

	if len(errs) > 0 {
		for _, err := range errs {
			logging.Error(err.Error())
		}
		logging.Flush()
		return
	}

	logging.Ok("RELOADED CONFIG")

	shell := filesystem.NewShell()
	shell.Populate(world)

	logging.Debug(shell.Tree("", 6))

	go func() {
		for {
			for _, s := range world.Subsystems() {
				logging.Debug(s.String())
			}
			time.Sleep(2 * time.Second)
		}
	}()

	engine.RunGame(world, startTick)

}
