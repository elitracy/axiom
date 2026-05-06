package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"strings"

	"github.com/elias/axiom/engine/game"
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/utils"
)

var globalGame *game.Game

//export AxiomInit
func AxiomInit(configPath *C.char) *C.char {
	tick := utils.NewTick()
	logging.Init("logging/logs/debug.log", tick)

	globalGame = game.NewGame(tick)
	err := globalGame.Init(C.GoString(configPath))

	if err != nil {
		return C.CString(err.Error())
	}

	globalGame.StartTickLoop()
	return C.CString("")
}

//export AxiomExecute
func AxiomExecute(input *C.char) *C.char {

	inputString := C.GoString(input)
	parts := strings.Fields(inputString)

	var cmd string
	switch len(parts) {
	case 0:
		return C.CString("")
	case 1:
		cmd = globalGame.Cmd(parts[0])
	default:
		cmd = globalGame.Cmd(parts[0], parts[1:]...)
	}

	return C.CString(cmd)
}
