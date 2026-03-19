package machines

import (
	"testing"

	"github.com/elias/axiom/engine/systems"
	"github.com/stretchr/testify/assert"
)

const ambientTemp = 20.0

func hvacInput() HvacInput {
	return HvacInput{
		PowerAvailable: 1200,
		HeatSources:    []float64{},
	}
}

func TestHVAC_NewHvac(t *testing.T) {
	hvac := NewHvac(20.0)

	assert.Equal(t, hvac.maxTemperatureDelta, 0.02)
	assert.Equal(t, hvac.targetTemperature, 20.0)
	assert.Equal(t, hvac.requiredPower, 1200.0)
	assert.Equal(t, hvac.temperature.Value(), ambientTemp)
	assert.Equal(t, hvac.health.Status(), systems.Online)
}

func TestHVAC_StatusOffline(t *testing.T) {
	hvac := NewHvac(-50.0)

	output := hvac.Tick(hvacInput())
	assert.Equal(t, output.Status, systems.Offline)
}

func TestHVAC_StatusCritical(t *testing.T) {
	hvac := NewHvac(40.0)

	output := hvac.Tick(hvacInput())
	assert.Equal(t, output.Status, systems.Critical)
	t.Logf("TEMP: %.2f", output.Temperature)

	hvac = NewHvac(0.0)

	output = hvac.Tick(hvacInput())
	t.Logf("TEMP: %.2f", output.Temperature)
	assert.Equal(t, output.Status, systems.Critical)
}
