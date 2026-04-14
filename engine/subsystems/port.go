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

type CorePort struct {
	id        PortID
	name      string
	subsystem Subsystem
}

func (p CorePort) ID() PortID           { return p.id }
func (p CorePort) Name() string         { return p.name }
func (p CorePort) Subsystem() Subsystem { return p.subsystem }

func newCorePort(id PortID, name string, subsystem Subsystem) CorePort {
	return CorePort{
		id:        id,
		name:      name,
		subsystem: subsystem,
	}
}

type InputPort struct {
	CorePort
	value    utils.Unit
	channel  string
	received bool
}

func (p InputPort) String() string {
	return fmt.Sprintf("%s[%d] %s", p.Name(), p.ID(), p.channel)
}

func newInputPort(id PortID, name string, subsystem Subsystem, channel string) *InputPort {
	return &InputPort{
		CorePort: newCorePort(id, name, subsystem),
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
	CorePort
	component *components.Component
}

func (p OutputPort) String() string {
	return fmt.Sprintf("%s[%d] %s", p.Name(), p.ID(), p.component.Name())
}

func newOutputPort(id PortID, name string, subsystem Subsystem, component *components.Component) *OutputPort {
	return &OutputPort{
		CorePort:  newCorePort(id, name, subsystem),
		component: component,
	}
}

func (p *OutputPort) Component() *components.Component { return p.component }
