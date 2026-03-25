package components

import "github.com/elias/axiom/engine/utils"

type ComponentType int64
type ComponentID int64

var currentComponentID ComponentID = 0

func newComponentID() ComponentID {
	id := currentComponentID
	currentComponentID++

	return id
}

const (
	Temperature ComponentType = iota
	Effort
	Power
)

type Component struct {
	id            ComponentID
	componentType ComponentType
	value         utils.Norm
}

func NewComponent(componentType ComponentType, value utils.Norm) *Component {
	return &Component{
		id:            newComponentID(),
		componentType: componentType,
		value:         value,
	}
}

func (c Component) Type() ComponentType                   { return c.componentType }
func (c Component) Value() utils.Norm                     { return c.value }
func (c *Component) SetValue(value utils.Norm) *Component { c.value = value; return c }

func (c *Component) AddValue(value utils.Norm) *Component { c.value += value; return c }
