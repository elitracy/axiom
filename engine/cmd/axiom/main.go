package main

import (
	"os"

	"github.com/elias/axiom/engine"
	"github.com/elias/axiom/engine/commands"
	"github.com/elias/axiom/engine/filesystem"
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/parser"
	"github.com/elias/axiom/engine/state"
)

func main() {

	args := os.Args[1:]
	if len(args) != 1 {
		logging.Error("No path provided")
		logging.Flush()
		return
	}

	initialConfig := args[0]

	startTick := engine.NewTick()
	stationConfig := parser.NewParserConfig()
	parser := parser.NewParser(stationConfig)
	world := state.NewWorldState()
	shell := filesystem.NewShell()
	commandEngine := commands.NewCommandEngine(world, shell, stationConfig)

	logging.Init("logging/logs/debug.log", startTick)
	gamelog := state.NewGameLogger(512)
	shell.Populate(world)

	file, err := os.ReadFile(initialConfig)
	if err != nil {
		logging.Error("Could not read file: %s", initialConfig)
		gamelog.Write(err.Error())
		logging.Flush()
		return
	}

	parser.Parse(file)

	logging.Ok("STARTING AXIOM")

	_, err = commandEngine.Execute("reload")

	if err != nil {
		logging.Error(err.Error())
		gamelog.Write(err.Error())
		logging.Flush()
		return
	}

	logging.Ok("RELOADED CONFIG")

	_, err = commandEngine.Execute("write", "/usr/conf/station.ax", string(file))
	if err != nil {
		logging.Error("Could not read file: %s", initialConfig)
		gamelog.Write(err.Error())
		logging.Flush()
		return
	}

	tree, err := commandEngine.Execute("tree", ".", "6")
	if err != nil {
		logging.Error("Could not tree: .")
		gamelog.Write(err.Error())
		logging.Flush()
		return
	}
	logging.Debug("TREE: %s", tree)

	newConf := string(file) + "\nsystem fooReactor type=power"
	newConf += "\nset fooReactor.power-out 0.2"
	newConf += "\nconnect coolant_loop.out.valve-2 -> fooReactor.in.valve-2 0.5"
	newConf += "\nconnect fooReactor.out.socket-1 -> ac.in.socket-2 0.5"
	newConf += "\nconnect fooReactor.out.valve-1 -> ac.in.valve-2 0.5"

	_, err = commandEngine.Execute("write", "/usr/conf/station.ax", newConf)

	conf := shell.Cat("/usr/conf/station.ax")
	parser.Parse([]byte(conf))

	_, err = commandEngine.Execute("reload")
	if err != nil {
		logging.Error(err.Error())
		gamelog.Write(err.Error())
		logging.Flush()
		return
	}

	logging.Ok("RELOADED CONFIG")

	// logging.Debug(shell.Tree("", 6))

	// go func() {
	// 	for {
	// 		for _, s := range world.Subsystems() {
	// 			logging.Debug(s.String())
	// 		}
	// 		time.Sleep(2 * time.Second)
	// 	}
	// }()

	status, err := commandEngine.Execute("status", "coolant_loop")
	if err != nil {
		logging.Error(err.Error())
		gamelog.Write(err.Error())
		logging.Flush()
		return
	}
	logging.Debug(status)

	engine.RunGame(world, startTick)

}
