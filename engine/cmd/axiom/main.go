package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
		return fmt.Sprintf("error: %s", err.Error())
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

	file, err := os.Create("telemetry.csv")
	if err != nil {
		logging.Error("Could not open file: telemetry.csv")
		game.log.Print(err.Error())
		logging.Flush()
		return
	}

	game.world.AddWriter("telemetry.csv", file, startTick)

	logging.Init("logging/logs/debug.log", startTick)

	configFile, err := os.ReadFile(initialConfig)
	if err != nil {
		logging.Error("Could not read file: %s", initialConfig)
		game.log.Print(err.Error())
		logging.Flush()
		return
	}

	game.cmd("write", "/usr/conf/station.ax", string(configFile))

	logging.Ok("===STARTING AXIOM===")

	game.cmd("reload")

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("> ")
			if !scanner.Scan() {
				break
			}

			parts := strings.Fields(scanner.Text())
			if len(parts) == 0 {
				fmt.Println()
				continue
			}

			if parts[0] == "exit" {
				fmt.Println("Shutting down...")
				os.Exit(0)
			}

			cmd := game.cmd(parts[0], parts[1:]...)

			fmt.Println(cmd)
			fmt.Println()
			logging.Info(cmd)
		}
	}()

	engine.RunGame(game.world, startTick)

}
