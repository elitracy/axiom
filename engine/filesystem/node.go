package filesystem

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

type Node struct {
	name     string
	parent   *Node
	children map[string]*Node

	createdAt time.Time
	updatedAt time.Time

	isDir    bool
	writable bool
	reader   func() string

	content string
}

func NewDir(path string) *Node {
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")

	name := parts[len(parts)-1]
	// if len(parts) == 1 && parts[0] == "" {
	// 	name = "/"
	// }

	node := &Node{
		name:      name,
		children:  make(map[string]*Node),
		writable:  false,
		isDir:     true,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}

	return node
}

func NewFile(path string) *Node {
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")

	node := &Node{
		name:      parts[len(parts)-1],
		writable:  true,
		isDir:     false,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}

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
func (n *Node) IsDir() bool          { return n.isDir }

func (n *Node) GetChild(path string) *Node {
	path = strings.Trim(path, "/")
	if path == "" {
		return n
	}

	parts := strings.Split(path, "/")

	nextChild := parts[0]
	remaining := strings.Join(parts[1:], "/")

	for _, child := range n.children {
		if child.name == nextChild {
			return child.GetChild(remaining)
		}
	}

	return nil

}

func (n *Node) ls(path string) string {
	path = strings.Trim(path, "/")

	if path == ".." {
		if n.Parent() == nil {
			return "Invalid Path"
		}

		return n.Parent().ls("")
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
	return child.ls(remaining)
}

func (n Node) read() string {
	if n.reader != nil {
		n.reader()
	}

	return n.content
}

func (n *Node) Write(content string) {
	if n.writable {
		n.content = content
		n.updatedAt = time.Now()
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
	output += n.updatedAt.Format(" Jan 01 15:04")
	output += fmt.Sprintf(" %s", n.Name())

	return output
}

func (n Node) pwd() string {
	if n.parent == nil {
		return n.Name()
	}
	return n.parent.pwd() + n.Name()
}
