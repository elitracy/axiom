package main

import (
	"github.com/elias/axiom/engine/filesystem"
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/utils"
)

func main() {

	startTick := utils.NewTick()
	logging.Init("logging/logs/debug.log", startTick)

	shell := filesystem.NewShell()
	shell.Populate(nil)

	logging.Flush()
}
