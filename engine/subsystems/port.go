package subsystems

import (
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type PortID int64

var (
	currentPortID = 0
)

type PortType int

const (
	Input PortType = iota
	Output
)

func newPortID() PortID {
	id := currentPortID
	currentPortID++
	return PortID(id)
}

type Port interface {
	ID() PortID
	Name() string
	Component() *components.Component
	Subsystem() Subsystem
}

type CorePort struct {
	id        PortID
	name      string
	component *components.Component
	subsystem Subsystem
}

func (p CorePort) ID() PortID                       { return p.id }
func (p CorePort) Name() string                     { return p.name }
func (p CorePort) Component() *components.Component { return p.component }
func (p CorePort) Subsystem() Subsystem             { return p.subsystem }

func newCorePort(name string, subsystem Subsystem, component *components.Component) *CorePort {
	return &CorePort{
		id:        newPortID(),
		name:      name,
		component: component,
		subsystem: subsystem,
	}
}

type InputPort struct {
	*CorePort
	input *utils.Unit
}

func NewInputPort(name string, subsystem Subsystem, component *components.Component) *InputPort {
	return &InputPort{
		CorePort: newCorePort(name, subsystem, component),
	}
}

func (p InputPort) Input() *utils.Unit         { return p.input }
func (p *InputPort) SetInput(value utils.Unit) { *p.input = value }

type OutputPort struct {
	*CorePort
}

func NewOutputPort(name string, subsystem Subsystem, component *components.Component) *OutputPort {
	return &OutputPort{
		CorePort: newCorePort(name, subsystem, component),
	}
}
