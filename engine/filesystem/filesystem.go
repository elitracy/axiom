package filesystem

import (
	"strings"
)

type Shell struct {
	cwd  *Node
	root *Node
}

func NewShell(root *Node) *Shell {
	return &Shell{
		cwd:  root,
		root: root,
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
