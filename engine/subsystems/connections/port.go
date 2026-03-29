package connections

import "github.com/elias/axiom/engine/subsystems/components"

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
}

func NewPort(name string, component *components.Component) *Port {
	return &Port{
		id:        newPortID(),
		name:      name,
		component: component,
	}
}
