package machines

import (
	"math"
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

	assert.Equal(t, hvac.maxTempPush, 0.01)
	assert.Equal(t, hvac.targetTemperature, 20.0)
	assert.Equal(t, hvac.requiredPower, 1200.0)
	assert.Equal(t, hvac.temperature.Value(), ambientTemp)
	assert.Equal(t, hvac.health.Status(), systems.Online)
}

func TestHvacStatus(t *testing.T) {
	tests := []struct {
		name      string
		upperTemp float64
		lowerTemp float64
		expected  systems.Status
	}{
		{"Hvac Status (offline)", 100, -100, systems.Offline},
		{"Hvac Status (critical)", 45, 0, systems.Critical},
		{"Hvac Status (degraded)", 35, 15, systems.Degraded},
		{"Hvac Status (online)", 29, 21, systems.Online},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hvac := NewHvac(tt.upperTemp)
			got := hvac.Status()
			assert.Equal(t, tt.expected, got)

			hvac = NewHvac(tt.lowerTemp)
			got = hvac.Status()
			assert.Equal(t, tt.expected, got)

		})
	}
}

func TestCalculateEffectiveness(t *testing.T) {
	tests := []struct {
		name           string
		powerAvailable float64
		powerRequired  float64
		expected       float64
	}{
		{"neg power", -1000, 1200, 0.0},
		{"no power", 0, 1200, 0.0},
		{"bad power", 100, 1200, 0.0},
		{"low power", 600, 1200, 0.5},
		{"medium power", 900, 1200, 0.75},
		{"high power", 1200, 1200, 1.0},
		{"excess power", 2000, 1200, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateEffectiveness(tt.powerAvailable, tt.powerRequired)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestCalculateHeatPush(t *testing.T) {
	tests := []struct {
		name        string
		heatLoad    float64
		targetTemp  float64
		heatingRate float64
		expected    float64
	}{
		{"Heat Load (excess)", 200.0, 0.0, 0.1, 20},
		{"Heat Load (max)", 100.0, 0.0, 0.1, 10},
		{"Heat Load (mid)", 40., 0.0, 0.1, 4},
		{"Heat Load (min)", 0.0, 0.0, 0.1, 0},
		{"Heat Load (neg)", -100.0, 0.0, 0.1, -10},

		{"Target Temp (excess)", 100.0, 200.0, 0.1, -10},
		{"Target Temp (max)", 100.0, 100.0, 0.1, 0},
		{"Target Temp (mid)", 100.0, 50.0, 0.1, 5},
		{"Target Temp (min)", 100.0, 0.0, 0.1, 10},
		{"Target Temp (neg)", 100.0, -100.0, 0.1, 20},

		{"Heating Rate (execess)", 100.0, 0.0, 2.0, 200},
		{"Heating Rate (max)", 100.0, 0.0, 1.0, 100},
		{"Heating Rate (mid)", 100.0, 0.0, 0.5, 50},
		{"Heating Rate (min)", 100.0, 0.0, 0.1, 10},
		{"Heating Rate (neg)", 100.0, 0.0, -0.5, -50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateHeatPush(tt.heatLoad, tt.targetTemp, tt.heatingRate)
			assert.InDelta(t, tt.expected, got, .001)
		})
	}
}

func TestCalculateHeatPush_Properties(t *testing.T) {
	// Symmetry
	pushUp := calculateHeatPush(100.0, 0.0, 0.1)
	pushDown := calculateHeatPush(-100.0, 0.0, 0.1)

	assert.InDelta(t, math.Abs(pushUp), math.Abs(pushDown), 0.0001)
	assert.Greater(t, pushUp, 0.0)
	assert.Less(t, pushDown, 0.0)

	// Monotonicity

	pushBig := calculateHeatPush(100, 0.0, 0.1)
	pushSmall := calculateHeatPush(10, 0.0, 0.1)

	assert.Greater(t, pushBig, pushSmall)

	// Equilibrium
	push := calculateHeatPush(20.0, 20.0, 0.1)
	assert.Zero(t, push)

	// Clamping
	// N/A
}

func TestCalculateAverageHeat(t *testing.T) {
	tests := []struct {
		name        string
		heatSources []float64
		expected    float64
	}{
		{"Average Heat (large difference)", []float64{-100, 100}, 0},
		{"Average Heat (small difference)", []float64{0, 10}, 5},
		{"Average Heat (no difference)", []float64{10, 10}, 10},
		{"Average Heat (no heat)", []float64{}, 0},
		{"Average Heat (one heat)", []float64{100}, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateAverageHeat(tt.heatSources)
			assert.InDelta(t, tt.expected, got, .001)
		})
	}
}

func TestHvacTick_NoPower_DriftToHeat(t *testing.T) {
	hvac := NewHvac(20.0)
	input := hvacInput()
	input.HeatSources = append(input.HeatSources, 100, 150, 200)
	input.PowerAvailable = 0

	out1 := hvac.Tick(input)
	out2 := hvac.Tick(input)
	assert.Greater(t, out2.Temperature, out1.Temperature)

	hvac = NewHvac(20.0)
	input = hvacInput()
	input.HeatSources = append(input.HeatSources, -100, -150, -200)
	input.PowerAvailable = 0

	out1 = hvac.Tick(input)
	out2 = hvac.Tick(input)
	assert.Less(t, out2.Temperature, out1.Temperature)
}

func TestHvacTick_LessPower_ConvergesLessEfficiently(t *testing.T) {
	hvacHighPower := NewHvac(20.0)
	inputHighPower := hvacInput()
	inputHighPower.HeatSources = append(inputHighPower.HeatSources, 100)

	hvacHighPower.Tick(inputHighPower)
	outHighPower := hvacHighPower.Tick(inputHighPower)

	hvacLowPower := NewHvac(20.0)
	inputLowPower := hvacInput()
	inputLowPower.PowerAvailable = 1200 * hvacEfficiencyMediumThreshold
	inputLowPower.HeatSources = append(inputLowPower.HeatSources, 100)

	hvacLowPower.Tick(inputLowPower)
	outLowPower := hvacLowPower.Tick(inputLowPower)

	assert.Greater(t, outLowPower.Temperature, outHighPower.Temperature)
}

func TestHvacTick_PowerHeatTargetTemps(t *testing.T) {

	tests := []struct {
		name           string
		heatSource     float64
		powerAvailable float64
		expected       float64
	}{
		// Low heat — all power levels hold
		{"High Power - Low Heat", 35, 1200, 21.0},
		{"Med Power - Low Heat", 35, 600, 22.0},
		{"Low Power - Low Heat", 35, 300, 23.0},

		// Med heat — low power starts struggling
		{"High Power - Med Heat", 70, 1200, 23.0},
		{"Med Power - Med Heat", 70, 600, 26.0},
		{"Low Power - Med Heat", 70, 300, 30.5},

		// High heat — med power struggles too
		{"High Power - High Heat", 110, 1200, 25.5},
		{"Med Power - High Heat", 110, 600, 30.5},
		{"Low Power - High Heat", 110, 300, 39.0},

		// Extreme heat — everyone hurts
		{"High Power - Extreme Heat", 190, 1200, 30.5},
		{"Med Power - Extreme Heat", 190, 600, 40.0},
		{"Low Power - Extreme Heat", 190, 300, 50.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hvac := NewHvac(20.0)
			input := hvacInput()
			input.HeatSources = append(input.HeatSources, tt.heatSource)
			input.PowerAvailable = tt.powerAvailable

			for range 100 {
				hvac.Tick(input)
			}

			assert.InDelta(t, hvac.temperature.Value(), tt.expected, 2)
		})

	}

}
