package filesystem

import (
	"fmt"
	"log"
	"slices"
	"strings"
)

type Node struct {
	name     string
	children map[string]*Node
	parent   *Node

	content  string
	writable bool
	reader   func() string
}

func NewNode(path string) *Node {

	isDir := path[len(path)-1] == '/'

	parts := strings.Split(path, "/")
	log.Printf("%v", parts)

	node := &Node{
		name:     parts[len(parts)-2],
		writable: true,
	}

	if isDir {
		node.children = make(map[string]*Node)

	}
	log.Printf(node.String())

	return node
}

func (n *Node) Name() string {
	if n.IsDir() {
		return n.name + "/"
	}

	return n.name
}
func (n *Node) AddChild(node *Node)  { n.children[node.Name()] = node; node.SetParent(n) }
func (n *Node) Parent() *Node        { return n.parent }
func (n *Node) SetParent(node *Node) { n.parent = node }
func (n *Node) IsDir() bool          { return n.children == nil }

func (n *Node) GetChild(path string) *Node {
	if n.name == path {
		return n
	}

	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return nil
	}

	nextChild := parts[0]
	remaining := strings.Join(parts[1:], "/")

	for _, child := range n.children {
		if child.name == nextChild {
			return child.GetChild(remaining)
		}
	}

	return nil

}

func (n *Node) Ls(path string) string {
	log.Printf("CHILDREN: %v", n.children)

	if path == ".." {
		if n.Parent() == nil {
			return "Invalid Path"
		}

		return n.Parent().Ls("")
	}

	if path == "" || path == "." {

		children := []*Node{}
		for _, child := range n.children {
			children = append(children, child)
		}
		slices.SortFunc(children, func(a, b *Node) int {
			return strings.Compare(a.name, b.name)

		})

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

func (n Node) Read() string {
	if n.reader != nil {
		n.reader()
	}

	return n.content
}

func (n Node) Write(content string) {
	if n.writable {
		n.content = content
	}
}

func (n Node) String() string {
	output := ""
	if n.IsDir() {
		output += "d"
	} else {
		output += "."
	}

	output += "r"

	if n.writable {
		output += "w-"
	} else {
		output += "--"
	}

	output += fmt.Sprintf(" axiom")
	output += fmt.Sprintf(" %s", n.name)

	return output
}
