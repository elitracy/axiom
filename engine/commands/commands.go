package commands

import (
	"github.com/elias/axiom/engine/filesystem"
	"github.com/elias/axiom/engine/parser"
)

type config any

type state interface {
	ValidateConfig(parser.ParserConfig) []error
	ApplyConfig(parser.ParserConfig) error
}

// NOTE: only allowing global config ATM
func Reload(state state, config parser.ParserConfig) []error {
	errs := state.ValidateConfig(config)
	if errs != nil {
		return errs
	}

	err := state.ApplyConfig(config)
	if err != nil {
		return []error{err}
	}

	return nil
}

func Cat(shell filesystem.Shell, path string) string {
	return shell.Cat(path)
}

func Ls(shell filesystem.Shell, path string) string {
	return shell.Ls(path)
}

func Tree(shell filesystem.Shell, path string, depth int) string {
	return shell.Tree(path, depth)
}

func Write(shell *filesystem.Shell, path string, content string) {
	node := shell.Find(path)

	node.Write(content)
}
