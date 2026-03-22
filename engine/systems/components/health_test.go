package components_test

import (
	"testing"

	"github.com/elias/axiom/engine/systems"
	"github.com/elias/axiom/engine/systems/components"
	"github.com/stretchr/testify/assert"
)

func TestHealthStatus(t *testing.T) {
	health := components.NewHealthComponent(0.5)

	assert.Equal(t, health.Name(), "Health (%)")
	assert.Equal(t, health.Min(), 0.0)
	assert.Equal(t, health.Norm(), 0.5)
	assert.Equal(t, health.Value(), 0.5)
	assert.Equal(t, health.Max(), 1.0)
	assert.Equal(t, health.String(), "0: Health (%) 0.500 (50.00%)")

	health.SetNorm(0.0)
	assert.Equal(t, health.Status(), systems.Offline)

	health.SetNorm(0.3)
	assert.Equal(t, health.Status(), systems.Critical)

	health.SetNorm(0.7)
	assert.Equal(t, health.Status(), systems.Degraded)

	health.SetNorm(1.0)
	assert.Equal(t, health.Status(), systems.Online)
}
