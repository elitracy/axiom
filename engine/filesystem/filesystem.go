package filesystem

import (
	"fmt"
	"log"
	"slices"
	"strings"
)

type Shell struct {
	cwd *Node
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

	dir, remaining, _ := strings.Cut(path, "/")

	node, exists := s.cwd.children[dir]

	// bad cd
	if !exists {
		return
	}

	// at last point
	if remaining == "" {
		s.cwd = node
	}

	if exists && node.IsDir() {
		if remaining == "" {
			s.Cd(path)
			return
		}

		s.Cd(dir)
	}

}

type Node struct {
	name     string
	children map[string]*Node
	parent   *Node

	reader func() string
}

func NewNode(path string) *Node {

	isDir := path[len(path)-1] == '/'

	parts := strings.Split(path, "/")

	node := &Node{
		name: parts[len(parts)-1],
	}

	if isDir {
		log.Printf("NODE: %v", node.name)
		node.children = make(map[string]*Node)

	}

	return node
}

func (n *Node) Name() string         { return n.name + "/" }
func (n *Node) AddChild(node *Node)  { n.children[node.Name()] = node; node.SetParent(n) }
func (n *Node) Parent() *Node        { return n.parent }
func (n *Node) SetParent(node *Node) { n.parent = node }
func (n *Node) IsDir() bool          { return n.children == nil }

func (n *Node) Ls(path string) string {

	if path == ".." {
		if n.Parent() == nil {
			return "Invalid Path"
		}

		return n.Parent().Ls("")
	}

	if path == "" || path == "." {
		children := []string{}
		for _, child := range n.children {
			children = append(children, child.Name())
		}

		slices.Sort(children)

		output := ""
		for _, child := range children {
			output += fmt.Sprintf("%s\n", child)
		}

		return output
	}

	pathParts := strings.Split(path, "/")

	child, exists := n.children[pathParts[0]]

	if !exists {
		return "Invalid path"
	}

	remaining := strings.Join(pathParts[1:], "/")
	return child.Ls(remaining)
}
