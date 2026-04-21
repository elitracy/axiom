package commands

import (
	"fmt"
	"strings"

	"github.com/elias/axiom/engine/filesystem"
	"github.com/elias/axiom/engine/parser"
	"github.com/elias/axiom/engine/state"
)

type CommandEngine struct {
	state  *state.WorldState
	shell  *filesystem.Shell
	parser *parser.Parser
}

func NewCommandEngine(state *state.WorldState, shell *filesystem.Shell, config parser.ParserConfig) *CommandEngine {
	return &CommandEngine{
		state:  state,
		shell:  shell,
		parser: parser.NewParser(config),
	}
}

func (ce *CommandEngine) Execute(cmd string, args ...string) (string, []error) {
	cmd = strings.Trim(cmd, " ")
	cmd = strings.ToLower(cmd)

	switch cmd {
	case "diagnose":
		return ce.diagnose(args)
	}

	return "", nil
}

// NOTE: only allowing global config ATM
func (ce *CommandEngine) diagnose(args []string) (string, []error) {
	if len(args) != 0 {
		return "", []error{fmt.Errorf("usage: diagnose")}
	}

	node := ce.shell.Find("station.ax")
	if node == nil {
		return "", []error{fmt.Errorf("config station.ax not found")}
	}

	err := ce.parser.Parse([]byte(node.Read()))
	if err != nil {
		return "", []error{err}
	}

	errs := ce.state.ValidateConfig(ce.parser.Config)
	if errs != nil {
		return "", errs
	}

	return "config valid", nil
}

func (ce *CommandEngine) reload(args []string) (string, []error) {
	if len(args) != 0 {
		return "", []error{fmt.Errorf("usage: reload")}
	}

	node := ce.shell.Find("station.ax")
	if node == nil {
		return "", []error{fmt.Errorf("config station.ax not found")}
	}

	err := ce.parser.Parse([]byte(node.Read()))
	if err != nil {
		return "", []error{err}
	}

	err = ce.state.ApplyConfig(ce.parser.Config)
	if err != nil {
		return "", []error{err}
	}

	ce.shell.Populate(ce.state)

	return "station reloaded", nil
}

func (ce *CommandEngine) cat(args []string) (string, []error) {
	if len(args) != 1 {
		return "", []error{fmt.Errorf("usage: cat <path>")}
	}

	path := args[0]
	return ce.shell.Cat(path), nil
}

func Ls(shell *filesystem.Shell, path string) string {
	return shell.Ls(path)
}

func Tree(shell *filesystem.Shell, path string, depth int) string {
	return shell.Tree(path, depth)
}

func Write(shell *filesystem.Shell, path string, content string) {
	node := shell.GetChild(path)
	node.Write(content)
}

func Status(shell *filesystem.Shell, subsystem string) string {
	node := shell.Find(subsystem)

	if node == nil {
		return ""
	}

	return "Healthy"

}
