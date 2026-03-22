package cooling

import (
	"testing"

	"github.com/elias/axiom/engine/materials"
	"github.com/elias/axiom/engine/systems"
	"github.com/stretchr/testify/assert"
)

func testFluid() materials.Fluid {
	return materials.Fluid{
		Name:                 "Test Fluid",
		MinTemperature:       0.0,
		MaxTemperature:       100.0,
		MinViscosity:         0.1,
		MaxViscosity:         10.0,
		HeatAbsorptionRate:   0.01,
		MaxTemperatureDelta:  0.02,
		ThermalExpansionRate: 0.1,
	}
}

func testMetal() materials.Metal {
	return materials.Metal{
		Name:                "Test Metal",
		MinTemperature:      0.0,
		MaxTemperature:      100.0,
		MinPressure:         0.0,
		MaxPressure:         100.0,
		HeatAbsorptionRate:  0.005,
		MaxTemperatureDelta: 0.01,
	}
}

var testInput = CoolantInput{
	LoadTemperature: 100.0,
}

var testOutput = CoolantInput{
	LoadTemperature: 100.0,
}

func TestCoolantCreation(t *testing.T) {
	coolant := NewCoolantLoop(testFluid(), testMetal())

	assert.Equal(t, coolant.ID(), systems.SystemID{ID: 0})
	assert.Equal(t, coolant.Name(), "Cooling Loop")
	assert.Equal(t, coolant.Status(), systems.Online)
	assert.Equal(t, coolant.volume.Norm(), 1.0)
	assert.Equal(t, coolant.temperature.Norm(), 0.01)
	assert.Equal(t, coolant.viscosity.Norm(), 0.99)
	assert.InDelta(t, coolant.pressure.Norm(), 0.051, 0.001)
}

func TestCoolant_CalculateViscosity(t *testing.T) {
	fluid := testFluid()
	coolant := NewCoolantLoop(fluid, testMetal())

	coolant.temperature.SetNorm(1.0)
	viscosity := coolant.calculateViscosityNorm()
	assert.Equal(t, viscosity, 0.0)

	coolant.temperature.SetNorm(0.5)
	viscosity = coolant.calculateViscosityNorm()
	assert.Equal(t, viscosity, 0.5)

	coolant.temperature.SetNorm(0.0)
	viscosity = coolant.calculateViscosityNorm()
	assert.Equal(t, viscosity, 1.0)
}

func TestCoolant_HighHeatDegradesVolumeOfFluid(t *testing.T) {
	coolant := NewCoolantLoop(testFluid(), testMetal())
	coolant.temperature.SetNorm(1.0)

	var testInput = CoolantInput{
		LoadTemperature: 100.0,
	}

	coolant.Tick(testInput)

	assert.Equal(t, coolant.volume.Norm(), 1.0-volumeLossPerTick)
}

func TestCoolant_TemperatureDelta(t *testing.T) {
	tests := []struct {
		name         string
		maxTempDelta float64
		expectedTemp float64
	}{
		{"No delta", 0.0, 0.0},
		{"Low delta", 0.01, 1.0},
		{"Med delta", 0.1, 10.0},
		{"High delta", 0.5, 50.0},
		{"Max delta", 1.0, 100.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fluid := materials.Fluid{
				Name:                 "Test Fluid",
				MinTemperature:       0.0,
				MaxTemperature:       100.0,
				MinViscosity:         0.1,
				MaxViscosity:         10.0,
				HeatAbsorptionRate:   1,
				MaxTemperatureDelta:  tt.maxTempDelta,
				ThermalExpansionRate: 0.1,
			}

			metal := materials.Metal{
				Name:                "Test Metal",
				MinTemperature:      0.0,
				MaxTemperature:      100.0,
				MinPressure:         0.0,
				MaxPressure:         100.0,
				HeatAbsorptionRate:  0.005,
				MaxTemperatureDelta: tt.maxTempDelta,
			}

			input := CoolantInput{
				LoadTemperature: 100,
			}

			coolant := NewCoolantLoop(fluid, metal)
			coolant.viscosity.SetNorm(0.0)
			coolant.temperature.SetNorm(0.0)
			coolant.pressure.SetNorm(1.0)
			coolant.basePressure = 0.0

			output := coolant.Tick(input)

			assert.InDelta(t, output.Temperature, tt.expectedTemp, 0.001)

		})
	}
}

