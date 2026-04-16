package filesystem

import (
	"strings"

	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/simulation"
	"github.com/elias/axiom/engine/utils"
)

type Shell struct {
	cwd  *Node
	root *Node
}

type worldState interface {
	Subsystems() []simulation.Subsystem
}

func NewShell() *Shell {
	return &Shell{}
}

func (s *Shell) Populate(ws worldState) {
	if ws == nil {
		return
	}

	root := NewDir("/")
	sys := NewDir("sys")
	usr := NewDir("usr")
	root.AddChild(sys)
	root.AddChild(usr)

	conf := NewDir("conf")
	bin := NewDir("bin")
	usr.AddChild(conf)
	usr.AddChild(bin)

	logs := NewDir("logs")
	systems := NewDir("systems")
	sys.AddChild(systems)
	sys.AddChild(logs)

	power := NewDir("power")
	cooling := NewDir("cooling")
	machines := NewDir("machines")

	systems.AddChild(power)
	systems.AddChild(cooling)
	systems.AddChild(machines)

	s.root = root
	s.cwd = root

	logging.Debug("SUBS: %v", ws.Subsystems())

	for _, subsystem := range ws.Subsystems() {
		dir := NewDir(subsystem.Name())
		status := NewFile("status")
		components := NewDir("components")
		dir.AddChild(status)
		dir.AddChild(components)

		switch subsystem.Type() {
		case utils.Power:
			power.AddChild(dir)
		case utils.Cooling:
			cooling.AddChild(dir)
		case utils.Machine:
			machines.AddChild(dir)
		}
	}
}

func (s *Shell) Ls(path string) string {
	if path == "" {
		return s.cwd.ls(path)
	}

	if path[0] == '/' {
		s.root.ls(path)
	}

	return s.cwd.ls(path)
}

func (s *Shell) Cd(path string) {

	if path == "." {
		return
	}

	if path == ".." && s.cwd.Parent() != nil {
		s.cwd = s.cwd.Parent()
	}

	node := s.cwd.GetChild(path)

	if node != nil {
		s.cwd = node
	}

}

func (s Shell) Cat(path string) string {
	path = strings.Trim(path, "/")
	node := s.cwd.GetChild(path)

	if node == nil {
		return ""
	}

	return node.read()
}

func (s Shell) Pwd() string {
	return s.cwd.pwd()
}
