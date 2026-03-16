package cooling_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/elias/axiom/engine/materials"
	"github.com/elias/axiom/engine/systems"
	"github.com/elias/axiom/engine/systems/cooling"
	"github.com/stretchr/testify/assert"
)

func testFluid() materials.Fluid {
	return materials.Fluid{
		Name:                    "Test Fluid",
		MinTemperature:          0.0,
		MaxTemperature:          100.0,
		MinViscosity:            0.5,
		MaxViscosity:            1.0,
		HeatAbsorptionRate:      0.01,
		ThermalConductivityRate: 0.05,
		ThermalExpansionRate:    0.5,
	}
}

func testMetal() materials.Metal {
	return materials.Metal{
		Name:                    "Test Metal",
		MinTemperature:          0.0,
		MaxTemperature:          100.0,
		MinPressure:             1.0,
		MaxPressure:             100.0,
		HeatAbsorptionRate:      0.01,
		ThermalConductivityRate: 0.01,
	}
}

var testInput = cooling.CoolantInput{
	LoadTemperature: 100.0,
}

var testOutput = cooling.CoolantInput{
	LoadTemperature: 100.0,
}

func TestCoolantCreation(t *testing.T) {
	coolant := cooling.NewCoolantLoop(testFluid(), testMetal())

	assert.Equal(t, coolant.ID(), systems.SystemID{ID: 0})
	assert.Equal(t, coolant.Name(), "Cooling Loop")
	assert.Equal(t, coolant.Status(), systems.Online)

	assert.Contains(t, coolant.String(), strconv.Itoa(coolant.ID().ID))
	assert.Contains(t, coolant.String(), coolant.Name())
}

func TestCoolant_HighThermalConductivity(t *testing.T) {
	fluid := testFluid()
	fluid.ThermalConductivityRate = 0.1

	coolant := cooling.NewCoolantLoop(fluid, testMetal())

	var output cooling.CoolantOutput
	for range 5 {
		output = coolant.Tick(testInput)
	}

	assert.Equal(t, math.Round(output.Temperature), 50.0)

	for range 5 {
		output = coolant.Tick(testInput)
	}
	assert.Equal(t, math.Round(output.Temperature), 100.0)
}

func TestCoolant_LowThermalConductivity(t *testing.T) {
	fluid := testFluid()
	fluid.ThermalConductivityRate = 0.001

	coolant := cooling.NewCoolantLoop(fluid, testMetal())

	var output cooling.CoolantOutput
	for range 5 {
		output = coolant.Tick(testInput)
	}

	assert.Equal(t, output.Temperature, .500)

	for range 5 {
		output = coolant.Tick(testInput)
	}
	assert.Equal(t, math.Round(output.Temperature), 1.00)
}

func TestCoolant_HighHeatAbsorption(t *testing.T) {
	fluid := testFluid()
	fluid.HeatAbsorptionRate = 0.01

	input := cooling.CoolantInput{LoadTemperature: 5.0}

	coolant := cooling.NewCoolantLoop(fluid, testMetal())

	var output cooling.CoolantOutput
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

	input := cooling.CoolantInput{LoadTemperature: 5.0}
	coolant := cooling.NewCoolantLoop(fluid, testMetal())

	var output cooling.CoolantOutput
	for range 5 {
		output = coolant.Tick(input)
	}

	assert.Equal(t, output.Temperature, 0.05)

	for range 5 {
		output = coolant.Tick(input)
	}

	assert.Equal(t, output.Temperature, 0.10)
}
