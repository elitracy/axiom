package subsystems

import (
	"fmt"

	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type PortID int64

type PortType int

const (
	PortInput PortType = iota
	PortOutput
)

type port struct {
	id   PortID
	name string
}

type OutputPort struct {
	port
	component *components.Component
}

type InputPort struct {
	port
	value    utils.Unit
	channel  string
	received bool
}

func newPort(id PortID, name string) port {
	return port{
		id:   id,
		name: name,
	}
}

func newInputPort(id PortID, name string, channel string) *InputPort {
	return &InputPort{
		port:    newPort(id, name),
		channel: channel,
	}
}

func newOutputPort(id PortID, name string, component *components.Component) *OutputPort {
	return &OutputPort{
		port:      newPort(id, name),
		component: component,
	}
}

func (p port) ID() PortID   { return p.id }
func (p port) Name() string { return p.name }

func (p InputPort) String() string {
	return fmt.Sprintf("%s[%d] %s", p.Name(), p.ID(), p.channel)
}

func (p *InputPort) Clear() {
	p.value = 0
	p.received = false
}

func (p *InputPort) SetValue(value utils.Unit) {
	p.value = value
	p.received = true
}

func (p OutputPort) String() string {
	return fmt.Sprintf("%s[%d] %s", p.Name(), p.ID(), p.component.Name())
}

func (p *OutputPort) Component() *components.Component { return p.component }
