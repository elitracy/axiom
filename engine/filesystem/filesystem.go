package filesystem

type Shell struct {
	cwd *Node
}

func NewShell(root *Node) *Shell {
	return &Shell{
		cwd: root,
	}
}

func (s *Shell) Ls(path string) string {
	return s.cwd.Ls(path)
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
	node := s.cwd.GetChild(path)

	if node == nil {
		return ""
	}

	return node.Read()
}
