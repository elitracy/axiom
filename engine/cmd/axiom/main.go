package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/elias/axiom/engine/game"
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/utils"
)

func main() {
	tick := utils.NewTick()
	logging.Init("logging/logs/debug.log", tick)

	args := os.Args
	if len(args) != 2 {
		logging.Error("No path provided")
		logging.Flush()
		return
	}

	initialConfig := args[1]

	g := game.NewGame(tick)

	err := g.Init(initialConfig)
	if err != nil {
		logging.Error(err.Error())
		g.Log().Print(err.Error())
		logging.Flush()
	}

	g.StartTickLoop()

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

		cmd := g.Cmd(parts[0], parts[1:]...)

		fmt.Printf("%s\n\n", cmd)
		logging.Info(cmd)
	}

}
