package components_test

import (
	"testing"

	"github.com/elias/axiom/engine/systems/components"
	"github.com/elias/axiom/engine/utils"
	"github.com/stretchr/testify/assert"
)

func TestComponentNormalized(t *testing.T) {
	component := components.NewComponent("Test Component", 0.5, 0.0, 1.0)

	assert.Equal(t, component.Name(), "Test Component")
	assert.Equal(t, component.Min(), 0.0)
	assert.Equal(t, component.Norm(), 0.5)
	assert.Equal(t, component.Value(), 0.5)
	assert.Equal(t, component.Max(), 1.0)
	assert.Equal(t, component.String(), "0: Test Component 0.500 (50.00%)")

	component.SetNorm(0.8)
	assert.Equal(t, component.Norm(), 0.8)
	assert.Equal(t, component.Value(), 0.8)

	component.SetNorm(5.0)
	assert.Equal(t, component.Norm(), 1.0)
	assert.Equal(t, component.Value(), 1.0)

	component.SetNorm(-5.0)
	assert.Equal(t, component.Norm(), 0.0)
	assert.Equal(t, component.Value(), 0.0)
}

func TestComponentClamp(t *testing.T) {
	component := components.NewComponent("Test Component", 0.5, 0.0, 1.0)

	component.SetNorm(5.0)
	assert.Equal(t, component.Norm(), 1.0)
	assert.Equal(t, component.Value(), 1.0)
	assert.Equal(t, component.String(), "0: Test Component 1.000 (100.00%)")

	component.SetNorm(-5.0)
	assert.Equal(t, component.Norm(), 0.0)
	assert.Equal(t, component.Value(), 0.0)
	assert.Equal(t, component.String(), "0: Test Component 0.000 (0.00%)")
}

func TestComponentCurve(t *testing.T) {
	curve := func(value float64) float64 {
		return utils.Sigmoid(value, 0.5, 10)
	}
	component := components.NewComponent("Test Component", 0.5, 0.0, 1.0, curve)

	assert.Equal(t, component.Name(), "Test Component")
	assert.Equal(t, component.Min(), 0.0)
	assert.Equal(t, component.Norm(), 0.5)
	assert.Equal(t, component.Value(), 0.5)
	assert.Equal(t, component.Max(), 1.0)

	component.SetNorm(1.0)
	assert.Equal(t, component.Norm(), 1.0)
	assert.Equal(t, component.Value(), curve(1.0))

	component.SetNorm(0.0)
	assert.Equal(t, component.Norm(), 0.0)
	assert.Equal(t, component.Value(), curve(0.0))
}

func TestComponentScaling(t *testing.T) {
	component := components.NewComponent("Test Component", 0.5, -100, 100)

	assert.Equal(t, component.Name(), "Test Component")
	assert.Equal(t, component.Min(), -100.0)
	assert.Equal(t, component.Norm(), 0.5)
	assert.Equal(t, component.Value(), 0.0)
	assert.Equal(t, component.Max(), 100.0)
}

func TestComponentCurveScaling(t *testing.T) {
	curve := func(value float64) float64 {
		return utils.Sigmoid(value, 0.5, 10)
	}
	component := components.NewComponent("Test Component", 0.5, -100, 100, curve)
	scale := func(x float64) float64 { return x*(component.Max()-component.Min()) + component.Min() }

	assert.Equal(t, component.Name(), "Test Component")
	assert.Equal(t, component.Min(), -100.0)
	assert.Equal(t, component.Norm(), 0.5)
	assert.Equal(t, component.Value(), scale(curve(0.5)))
	assert.Equal(t, component.Max(), 100.0)

	component.SetNorm(0.0)
	assert.Equal(t, component.Value(), scale(curve(0.0)))

	component.SetNorm(1.0)
	assert.Equal(t, component.Value(), scale(curve(1.0)))
}
