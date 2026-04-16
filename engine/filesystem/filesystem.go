package filesystem

import (
	"strings"

	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/subsystems/components"
)

type Shell struct {
	cwd  *Node
	root *Node
}

type worldState interface {
	Subsystems() []subsystem
}

type subsystem interface {
	Name() string
	Components() map[string]*components.Component
}

func NewShell() *Shell {
	return &Shell{}
}

func (s *Shell) Populate(ws worldState) {

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

	systems.AddChild(NewDir("power"))
	systems.AddChild(NewDir("cooling"))
	systems.AddChild(NewDir("machines"))

	s.root = root
	s.cwd = root
}

func (s *Shell) Ls(path string) string {
	if path == "" {
		return s.cwd.ls(path)
	}

	if path[0] == '/' {
		logging.Debug("HERE")
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
