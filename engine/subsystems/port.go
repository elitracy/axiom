package subsystems

import (
	"github.com/elias/axiom/engine/subsystems/components"
)

type PortID int64

var (
	currentPortID = 0
)

func newPortID() PortID {
	id := currentPortID
	currentPortID++
	return PortID(id)
}

type Port struct {
	id        PortID
	name      string
	component *components.Component
	subsystem Subsystem
}

func (p Port) ID() PortID                       { return p.id }
func (p Port) Name() string                     { return p.name }
func (p Port) Component() *components.Component { return p.component }
func (p Port) Subsystem() Subsystem { return p.subsystem}

func NewPort(name string, subsystem Subsystem,  component *components.Component) *Port {
	return &Port{
		id:        newPortID(),
		name:      name,
		component: component,
		subsystem: subsystem,
	}
}
