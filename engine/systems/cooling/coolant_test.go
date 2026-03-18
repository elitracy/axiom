package cooling

import (
	"math"
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
	assert.Equal(t, coolant.temperature.Norm(), testFluid().MinTemperature)
	assert.Equal(t, coolant.viscosity.Norm(), 1.0)
	assert.Equal(t, coolant.pressure.Norm(), 0.05)
}

func TestCoolant_CalculatePressureTemperature(t *testing.T) {
	fluid := testFluid()
	fluid.ThermalExpansionRate = 1
	coolant := NewCoolantLoop(fluid, testMetal())
	coolant.basePressure = 0.0

	pressure := coolant.calculatePressureNorm()
	assert.Equal(t, pressure, 0.0)

	coolant.temperature.SetNorm(0.5)
	pressure = coolant.calculatePressureNorm()
	assert.Equal(t, pressure, 0.5)

	coolant.temperature.SetNorm(1.0)
	pressure = coolant.calculatePressureNorm()
	assert.Equal(t, pressure, 1.0)
}

func TestCoolant_CalculatePressureExpansion(t *testing.T) {
	fluid := testFluid()
	fluid.ThermalExpansionRate = 0.0
	coolant := NewCoolantLoop(fluid, testMetal())
	coolant.basePressure = 0.0
	coolant.temperature.SetNorm(1.0)

	pressure := coolant.calculatePressureNorm()
	assert.Equal(t, pressure, 0.0)

	fluid.ThermalExpansionRate = 0.5
	coolant = NewCoolantLoop(fluid, testMetal())
	coolant.basePressure = 0.0
	coolant.temperature.SetNorm(1.0)
	pressure = coolant.calculatePressureNorm()
	assert.Equal(t, pressure, 0.5)

	fluid.ThermalExpansionRate = 1.0
	coolant = NewCoolantLoop(fluid, testMetal())
	coolant.basePressure = 0.0
	coolant.temperature.SetNorm(1.0)
	pressure = coolant.calculatePressureNorm()
	assert.Equal(t, pressure, 1.0)
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

func TestCoolant_NoTemperatureDelta(t *testing.T) {
	fluid := testFluid()
	fluid.HeatAbsorptionRate = 1.0
	fluid.MaxTemperatureDelta = 0.0

	coolant := NewCoolantLoop(fluid, testMetal())
	coolant.viscosity.SetNorm(0.0)
	coolant.temperature.SetNorm(0.0)
	coolant.pressure.SetNorm(1.0)
	coolant.basePressure = 0.0

	output := coolant.Tick(testInput)

	assert.Equal(t, math.Round(output.Temperature), 0.0)
}

func TestCoolant_LowTemperatureDelta(t *testing.T) {
	fluid := testFluid()
	fluid.HeatAbsorptionRate = 1.0
	fluid.MaxTemperatureDelta = 0.01

	coolant := NewCoolantLoop(fluid, testMetal())
	coolant.viscosity.SetNorm(0.0)
	coolant.temperature.SetNorm(0.0)
	coolant.pressure.SetNorm(1.0)
	coolant.basePressure = 0.0

	output := coolant.Tick(testInput)

	assert.Equal(t, math.Round(output.Temperature), 1.0)
}

func TestCoolant_HighTemperatureDelta(t *testing.T) {
	fluid := testFluid()
	fluid.HeatAbsorptionRate = 1.0
	fluid.MaxTemperatureDelta = 0.5

	coolant := NewCoolantLoop(fluid, testMetal())
	coolant.viscosity.SetNorm(0.0)
	coolant.temperature.SetNorm(0.0)
	coolant.pressure.SetNorm(1.0)
	coolant.basePressure = 0.0

	output := coolant.Tick(testInput)

	assert.Equal(t, math.Round(output.Temperature), 50.0)
}

func TestCoolant_NoHeatAbsorptionRate(t *testing.T) {
	fluid := testFluid()
	fluid.HeatAbsorptionRate = 0.0
	fluid.MaxTemperatureDelta = 1.0

	coolant := NewCoolantLoop(fluid, testMetal())

	output := coolant.Tick(testInput)

	assert.Equal(t, math.Round(output.Temperature), 0.0)
}

func TestCoolant_LowHeatAbsorptionRate(t *testing.T) {
	fluid := testFluid()
	fluid.HeatAbsorptionRate = 0.2
	fluid.MaxTemperatureDelta = 1.0
	fluid.ThermalExpansionRate = 1.0

	coolant := NewCoolantLoop(fluid, testMetal())
	coolant.temperature.SetNorm(0.0)
	coolant.viscosity.SetNorm(0.0)
	coolant.pressure.SetNorm(1.0)

	var testInput = CoolantInput{
		LoadTemperature: 100.0,
	}
	output := coolant.Tick(testInput)

	assert.Equal(t, math.Round(output.Temperature), 20.0)
}

func TestCoolant_HighHeatAbsorptionRate(t *testing.T) {
	fluid := testFluid()
	fluid.HeatAbsorptionRate = 1.0
	fluid.MaxTemperatureDelta = 1.0
	fluid.ThermalExpansionRate = 1.0

	coolant := NewCoolantLoop(fluid, testMetal())
	coolant.temperature.SetNorm(0.0)
	coolant.viscosity.SetNorm(0.0)
	coolant.pressure.SetNorm(1.0)

	var testInput = CoolantInput{
		LoadTemperature: 100.0,
	}
	output := coolant.Tick(testInput)

	assert.Equal(t, math.Round(output.Temperature), 100.0)
}
