package power

import (
	"testing"

	"github.com/elias/axiom/engine/materials"
	"github.com/elias/axiom/engine/systems"
	"github.com/elias/axiom/engine/systems/components"
	"github.com/stretchr/testify/assert"
)

func testMetal() materials.Metal {
	return materials.Metal{
		Name:                "Test Metal",
		MinTemperature:      0.0,
		MaxTemperature:      100.0,
		MinPressure:         0.0,
		MaxPressure:         100.0,
		HeatAbsorptionRate:  0.01,
		MaxTemperatureDelta: 0.01,
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
	power := testPower(testMetal())
	power.power.SetNorm(1.0)

	power.fuel.SetNorm(0.0)

	power.Tick(testInput())

	assert.Equal(t, power.fuel.Norm(), 0.0)
	assert.Equal(t, power.power.Norm(), 1.0-1.0/systems.TICKS_TILL_DEATH_DEBUG)

}

func TestPower_NoPowerDissipatesHeat(t *testing.T) {
	power := testPower(testMetal())
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
	power := testPower(testMetal())
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
	power := testPower(testMetal())
	power.temperature.SetNorm(1.0)

	power.Tick(testInput())

	assert.Equal(t, power.health.Norm(), 1.0-percentHealthLostPerTick)
}

func TestPower_NoHealthNoPower(t *testing.T) {
	power := testPower(testMetal())
	power.health.SetNorm(0.0)

	power.Tick(testInput())

	assert.Equal(t, power.power.Norm(), 0.0)
}

func TestPower_CoolantLowersTemperature(t *testing.T) {
	testMetal := testMetal()
	testMetal.HeatAbsorptionRate = 1.0
	testMetal.MaxTemperatureDelta = 1.0

	power := testPower(testMetal)
	power.temperature.SetNorm(1.0)
	power.heatGenerationRate = 0.0

	power.Tick(testInput())

	assert.Equal(t, power.temperature.Norm(), 0.0)
}

func TestPower_ProducingPowerGeneratesHeat(t *testing.T) {
	testMetal := testMetal()
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
	testMetal := testMetal()

	power := testPower(testMetal)

	power.Tick(testInput())

	assert.Equal(t, power.fuel.Norm(), 1.0-percentFuelLostPerTick)
}
