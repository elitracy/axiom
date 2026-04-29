package commands

import (
	"fmt"
	"strconv"
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

func (ce *CommandEngine) Execute(cmd string, args ...string) (string, error) {
	cmd = strings.Trim(cmd, " ")
	cmd = strings.ToLower(cmd)

	switch cmd {
	case "diagnose":
		val, errs := ce.diagnose(args)
		var err error

		if errs != nil {
			var errMsgs []string
			for _, err := range errs {
				errMsgs = append(errMsgs, err.Error())
			}

			err = fmt.Errorf("%s", strings.Join(errMsgs, "\n"))
		}

		return val, err
	case "reload":
		return ce.reload(args)
	case "cat":
		return ce.cat(args)
	case "ls":
		return ce.ls(args)
	case "tree":
		return ce.tree(args)
	case "write":
		return "", ce.write(args)
	case "status":
		return ce.status(args)

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

func (ce *CommandEngine) reload(args []string) (string, error) {
	if len(args) != 0 {
		return "", fmt.Errorf("usage: reload")
	}

	node := ce.shell.Find("station.ax")
	if node == nil {
		return "", fmt.Errorf("config station.ax not found")
	}

	err := ce.parser.Parse([]byte(node.Read()))
	if err != nil {
		return "", err
	}

	err = ce.state.ApplyConfig(ce.parser.Config)
	if err != nil {
		return "", err
	}

	ce.shell.Populate(ce.state)

	return "station reloaded", nil
}

func (ce *CommandEngine) cat(args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("usage: cat <path>")
	}

	path := args[0]
	return ce.shell.Cat(path), nil
}

func (ce *CommandEngine) ls(args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("usage: ls <path>")
	}

	return ce.shell.Ls(args[0]), nil
}

func (ce *CommandEngine) tree(args []string) (string, error) {
	switch len(args) {
	case 0:
		return "", fmt.Errorf("usage: tree <path> <depth?>")
	case 1:
		return ce.shell.Tree(args[0], 2), nil
	case 2:
		depth, err := strconv.Atoi(args[1])
		if err != nil {
			return "", fmt.Errorf("usage: tree <path> <depth?>")
		}
		return ce.shell.Tree(args[0], depth), nil
	default:
		return "", fmt.Errorf("usage: tree <path> <depth?>")
	}
}

func (ce *CommandEngine) write(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: write <path> <content>")
	}
	path := args[0]
	content := args[1]

	node := ce.shell.GetChild(path)

	if node == nil {
		return fmt.Errorf("%s does not exist", path)
	}

	node.Write(content)

	return nil
}

func (ce *CommandEngine) status(args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("usage: status <subsystem>")
	}

	subsystem := args[0]

	node := ce.shell.Find(subsystem)
	if node == nil {
		return "", fmt.Errorf("subsystem: %s does not exist", subsystem)
	}

	return "Healthy", nil
}
