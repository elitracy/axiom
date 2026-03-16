package components

import (
	"fmt"

	"github.com/elias/axiom/engine/utils"
)

type ComponentID struct {
	id int64
}

type ComponentReader interface {
	ID() ComponentID
	Name() string
	Min() float64
	Max() float64
	Norm() float64
	Value() float64
	String() string
}

type Component interface {
	ComponentReader
	SetNorm(float64)
}

// ComponentCore is the core struct of all components.
type ComponentCore struct {
	id        ComponentID
	name      string
	min       float64
	max       float64
	value     float64
	normValue float64
	curve     func(float64) float64
}

// Creates a new component
// initialValue is the starting value for the component, min is the minimum value, max is the maximum value, valueCurve is the function to calcuate the current value
func NewComponent(name string, normInitialValue, appliedMin, appliedMax float64, appliedValueCurve ...func(float64) float64) *ComponentCore {

	var curve func(float64) float64
	if len(appliedValueCurve) > 0 {
		curve = appliedValueCurve[0]
	}

	comp := &ComponentCore{
		// TODO: generate component IDS dynamically
		id:    ComponentID{id: 0},
		name:  name,
		min:   appliedMin,
		max:   appliedMax,
		curve: curve,
	}

	comp.SetNorm(normInitialValue)

	return comp

}

// Returns the component's id
func (c *ComponentCore) ID() ComponentID { return c.id }

// Returns the component's name
func (c *ComponentCore) Name() string { return c.name }

// Returns the minimum value for the component
func (c *ComponentCore) Min() float64 { return c.min }

// Returns the maximum value for the component
func (c *ComponentCore) Max() float64 { return c.max }

// Returns the component's current value
func (c *ComponentCore) Value() float64 { return c.value }

// Returns the component's value normalized by the component's min and max
func (c *ComponentCore) Norm() float64 { return c.normValue }

// Sets the components applied value with a normalized value
func (c *ComponentCore) SetNorm(normValue float64) {
	c.normValue = utils.Clamp(0.0, normValue, 1.0)
	if c.curve != nil {
		c.value = c.curve(normValue)
	} else {
		c.value = normValue
	}

	c.value = c.value*(c.Max()-c.Min()) + c.Min()
	c.value = utils.Clamp(c.min, c.value, c.max)
}

// Returns the stringified info of the component
func (c *ComponentCore) String() string {
	output := fmt.Sprintf("%d: %s %.3f (%.2f%%)", c.ID().id, c.Name(), c.Value(), c.Norm()*100)
	return output
}
