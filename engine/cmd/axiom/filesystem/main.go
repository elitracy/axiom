package main

import (
	"log"

	"github.com/elias/axiom/engine/filesystem"
)

func main() {

	home := filesystem.NewNode("root/")
	systems := filesystem.NewNode("systems/")
	sensors := filesystem.NewNode("sensors/")

	home.AddChild(systems)
	home.AddChild(sensors)

	power := filesystem.NewNode("power/")
	coolant := filesystem.NewNode("coolant/")
	lifesupport := filesystem.NewNode("life_support/")

	systems.AddChild(power)
	systems.AddChild(coolant)
	systems.AddChild(lifesupport)

	ls := home.Ls("")
	log.Printf("\n" + ls)
}
