package filesystem

import (
	"fmt"
	"slices"
	"strings"
)

type Node interface {
	Name() string
	Parent() *DirNode
	SetParent(*DirNode)
}

type DirNode struct {
	name     string
	children map[string]Node
	parent   *DirNode
}

func NewDirNode(name string) *DirNode {
	return &DirNode{
		name:     name,
		children: make(map[string]Node),
	}
}

func (n *DirNode) Name() string                    { return n.name + "/" }
func (n *DirNode) AddChild(name string, node Node) { n.children[name] = node; node.SetParent(n) }
func (n *DirNode) Parent() *DirNode                { return n.parent }
func (n *DirNode) SetParent(node *DirNode)         { n.parent = node }

func (n *DirNode) Ls(path string) string {

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

	child, isChild := n.children[pathParts[0]]
	childDirNode, isDirNode := child.(*DirNode)

	if !isChild || !isDirNode {
		return "Invalid path"
	}

	remaining := strings.Join(pathParts[1:], "/")
	return childDirNode.Ls(remaining)

}

type FileNode struct {
	name   string
	parent *DirNode
	reader func() string
}

func (n *FileNode) Name() string            { return n.name }
func (n *FileNode) Parent() *DirNode        { return n.parent }
func (n *FileNode) SetParent(node *DirNode) { n.parent = node }
func NewFileNode(name string) *FileNode {
	return &FileNode{
		name: name,
	}

}
