package main

import (
	"fmt"

	"github.com/elias/axiom/engine/filesystem"
)

func main() {

	root := filesystem.NewDir("/")
	shell := filesystem.NewShell(root)

	systems := filesystem.NewDir("systems/")
	sensors := filesystem.NewDir("sensors/")
	file := filesystem.NewFile("file.txt")

	root.AddChild(systems)
	root.AddChild(sensors)

	power := filesystem.NewDir("power/")
	coolant := filesystem.NewDir("coolant/")
	lifesupport := filesystem.NewDir("life_support/")

	systems.AddChild(power)
	systems.AddChild(coolant)
	systems.AddChild(lifesupport)

	power.AddChild(file)
	shell.Cd("systems/power/")

	file.Write("test input\nnew line as well")
	// fmt.Println(shell.Cat("file.txt"))
	fmt.Printf(shell.Ls(""))

}
