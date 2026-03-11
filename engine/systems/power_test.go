package systems_test

import (
	"testing"

	"github.com/elias/axiom/engine/systems"
)

func TestPowerTick(t *testing.T) {
	ambientTemp := 20.0
	coolantFlowRate := 1.5

	p := systems.NewPowerSystem(ambientTemp)
	p.Tick(coolantFlowRate, ambientTemp)

	if p.Sensors()["temp"] >= ambientTemp {
		t.Error("Temperature is expected to decrease if coolant flow is greater than heat rate.")
	}

	if p.Sensors()["fuel_level"] != .99 {
		t.Error("Fuel should be consumed as power is generated.")
	}

	if p.Sensors()["output_level"] != 1.0 {
		t.Error("OUtput level should not be changed (for now)")
	}
}
