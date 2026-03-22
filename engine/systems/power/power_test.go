package power

import (
	"testing"

	"github.com/elias/axiom/engine/materials"
	"github.com/elias/axiom/engine/systems"
	"github.com/elias/axiom/engine/systems/components"
	"github.com/stretchr/testify/assert"
)

func testMetalNeutral() materials.Metal {
	return materials.Metal{
		Name:                "Test Metal",
		MinTemperature:      0,
		MaxTemperature:      100,
		MinPressure:         0,
		MaxPressure:         100,
		HeatAbsorptionRate:  0,
		MaxTemperatureDelta: 0,
	}
}

func testMetalStandard() materials.Metal {
	return materials.Metal{
		Name:                "Test Metal",
		MinTemperature:      0.0,
		MaxTemperature:      100.0,
		MinPressure:         0.0,
		MaxPressure:         100.0,
		HeatAbsorptionRate:  0.05,
		MaxTemperatureDelta: 0.05,
	}
}

func testPower(metal materials.Metal) *PowerCore {

	power := &PowerCore{
		SystemCore:         systems.NewSystemCore("Test PowerCore"),
		housingMaterial:    metal,
		power:              components.NewComponent("Power", 1.0, 10.0, 100),
		fuel:               components.NewComponent("Fuel", 1.0, 1.0, 1.0),
		temperature:        components.NewComponent("Temperature", 0.0, metal.MinTemperature, metal.MaxTemperature),
		health:             components.NewHealthComponent(1.0),
		heatGenerationRate: 0.01,
		powerGrowthRate:    0.0,
	}

	return power
}

func testInput() PowerInput {
	return PowerInput{
		CoolantTemperature: 0.0,
		AmbientTemperature: 0.0,
	}
}

func TestPower_NoFuelLosesPower(t *testing.T) {
	power := testPower(testMetalNeutral())
	power.power.SetNorm(1.0)

	power.fuel.SetNorm(0.0)

	power.Tick(testInput())

	assert.Equal(t, power.fuel.Norm(), 0.0)
	assert.Equal(t, power.power.Norm(), 1.0-power.powerGrowthRate)

}

func TestPower_NoPowerDissipatesHeat(t *testing.T) {
	power := testPower(testMetalNeutral())
	power.temperature.SetNorm(1.0)
	power.power.SetNorm(0.0)

	testInput := PowerInput{
		AmbientTemperature: 0.0,
		CoolantTemperature: 100.0,
	}
	power.Tick(testInput)

	assert.Equal(t, power.power.Norm(), 0.0)
	assert.Equal(t, power.temperature.Norm(), 1.0-power.housingMaterial.HeatAbsorptionRate)
}

func TestPower_NoPowerHighAmbientTempKeepsHeat(t *testing.T) {
	power := testPower(testMetalNeutral())
	power.temperature.SetNorm(1.0)
	power.power.SetNorm(0.0)

	testInput := PowerInput{
		AmbientTemperature: 100.0,
		CoolantTemperature: 100.0,
	}

	power.Tick(testInput)

	assert.Equal(t, power.power.Norm(), 0.0)
	assert.Equal(t, power.temperature.Norm(), 1.0)
}

func TestPower_HighTemperatureDegradesHealth(t *testing.T) {
	power := testPower(testMetalNeutral())
	power.temperature.SetNorm(1.0)

	power.Tick(testInput())

	assert.Equal(t, power.health.Norm(), 1.0-healthLostPerTick)
}

func TestPower_NoHealthNoPower(t *testing.T) {
	power := testPower(testMetalNeutral())
	power.health.SetNorm(0.0)

	power.Tick(testInput())

	assert.Equal(t, power.power.Norm(), 0.0)
}

func TestPower_CoolantLowersTemperature(t *testing.T) {
	testMetal := testMetalNeutral()
	testMetal.HeatAbsorptionRate = 1.0
	testMetal.MaxTemperatureDelta = 1.0

	power := testPower(testMetal)
	power.temperature.SetNorm(1.0)
	power.heatGenerationRate = 0.0

	power.Tick(testInput())

	assert.Equal(t, power.temperature.Norm(), 0.0)
}

