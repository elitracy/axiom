package commands

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/elias/axiom/engine/filesystem"
	"github.com/elias/axiom/engine/parser"
	"github.com/elias/axiom/engine/state"
)

type CommandEngine struct {
	state *state.State
	shell *filesystem.Shell
}

func NewCommandEngine(state *state.State, shell *filesystem.Shell) *CommandEngine {
	ce := &CommandEngine{
		state: state,
		shell: shell,
	}

	return ce
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
	case "inspect":
		return ce.inspect(args)
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
	parser := parser.NewParser(parser.NewParserConfig())

	err := parser.Parse([]byte(node.Read()))
	if err != nil {
		return "", []error{err}
	}

	errs := ce.state.ValidateConfig(parser.Config)
	if errs != nil {
		return "", errs
	}

	return "config valid", nil
}

func (ce *CommandEngine) reload(args []string) (string, error) {
	var node *filesystem.Node
	config := "station.ax"

	if len(args) == 1 {
		config = args[0]
	}

	node = ce.shell.Find(config)
	if node == nil {
		return "", fmt.Errorf("config %s not found", config)
	}

	parser := parser.NewParser(parser.NewParserConfig())

	err := parser.Parse([]byte(node.Read()))
	if err != nil {
		return "", err
	}

	err = ce.state.ApplyConfig(parser.Config)
	if err != nil {
		return "", err
	}

	ce.shell.ReloadSubsystems(ce.state)

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
	path := "."
	depth := 2

	switch len(args) {
	case 0:
		return "", fmt.Errorf("usage: tree <path> <depth?>")
	case 1:
		path = args[0]
	case 2:
		path = args[0]
		d, err := strconv.Atoi(args[1])
		if err != nil {
			return "", fmt.Errorf("usage: tree <path> <depth?>")
		}
		depth = d
	default:
		return "", fmt.Errorf("usage: tree <path> <depth?>")
	}

	return ce.shell.Tree(path, depth), nil
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

func (ce *CommandEngine) inspect(args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("usage: inspect <subsystem>")
	}

	subsystem := args[0]

	node := ce.shell.Find(subsystem)
	if node == nil {
		return "", fmt.Errorf("subsystem: %s does not exist", subsystem)
	}

	statusFile := node.FindChild("status")
	if statusFile == nil {
		return "", fmt.Errorf("%s status doesn't exist", subsystem)
	}

	var sortedComponents []*filesystem.Node
	componentsDir := node.FindChild("components")
	for _, child := range componentsDir.Children() {
		sortedComponents = append(sortedComponents, child)
	}

	slices.SortFunc(sortedComponents, func(a, b *filesystem.Node) int {
		return strings.Compare(a.Name(), b.Name())
	})

	status := statusFile.Read()
	var output string

	output = fmt.Sprintf("%s [%s]", subsystem, status)
	for _, component := range sortedComponents {
		compType := component.FindChild("type")
		compValue := component.FindChild("value")
		compPorts := component.FindChild("ports")

		output += fmt.Sprintf("\n==%s %s %s==", component.Name(), compType.Read(), compValue.Read())

		var sortedPorts []*filesystem.Node
		for _, port := range compPorts.Children() {
			sortedPorts = append(sortedPorts, port)
		}

		slices.SortFunc(sortedPorts, func(a, b *filesystem.Node) int {
			return strings.Compare(a.Name(), b.Name())
		})

		for _, port := range sortedPorts {
			conn := port.FindChild("connection")
			if conn != nil {
				output += fmt.Sprintf("\n%s %s", port.Name(), conn.Read())
			}
		}
	}

	return output, nil
}

func (ce *CommandEngine) status(args []string) (string, error) {
	if len(args) != 0 {
		return "", fmt.Errorf("usage: status")
	}

	systemsDir := ce.shell.Find("systems")
	if systemsDir == nil {
		return "", fmt.Errorf("subsystems directory does not exist")
	}

	var output string
	for _, subsystemTypeDir := range systemsDir.Children() {
		output += subsystemTypeDir.Name()
		for _, subsystemDir := range subsystemTypeDir.Children() {
			statusFile := subsystemTypeDir.FindChild("status")
			if statusFile == nil {
				return "", fmt.Errorf("%s status doesn't exist", subsystemDir)
			}

			var sortedComponents []*filesystem.Node
			componentsDir := subsystemDir.FindChild("components")
			for _, child := range componentsDir.Children() {
				sortedComponents = append(sortedComponents, child)
			}

			slices.SortFunc(sortedComponents, func(a, b *filesystem.Node) int {
				return strings.Compare(a.Name(), b.Name())
			})

			output += fmt.Sprintf("\n %s [%s]", subsystemDir.Name(), statusFile.Read())
			for _, component := range sortedComponents {
				compType := component.FindChild("type")
				compValue := component.FindChild("value")
				output += fmt.Sprintf("\n  %s (%s): %s", component.Name(), compType.Read(), compValue.Read())
			}
			output += "\n"
		}
	}

	return output, nil
}
