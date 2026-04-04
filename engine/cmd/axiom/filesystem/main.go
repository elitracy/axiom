package main

import (
	"log"

	"github.com/elias/axiom/engine/filesystem"
)

func main() {

	root := filesystem.NewNode("root/")
	shell := filesystem.NewShell(root)

	systems := filesystem.NewNode("systems/")
	sensors := filesystem.NewNode("sensors/")
	log.Printf(systems.Name())
	log.Printf(sensors.Name())

	root.AddChild(systems)
	root.AddChild(sensors)

	power := filesystem.NewNode("power/")
	coolant := filesystem.NewNode("coolant/")
	lifesupport := filesystem.NewNode("life_support/")

	systems.AddChild(power)
	systems.AddChild(coolant)
	systems.AddChild(lifesupport)

	log.Printf("\n" + shell.Ls(""))

}
