package subsystems

import (
	"fmt"

	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type PortID int64

var (
	currentPortID = 0
)

type PortType int

const (
	PortInput PortType = iota
	PortOutput
)

func newPortID() PortID {
	id := currentPortID
	currentPortID++
	return PortID(id)
}

type Port interface {
	ID() PortID
	Name() string
	Subsystem() Subsystem
	String() string
}

type CorePort struct {
	id        PortID
	name      string
	subsystem Subsystem
}

func (p CorePort) ID() PortID           { return p.id }
func (p CorePort) Name() string         { return p.name }
func (p CorePort) Subsystem() Subsystem { return p.subsystem }

func newCorePort(name string, subsystem Subsystem, component *components.Component) *CorePort {
	return &CorePort{
		id:        newPortID(),
		name:      name,
		subsystem: subsystem,
	}
}

type InputPort struct {
	*CorePort
	value    utils.Unit
	channel  string
	received bool
}

func (p InputPort) String() string {
	return fmt.Sprintf("%s[%d] %s", p.Name(), p.ID(), p.channel)
}

func NewInputPort(name string, subsystem Subsystem, channel string) *InputPort {
	return &InputPort{
		CorePort: newCorePort(name, subsystem, nil),
		channel:  channel,
	}
}

func (p *InputPort) Clear() {
	p.value = 0
	p.received = false
}

func (p *InputPort) SetValue(value utils.Unit) {
	p.value = value
	p.received = true
}

type OutputPort struct {
	*CorePort
	component *components.Component
}

func (p OutputPort) String() string {
	return fmt.Sprintf("%s[%d] %s", p.Name(), p.ID(), p.component.Name())
}

func NewOutputPort(name string, subsystem Subsystem, component *components.Component) *OutputPort {
	return &OutputPort{
		CorePort:  newCorePort(name, subsystem, component),
		component: component,
	}
}

func (p *OutputPort) Component() *components.Component { return p.component }
