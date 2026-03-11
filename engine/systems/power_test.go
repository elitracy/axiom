package systems_test

import (
	"testing"

	"github.com/elias/axiom/engine/systems"
)

func TestPowerOverheating(t *testing.T) {
	ambientTemp := 20.0
	coolantFlowRate := 0.8

	p := systems.NewPowerSystem(ambientTemp)
	p.Tick(coolantFlowRate, ambientTemp)

	if p.Sensors()["temp"] <= ambientTemp {
		t.Error("Temperature should rise while coolant flow < 1")
	}

	if p.Sensors()["fuel_level"] >= 1.0 {
		t.Error("Fuel should be consumed as power is generated.")
	}

	if p.Sensors()["output_level"] >= 1.0 {
		t.Error("Output level should follow sigmoid curve")
	}
}

func TestPowerCooling(t *testing.T) {
	ambientTemp := 20.0
	coolantFlowRate := 1.1

	p := systems.NewPowerSystem(ambientTemp)
	p.Tick(coolantFlowRate, ambientTemp)

	if p.Sensors()["temp"] >= ambientTemp {
		t.Error("Temperature should lower while coolant flow > 1")
	}

	if p.Sensors()["fuel_level"] >= 1.0 {
		t.Error("Fuel should be consumed as power is generated.")
	}

	if p.Sensors()["output_level"] >= 1.0 {
		t.Error("Output level should follow sigmoid curve")
	}
}
