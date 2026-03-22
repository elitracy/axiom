package machines

import (
	"math"
	"testing"

	"github.com/elias/axiom/engine/systems"
	"github.com/stretchr/testify/assert"
)

const ambientTemp = 25.0

func hvacInput() HvacInput {
	return HvacInput{
		PowerSupplied: 1200,
		HeatSources:   []float64{},
	}
}

func TestHVAC_NewHvac(t *testing.T) {
	hvac := NewHvac(25.0)

	assert.Equal(t, hvac.targetTemperature, 25.0)
	assert.Equal(t, hvac.powerCapacity, 1200.0)
	assert.InDelta(t, hvac.temperature.Value(), ambientTemp, 0.001)
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
		name          string
		powerSupplied float64
		powerRequired float64
		expected      float64
	}{
		{"neg power", -1000, 1200, 0},
		{"no power", 0, 1200, 0},
		{"low power", 600, 1200, 0.25},
		{"medium power", 900, 1200, 0.5625},
		{"high power", 1200, 1200, 1},
		{"excess power", 2000, 1200, 1},

		{"zero capacity", 100, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateEffectiveness(tt.powerSupplied, tt.powerRequired)
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
	hvac := NewHvac(25.0)
	input := hvacInput()
	input.HeatSources = append(input.HeatSources, 100, 150, 200)
	input.PowerSupplied = 0

	out1 := hvac.Tick(input)
	out2 := hvac.Tick(input)
	assert.Greater(t, out2.Temperature, out1.Temperature)

	hvac = NewHvac(25.0)
	input = hvacInput()
	input.HeatSources = append(input.HeatSources, -100, -150, -200)
	input.PowerSupplied = 0

	out1 = hvac.Tick(input)
	out2 = hvac.Tick(input)
	assert.Less(t, out2.Temperature, out1.Temperature)
}

func TestHvacTick_LessPower_ConvergesLessEfficiently(t *testing.T) {
	hvacHighPower := NewHvac(25.0)
	inputHighPower := hvacInput()
	inputHighPower.HeatSources = append(inputHighPower.HeatSources, 100)

	hvacHighPower.Tick(inputHighPower)
	outHighPower := hvacHighPower.Tick(inputHighPower)

	hvacLowPower := NewHvac(25.0)
	inputLowPower := hvacInput()
	inputLowPower.PowerSupplied = 1200 * hvacEfficiencyMediumThreshold
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
		{"20c - Low Power", 20, 300, 22},
		{"20c - Med Power", 20, 600, 23},
		{"20c - High Power", 20, 1200, 24},

		{"30c - Low Power", 30, 300, 29},
		{"30c - Med Power", 30, 600, 26},
		{"30c - High Power", 30, 1200, 25},

		{"50c - Low Power", 50, 300, 38},
		{"50c - Med Power", 50, 600, 30},
		{"50c - High Power", 50, 1200, 27},

		{"80c - Low Power", 80, 300, 50},
		{"80c - Med Power", 80, 600, 36},
		{"80c - High Power", 80, 1200, 28},

		{"140c - Low Power", 140, 300, 50},
		{"140c - Med Power", 140, 600, 50},
		{"140c - High Power", 140, 1200, 31},

		{"200c - Low Power", 200, 300, 50},
		{"200c - Med Power", 200, 600, 50},
		{"200c - High Power", 200, 1200, 36},

		{"500c - Low Power", 500, 300, 50},
		{"500c - Med Power", 500, 600, 50},
		{"500c - High Power", 500, 1200, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hvac := NewHvac(25.0)
			input := hvacInput()
			input.HeatSources = append(input.HeatSources, tt.heatSource)
			input.PowerSupplied = tt.powerAvailable

			for range 100 {
				hvac.Tick(input)
			}

			assert.InDelta(t, hvac.temperature.Value(), tt.expected, 2)
		})
	}

}
