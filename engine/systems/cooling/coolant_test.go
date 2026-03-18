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
		ThermalExpansionRate: 0.5,
	}
}

func testMetal() materials.Metal {
	return materials.Metal{
		Name:                "Test Metal",
		MinTemperature:      0.0,
		MaxTemperature:      100.0,
		MinPressure:         1.0,
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
	assert.Equal(t, coolant.volume, 1.0)
	assert.Equal(t, coolant.temperature, testFluid().MinTemperature)
	assert.Equal(t, coolant.viscosity, 1.0)
	assert.Equal(t, coolant.pressure, 0.0)
}

func TestCoolant_CalculatePressureTemperature(t *testing.T) {
	fluid := testFluid()
	fluid.ThermalExpansionRate = 1
	coolant := NewCoolantLoop(fluid, testMetal())

	pressure := coolant.calculatePressure()
	assert.Equal(t, pressure, 0.0)

	coolant.temperature.SetNorm(0.5)
	pressure = coolant.calculatePressure()
	assert.Equal(t, pressure, 0.5)

	coolant.temperature.SetNorm(1.0)
	pressure = coolant.calculatePressure()
	assert.Equal(t, pressure, 1.0)
}

func TestCoolant_CalculatePressureExpansion(t *testing.T) {
	fluid := testFluid()
	fluid.ThermalExpansionRate = 1.0
	coolant := NewCoolantLoop(fluid, testMetal())
	coolant.temperature.SetNorm(1.0)

	pressure := coolant.calculatePressure()
	assert.Equal(t, pressure, 0.0)

	fluid.ThermalExpansionRate = 0.5
	coolant = NewCoolantLoop(fluid, testMetal())
	pressure = coolant.calculatePressure()
	assert.Equal(t, pressure, 0.5)

	fluid.ThermalExpansionRate = 1.0
	coolant = NewCoolantLoop(fluid, testMetal())
	pressure = coolant.calculatePressure()
	assert.Equal(t, pressure, 1.0)
}

func TestCoolant_NoTemperatureDelta(t *testing.T) {
	fluid := testFluid()
	fluid.MaxTemperatureDelta = 0.0

	coolant := NewCoolantLoop(fluid, testMetal())

	var output CoolantOutput
	for range 10 {
		output = coolant.Tick(testInput)
	}

	assert.Equal(t, math.Round(output.Temperature), 0.0)
}

func TestCoolant_HighTemperatureDelta(t *testing.T) {
	fluid := testFluid()
	fluid.MaxTemperatureDelta = 0.5

	coolant := NewCoolantLoop(fluid, testMetal())

	output := coolant.Tick(testInput)

	assert.Equal(t, math.Round(output.Temperature), 50.0)

	output = coolant.Tick(testInput)
	assert.Equal(t, math.Round(output.Temperature), 100.0)
}

func TestCoolant_LowTemperatureDelta(t *testing.T) {
	fluid := testFluid()
	fluid.MaxTemperatureDelta = 0.01

	coolant := NewCoolantLoop(fluid, testMetal())

	var output CoolantOutput
	for range 5 {
		output = coolant.Tick(testInput)
	}

	assert.Equal(t, output.Temperature, 5.00)

	for range 5 {
		output = coolant.Tick(testInput)
	}
	assert.Equal(t, math.Round(output.Temperature), 10.0)
}

func TestCoolant_NoHeatAbsorption(t *testing.T) {
	fluid := testFluid()
	fluid.HeatAbsorptionRate = 0.0

	coolant := NewCoolantLoop(fluid, testMetal())

	var output CoolantOutput
	for range 10 {
		output = coolant.Tick(testInput)
	}

	assert.Equal(t, math.Round(output.Temperature), 0.0)
}

func TestCoolant_HighHeatAbsorption(t *testing.T) {
	fluid := testFluid()
	fluid.HeatAbsorptionRate = 0.01

	input := CoolantInput{LoadTemperature: 5.0}

	coolant := NewCoolantLoop(fluid, testMetal())

	var output CoolantOutput
	for range 5 {
		output = coolant.Tick(input)
	}

	assert.Equal(t, math.Round(output.Temperature), 25.00)

	for range 5 {
		output = coolant.Tick(input)
	}

	assert.Equal(t, math.Round(output.Temperature), 50.00)
}

func TestCoolant_LowHeatAbsorption(t *testing.T) {
	fluid := testFluid()
	fluid.HeatAbsorptionRate = 0.01

	input := CoolantInput{LoadTemperature: 5.0}
	coolant := NewCoolantLoop(fluid, testMetal())

	var output CoolantOutput
	for range 5 {
		output = coolant.Tick(input)
	}

	assert.Equal(t, output.Temperature, 0.05)

	for range 5 {
		output = coolant.Tick(input)
	}

	assert.Equal(t, output.Temperature, 0.10)
}
