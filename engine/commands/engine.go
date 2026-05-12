package commands

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/elias/axiom/engine/filesystem"
	"github.com/elias/axiom/engine/parser"
	"github.com/elias/axiom/engine/state"
)

type command struct {
	name     string
	helpMsg  string
	usageMsg string
	handler  func([]string) (string, error)
}

type CommandEngine struct {
	state    *state.State
	shell    *filesystem.Shell
	commands map[string]command
}

func NewCommandEngine(state *state.State, shell *filesystem.Shell) *CommandEngine {
	ce := &CommandEngine{
		state:    state,
		shell:    shell,
		commands: make(map[string]command),
	}
	ce.commands = map[string]command{
		"diagnose": {"diagnose", "checks the station config file for errors", "diagnose", ce.diagnose},
		"reload":   {"reload", "reloads the station using the station config file", "reload", ce.reload},
		"cat":      {"cat", "displays the contents of the provded file", "cat <path>", ce.cat},
		"ls":       {"ls", "displays the contents of the provided directory", "ls <path>", ce.ls},
		"tree":     {"tree", "displays the file/directory structure of the provided directory to the specified depth", "tree <path> <depth>", ce.tree},
		"write":    {"write", "writes the provided string to the specific path", "write <path> <contents>", ce.write},
		"inspect":  {"insepct", "returns the details of the specified subsystem", "inspect <subsystem>", ce.inspect},
		"status":   {"status", "returns the station's current status details", "status", ce.status},
		"help":     {"help", "returns help message", "help <cmd?>", ce.help},
		"exit":     {"exit", "closes the terminal", "exit", nil},
	}

	return ce
}

func (ce *CommandEngine) Execute(cmd string, args ...string) (string, error) {
	cmd = strings.Trim(cmd, " ")
	cmd = strings.ToLower(cmd)

	if command, exists := ce.commands[cmd]; exists {
		return command.handler(args)
	}

	return fmt.Sprintf("invalid command: %s", cmd), nil
}

// NOTE: only allowing global config ATM
func (ce *CommandEngine) diagnose(args []string) (string, error) {
	if len(args) != 0 {
		return "", fmt.Errorf(ce.commands["diagnose"].usageMsg)
	}

	node := ce.shell.Find("station.ax")
	if node == nil {
		return "", fmt.Errorf("config station.ax not found")
	}
	parser := parser.NewParser(parser.NewParserConfig())

	err := parser.Parse([]byte(node.Read()))
	if err != nil {
		return "", err
	}

	errs := ce.state.ValidateConfig(parser.Config)
	if errs != nil {
		var errMsgs []string
		for _, err := range errs {
			errMsgs = append(errMsgs, err.Error())
		}
		errsMsgsJoined := strings.Join(errMsgs, "\n")

		return "", errors.New(errsMsgsJoined)
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
		return "", fmt.Errorf(ce.commands["cat"].usageMsg)
	}

	path := args[0]
	return ce.shell.Cat(path), nil
}

func (ce *CommandEngine) ls(args []string) (string, error) {
	if len(args) > 1 {
		return "", fmt.Errorf(ce.commands["ls"].usageMsg)
	}

	path := ""
	if len(args) == 1 {
		path = args[0]
	}

	return ce.shell.Ls(path), nil
}

func (ce *CommandEngine) tree(args []string) (string, error) {
	path := "."
	depth := 2

	switch len(args) {
	case 0:
		return "", fmt.Errorf(ce.commands["tree"].usageMsg)
	case 1:
		path = args[0]
	case 2:
		path = args[0]
		d, err := strconv.Atoi(args[1])
		if err != nil {
			return "", fmt.Errorf(ce.commands["tree"].usageMsg)
		}
		depth = d
	default:
		return "", fmt.Errorf(ce.commands["tree"].usageMsg)
	}

	return ce.shell.Tree(path, depth), nil
}

func (ce *CommandEngine) write(args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf(ce.commands["write"].usageMsg)
	}

	path := args[0]
	content := args[1]
	node := ce.shell.GetChild(path)

	if node == nil {
		return "", fmt.Errorf("%s does not exist", path)
	}

	node.Write(content)

	return fmt.Sprintf("successfully wrote %s", path), nil
}

func (ce *CommandEngine) inspect(args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf(ce.commands["inspect"].usageMsg)
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
		return "", fmt.Errorf(ce.commands["status"].usageMsg)
	}

	systemsDir := ce.shell.Find("systems")
	if systemsDir == nil {
		return "", fmt.Errorf("subsystems directory does not exist")
	}

	var subsystemTypes []*filesystem.Node
	for _, s := range systemsDir.Children() {
		subsystemTypes = append(subsystemTypes, s)
	}

	slices.SortFunc(subsystemTypes, func(a, b *filesystem.Node) int {
		return strings.Compare(a.Name(), b.Name())
	})

	var output string
	for _, subsystemTypeDir := range subsystemTypes {
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

func (ce *CommandEngine) help(args []string) (string, error) {

	if len(args) == 0 {
		var cmds []command
		for _, cmd := range ce.commands {
			cmds = append(cmds, cmd)
		}

		slices.SortFunc(cmds, func(a, b command) int {
			return strings.Compare(a.name, b.name)
		})

		var helpMsgs []string
		for _, cmd := range cmds {
			helpMsgs = append(helpMsgs, fmt.Sprintf("%s - %s", cmd.usageMsg, cmd.helpMsg))
		}

		return strings.Join(helpMsgs, "\n"), nil
	}

	if len(args) == 1 {
		if cmd, exists := ce.commands[args[0]]; exists {
			return fmt.Sprintf("%s - %s", cmd.usageMsg, cmd.helpMsg), nil
		}
	}

	return "", fmt.Errorf(ce.commands["help"].usageMsg)
}
