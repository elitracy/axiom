package main

import (
	"github.com/elias/axiom/engine"
	"github.com/elias/axiom/engine/filesystem"
	"github.com/elias/axiom/engine/logging"
)

func main() {

	startTick := engine.NewTick()
	logging.Init("logging/logs/debug.log", startTick)

	shell := filesystem.NewShell()
	shell.Populate(nil)

	logging.Flush()
}
