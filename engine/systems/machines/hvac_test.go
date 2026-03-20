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

	assert.Equal(t, hvac.maxTemperatureNormDelta, 0.02)
	assert.Equal(t, hvac.targetTemperature, 20.0)
	assert.Equal(t, hvac.requiredPower, 1200.0)
	assert.Equal(t, hvac.temperature.Value(), ambientTemp)
	assert.Equal(t, hvac.health.Status(), systems.Online)
}

func TestHVAC_StatusOffline(t *testing.T) {
	hvac := NewHvac(-50.0)
	output := hvac.Tick(hvacInput())
	assert.Equal(t, output.Status, systems.Offline)

	hvac = NewHvac(100.0)
	output = hvac.Tick(hvacInput())
	assert.Equal(t, output.Status, systems.Offline)
}

func TestHVAC_StatusCritical(t *testing.T) {
	hvac := NewHvac(45.0)
	output := hvac.Tick(hvacInput())
	assert.Equal(t, output.Status, systems.Critical)

	hvac = NewHvac(5.0)
	output = hvac.Tick(hvacInput())
	assert.Equal(t, output.Status, systems.Critical)
}

func TestHVAC_StatusDegraded(t *testing.T) {
	hvac := NewHvac(35.0)
	output := hvac.Tick(hvacInput())
	assert.Equal(t, output.Status, systems.Degraded)

	hvac = NewHvac(15.0)
	output = hvac.Tick(hvacInput())
	assert.Equal(t, output.Status, systems.Degraded)
}

func TestHVAC_TickEffectivenessDropsWithPower(t *testing.T) {
	hvac := NewHvac(ambientTemp)

	input := hvacInput()
	input.PowerAvailable = hvac.requiredPower * 0.0
	input.HeatSources = append(input.HeatSources, 30.0)

	output := hvac.Tick(input)

	assert.Equal(t, output.Temperature, ambientTemp*(1+heatRate))
}
