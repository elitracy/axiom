package components

import (
	"fmt"

	"github.com/elias/axiom/engine/utils"
)

type ComponentType int64
type ComponentID int64

//go:generate stringer -type=ComponentType
const (
	Temperature ComponentType = iota
	Power
)

type Component struct {
	id            ComponentID
	name          string
	componentType ComponentType
	value         utils.Unit
}

func NewComponent(id ComponentID, name string, componentType ComponentType, value utils.Unit) *Component {
	return &Component{
		id:            id,
		name:          name,
		componentType: componentType,
		value:         value,
	}
}

func (c Component) ID() ComponentID                       { return c.id }
func (c Component) Type() ComponentType                   { return c.componentType }
func (c Component) Name() string                          { return c.name }
func (c Component) Value() utils.Unit                     { return c.value }
func (c *Component) SetValue(value utils.Unit) *Component { c.value = value.Clamp(); return c }

func (c *Component) AddValue(value utils.Unit) *Component {
	c.value += value
	c.value = c.value.Clamp()

	return c
}

func (c *Component) String() string {
	return fmt.Sprintf("%s[%s]: %.2f", c.Name(), c.Type(), c.Value())
}
