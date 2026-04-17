package main

import (
	"os"

	"github.com/elias/axiom/engine"
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
	content, err := parser.ReadFile(path)
	if err != nil {
		logging.Error(err.Error())
		logging.Flush()
		return
	}

	parser.Parse(content)

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

	shell := filesystem.NewShell()
	shell.Populate(world)

	logging.Debug(shell.Tree("", 0))

	// logging.Ok("POWER: %v", shell.Ls("sys/systems/power"))
	// logging.Ok("POWER: %v", shell.Ls("sys/systems/power/reactor"))
	// logging.Ok("")
	// logging.Ok("COOLING: %v", shell.Ls("sys/systems/cooling"))
	// logging.Ok("COOLING: %v", shell.Ls("sys/systems/cooling/coolant_loop"))
	// logging.Ok("")
	// logging.Ok("MACHINES: %v", shell.Ls("sys/systems/machines"))
	// logging.Ok("MACHINES: %v", shell.Ls("sys/systems/machines/ac"))

	// go func() {
	// 	for {
	// 		system, err := world.GetSubsystem("ac")
	// 		if err != nil {
	// 			logging.Error(err.Error())
	// 		} else {
	// 			logging.Info(system.String())
	// 		}
	// 		time.Sleep(2 * time.Second)
	// 	}
	// }()

	engine.RunGame(world, startTick)

}
