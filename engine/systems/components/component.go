package components

import (
	"fmt"

	"github.com/elias/axiom/engine/utils"
)

type ComponentID struct {
	id int64
}

type Component interface {
	ID() ComponentID
	Name() string
	Value() float64
	SetValue(float64)
	Min() float64
	Max() float64
	ApplyCurve(float64) float64
	String() string
}

// ComponentCore is the core struct of all components.
type ComponentCore struct {
	id         ComponentID
	name       string
	value      float64
	min        float64
	max        float64
	valueCurve func(float64) float64
}

// Creates a new component
// value is the starting value for the component, min is the minimum value, max is the maximum value, valueCurve is the function to calcuate the current value
func NewComponentCore(name string, value, min, max float64, valueCurve ...func(float64) float64) *ComponentCore {

	var curve func(float64) float64
	if len(valueCurve) > 0 {
		curve = valueCurve[0]
	}

	return &ComponentCore{
		// TODO: generate component IDS dynamically
		id:         ComponentID{id: 0},
		name:       name,
		value:      value,
		min:        min,
		max:        max,
		valueCurve: curve,
	}

}

// Returns the component's id
func (c *ComponentCore) ID() ComponentID { return c.id }

// Returns the component's name
func (c *ComponentCore) Name() string { return c.name }

// Returns the component's current value
func (c *ComponentCore) Value() float64 { return c.value }

// Sets the components value
func (c *ComponentCore) SetValue(value float64) { c.value = utils.Clamp(c.min, value, c.max) }

// Returns the minimum value for the component
func (c *ComponentCore) Min() float64 { return c.min }

// Returns the maximum value for the component
func (c *ComponentCore) Max() float64 { return c.max }

// Returns the current component value on the valueCurve function
func (c *ComponentCore) ApplyValueCurve() float64 {
	value := c.Value()
	if c.valueCurve != nil {

		value = c.valueCurve(c.Value())
	}

	return utils.Clamp(c.Min(), value, c.Max())
}

// Returns the stringified info of the component
func (c *ComponentCore) String() string {
	output := fmt.Sprintf("%v: %v (%.2f, %.2f, %.2f) %.2f", c.ID().id, c.Name(), c.Min(), c.Value(), c.Max(), c.ApplyValueCurve())
	return output
}
