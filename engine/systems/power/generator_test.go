package power

import (
	"testing"

	"github.com/elias/axiom/engine/systems"
	"github.com/stretchr/testify/assert"
)

func TestGenerator_NewGenerator(t *testing.T) {

	generator := NewGenerator(testMetalNeutral(), 20.0)

	assert.Equal(t, generator.Power().Norm(), 0.0)
	assert.Equal(t, generator.Fuel().Norm(), 1.0)
	assert.Equal(t, generator.Temperature().Norm(), 0.2)
	assert.Equal(t, generator.Status(), systems.Online)
}

func TestGenerator_Tick(t *testing.T) {
	generator := NewGenerator(testMetalNeutral(), 20.0)

	input := testInput()

	output := generator.Tick(input)

	assert.InDelta(t, output.Power, 120.0, 0.001)
	assert.InDelta(t, output.Temperature, 20.1, 0.1)
}