func TestPower_ProducingPowerGeneratesHeat(t *testing.T) {
	testMetal := testMetalNeutral()
	testMetal.HeatAbsorptionRate = 0.0
	testMetal.MaxTemperatureDelta = 1.0

	power := testPower(testMetal)
	power.temperature.SetNorm(0.0)
	power.heatGenerationRate = 0.5

	testInput := PowerInput{
		AmbientTemperature: 100.0,
		CoolantTemperature: 100.0,
	}

	power.Tick(testInput)

	assert.Equal(t, power.temperature.Norm(), 0.50)
}

func TestPower_ProducingPowerUsesFuel(t *testing.T) {
	testMetal := testMetalNeutral()

	power := testPower(testMetal)

	power.Tick(testInput())

	assert.Equal(t, power.fuel.Norm(), 1.0-fuelLostPerTick)
}
func TestPower_TickOutputTemperature(t *testing.T) {
	testMetal := testMetalNeutral()
	testMetal.HeatAbsorptionRate = 0.0
	testMetal.MaxTemperatureDelta = 1.0

	power := testPower(testMetal)
	power.temperature.SetNorm(0.0)
	power.heatGenerationRate = 0.5

	testInput := PowerInput{
		AmbientTemperature: 100.0,
		CoolantTemperature: 100.0,
	}

	output := power.Tick(testInput)

	assert.Equal(t, output.Temperature, 50.0)
}

func TestPower_TickOutputPower(t *testing.T) {
	testMetal := testMetalNeutral()
	testMetal.HeatAbsorptionRate = 0.0
	testMetal.MaxTemperatureDelta = 1.0

	power := testPower(testMetal)
	power.temperature.SetNorm(0.0)

	output := power.Tick(testInput())

	assert.Equal(t, output.Power, 100.0)
}

func TestPower_TickSteadyState(t *testing.T) {
	tests := []struct {
		name            string
		coolantTemp     float64
		ambientTemp     float64
		heatGenRate     float64
		powerGrowthRate float64
		startTempNorm   float64
		startPowerNorm  float64
		startFuelNorm   float64
		metal           materials.Metal
		ticks           int

		expectedTempDir  string // "up", "down", "stable"
		expectedPowerDir string // "up", "down", "stable"
	}{
		{"Full Fuel - No Coolant Effect", 999, 0, .01, .01, 0, 1, 1, testMetalStandard(), 100, "up", "stable"},
		{"Full Fuel - Weak Coolant Effect", 20, 0, .01, .01, 0, 1, 1, testMetalStandard(), 100, "up", "stable"},
		{"Full Fuel - Strong Coolant Effect", -20, 0, .01, .01, 0, 1, 1, testMetalStandard(), 100, "stable", "stable"},
		{"No Fuel - High Temp", 999, 0, .01, .01, 1, 1, 0, testMetalStandard(), 100, "down", "down"},
		{"No Fuel - No Power - High Temp - High Ambient", 999, 100, .01, .01, 1, 0, 0, testMetalStandard(), 100, "stable", "stable"},
		{"Full Fuel - Coolant Heat Equalibrium", 0, 100, .01, .01, 0.2, 1, 1, testMetalStandard(), 10, "stable", "stable"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			power := &PowerCore{
				SystemCore:         systems.NewSystemCore("Test Power"),
				housingMaterial:    tt.metal,
				power:              components.NewComponent("Test Power", tt.startPowerNorm, 0, 100),
				fuel:               components.NewComponent("Test Fuel", tt.startFuelNorm, 0, 1),
				temperature:        components.NewComponent("Test Temperature", tt.startTempNorm, tt.metal.MinTemperature, tt.metal.MaxTemperature),
				health:             components.NewHealthComponent(1),
				heatGenerationRate: tt.heatGenRate,
				powerGrowthRate:    tt.powerGrowthRate,
			}

			input := PowerInput{AmbientTemperature: tt.ambientTemp, CoolantTemperature: tt.coolantTemp}

			first := power.Tick(input)
			for range tt.ticks - 2 {
				power.Tick(input)
			}
			last := power.Tick(input)

			switch tt.expectedPowerDir {
			case "up":
				assert.Greater(t, last.Power, first.Power)
			case "down":
				assert.Less(t, last.Power, first.Power)
			case "stable":
				assert.InDelta(t, last.Power, first.Power, .001)
			}

			switch tt.expectedTempDir {
			case "up":
				assert.Greater(t, last.Temperature, first.Temperature)
			case "down":
				assert.Less(t, last.Temperature, first.Temperature)
			case "stable":
				assert.InDelta(t, last.Temperature, first.Temperature, .001)
			}

			t.Logf("TEMP: %v", last.Temperature)

		})
	}
}
