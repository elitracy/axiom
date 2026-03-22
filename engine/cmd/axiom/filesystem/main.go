package main

import (
	"log"

	"github.com/elias/axiom/engine/filesystem"
)

func main() {

	dir := filesystem.NewDirNode("TOP")
	subDirA := filesystem.NewDirNode("subDirA")
	subDirC := filesystem.NewDirNode("subDirC")
	subDirA.AddChild("fileC", filesystem.NewFileNode("fileC"))
	subDirA.AddChild("subDirC", subDirC)

	subDirC.AddChild("fileD", filesystem.NewFileNode("fileD"))

	dir.AddChild("subDirA", subDirA)
	dir.AddChild("subDirB", filesystem.NewDirNode("subDirB"))
	dir.AddChild("fileA", filesystem.NewFileNode("fileA"))
	dir.AddChild("fileB", filesystem.NewFileNode("fileB"))

	ls := dir.Ls("subDirA/..")
	log.Printf("\n" + ls)
}