func TestCoolant_HeatAbsorptionRate(t *testing.T) {
	tests := []struct {
		name               string
		heatAbsorptionRate float64
		expectedTemp       float64
	}{
		{"No absorption", 0.0, 0.0},
		{"Low absorption", 0.2, 20.0},
		{"Med absorption", 0.5, 50.0},
		{"High absorption", 1.0, 100.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fluid := materials.Fluid{
				Name:                 "Test Fluid",
				MinTemperature:       0.0,
				MaxTemperature:       100.0,
				MinViscosity:         0.1,
				MaxViscosity:         10.0,
				HeatAbsorptionRate:   tt.heatAbsorptionRate,
				MaxTemperatureDelta:  1,
				ThermalExpansionRate: 0.1,
			}

			metal := materials.Metal{
				Name:                "Test Metal",
				MinTemperature:      0.0,
				MaxTemperature:      100.0,
				MinPressure:         0.0,
				MaxPressure:         100.0,
				HeatAbsorptionRate:  0.005,
				MaxTemperatureDelta: 1,
			}

			input := CoolantInput{
				LoadTemperature: 100,
			}

			coolant := NewCoolantLoop(fluid, metal)
			coolant.viscosity.SetNorm(0.0)
			coolant.temperature.SetNorm(0.0)
			coolant.pressure.SetNorm(1.0)
			coolant.basePressure = 0.0

			output := coolant.Tick(input)

			assert.InDelta(t, output.Temperature, tt.expectedTemp, 0.001)

		})
	}
}

func TestCoolant_CalculatePressure(t *testing.T) {
	tests := []struct {
		name          string
		tempNorm      float64
		expansionRate float64
		basePressure  float64
		expected      float64
	}{
		{"Zero temp, zero base", 0.0, 1.0, 0.0, 0.0},
		{"Half temp, full expansion", 0.5, 1.0, 0.0, 0.5},
		{"Full temp, full expansion", 1.0, 1.0, 0.0, 1.0},
		{"Full temp, zero expansion", 1.0, 0.0, 0.0, 0.0},
		{"Full temp, half expansion", 1.0, 0.5, 0.0, 0.5},
		{"Half temp, half expansion", 0.5, 0.5, 0.0, 0.25},
		{"With base pressure", 0.0, 1.0, 0.1, 0.1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fluid := materials.Fluid{
				Name:                 "Test Fluid",
				MinTemperature:       0.0,
				MaxTemperature:       100.0,
				MinViscosity:         0.1,
				MaxViscosity:         10.0,
				HeatAbsorptionRate:   0.1,
				MaxTemperatureDelta:  0.01,
				ThermalExpansionRate: tt.expansionRate,
			}

			metal := materials.Metal{
				Name:                "Test Metal",
				MinTemperature:      0.0,
				MaxTemperature:      100.0,
				MinPressure:         0.0,
				MaxPressure:         100.0,
				HeatAbsorptionRate:  0.005,
				MaxTemperatureDelta: 1,
			}

			coolant := NewCoolantLoop(fluid, metal)
			coolant.basePressure = tt.basePressure
			coolant.temperature.SetNorm(tt.tempNorm)

			pressure := coolant.calculatePressureNorm()

			assert.InDelta(t, pressure, tt.expected, 0.001)
		})
	}
}

func TestCoolant_TickSteadyState(t *testing.T) {
	tests := []struct {
		name            string
		loadTemp        float64
		startTempNorm   float64
		ticks           int
		expectedTempDir string
	}{
		{"Hot load from cold", 100, 0.1, 50, "up"},
		{"Cold load from hot", 0, 1.0, 50, "stable"},
		{"No load", 0, 0.0, 50, "stable"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fluid := materials.Fluid{
				Name:                 "Test Fluid",
				MinTemperature:       0.0,
				MaxTemperature:       100.0,
				MinViscosity:         0.1,
				MaxViscosity:         10.0,
				HeatAbsorptionRate:   0.1,
				MaxTemperatureDelta:  0.01,
				ThermalExpansionRate: 0.1,
			}

			metal := materials.Metal{
				Name:                "Test Metal",
				MinTemperature:      0.0,
				MaxTemperature:      100.0,
				MinPressure:         0.0,
				MaxPressure:         100.0,
				HeatAbsorptionRate:  0.005,
				MaxTemperatureDelta: 1,
			}

			input := CoolantInput{
				LoadTemperature: tt.loadTemp,
			}

			coolant := NewCoolantLoop(fluid, metal)
			coolant.temperature.SetNorm(tt.startTempNorm)

			first := coolant.Tick(input)
			for range tt.ticks - 2 {

				coolant.Tick(input)
			}
			last := coolant.Tick(input)

			switch tt.expectedTempDir {
			case "up":
				assert.Greater(t, last.Temperature, first.Temperature)
			case "down":
				assert.Less(t, last.Temperature, first.Temperature)
			case "stable":
				assert.InDelta(t, last.Temperature, first.Temperature, 0.001)
			}

		})
	}
}
