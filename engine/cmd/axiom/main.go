package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/elias/axiom/engine/game"
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/utils"
)

const (
	minMsgSleep = 1
	maxMsgSleep = 3
)

func getStartupMessages() []string {
	startupMessages := []string{
		"BOOTING AXIOM",
		"SCANNING SYSTEM",
		"IDENTIFYING SUBSYSTEMS",
		"RELOADING COMPONENTS",
	}

	return startupMessages
}

func showStartupMessages() {
	for _, msg := range getStartupMessages() {
		fmt.Printf("%s\n", msg)

		sleep := rand.Intn(maxMsgSleep-minMsgSleep+1) + minMsgSleep
		time.Sleep(time.Duration(sleep) * time.Second)
	}

	time.Sleep(time.Duration(1) * time.Second)
	fmt.Printf("\nSYSTEM STATUS: ")
	time.Sleep(time.Duration(2) * time.Second)
	fmt.Printf("WARNING\n\n")

	time.Sleep(time.Duration(1) * time.Second)
	fmt.Print("> help\n")
	time.Sleep(time.Duration(1) * time.Second)

}

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
		return
	}

	// showStartupMessages()

	cmd := g.Cmd("help")

	fmt.Printf("%s\n\n", cmd)
	logging.Info(cmd)
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

		logging.Info("%v", parts)
	}

}
