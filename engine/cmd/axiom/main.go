package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/elias/axiom/engine"
	"github.com/elias/axiom/engine/commands"
	"github.com/elias/axiom/engine/filesystem"
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/state"
)

type Game struct {
	world  *state.State
	shell  *filesystem.Shell
	engine *commands.CommandEngine
	log    *state.GameLogger
}

func NewGame() *Game {
	world := state.NewState()
	shell := filesystem.NewShell()

	shell.Populate(world)
	engine := commands.NewCommandEngine(world, shell)
	log := state.NewGameLogger(512)

	logNode := shell.GetChild("sys/logs/station.log")

	logNode.SetReader(func() string {
		return strings.Join(log.Read(), "\n")
	})

	return &Game{world, shell, engine, log}

}

func (g *Game) cmd(cmd string, args ...string) string {
	val, err := g.engine.Execute(cmd, args...)
	if err != nil {
		logging.Error(err.Error())
		g.log.Print(err.Error())
		logging.Flush()
		os.Exit(1)
		return ""
	}

	return val
}

func main() {

	args := os.Args[1:]
	if len(args) != 1 {
		logging.Error("No path provided")
		logging.Flush()
		return
	}

	initialConfig := args[0]

	startTick := engine.NewTick()
	game := NewGame()
	logging.Init("logging/logs/debug.log", startTick)

	file, err := os.ReadFile(initialConfig)
	if err != nil {
		logging.Error("Could not read file: %s", initialConfig)
		game.log.Print(err.Error())
		logging.Flush()
		return
	}

	game.cmd("write", "/usr/conf/station.ax", string(file))

	logging.Ok("STARTING AXIOM")

	game.cmd("reload")

	ls := game.cmd("ls", "/sys/systems/")
	logging.Info("LS: %s", ls)

	tree := game.cmd("tree", ".", "6")
	logging.Info("TREE: %s", tree)

	newConf := string(file) + "\nsystem fooReactor type=power"
	newConf += "\nset fooReactor.power-out 0.2"
	newConf += "\nconnect coolant_loop.out.valve-2 -> fooReactor.in.valve-2 0.5"
	newConf += "\nconnect fooReactor.out.socket-1 -> ac.in.socket-2 0.5"
	newConf += "\nconnect fooReactor.out.valve-1 -> ac.in.valve-2 0.5"

	game.cmd("write", "/usr/conf/station.ax", newConf)

	game.cmd("reload")

	logging.Ok("RELOADED CONFIG")

	go func() {
		for {
			for _, s := range game.world.Subsystems() {
				status := game.cmd("status", s.Name())
				msg := fmt.Sprintf("%s: %s", s.Name(), status)
				// logging.Debug(s.String())
				game.log.Println(msg)
				logging.Info(strings.Join(game.log.Read(), ""))
			}
			game.log.Println("")
			time.Sleep(2 * time.Second)
		}
	}()

	status := game.cmd("status", "coolant_loop")
	logging.Info(status)

	log := game.cmd("cat", "/sys/logs/station.log")
	logging.Info("LOG: %s", log)

	engine.RunGame(game.world, startTick)

}
