package systems_test

import (
	"testing"

	"github.com/elias/axiom/engine/systems"
)

func TestCoolantTick(t *testing.T) {
	const ambientTemp = 20
	const startingHeat = 350
	c := systems.NewCoolantSystem(ambientTemp)

	c.Tick(startingHeat)

	if c.Health() != 99 {
		t.Error("Health should degrade as the system ticks.")
	}

	if c.Sensors()["temp"] <= ambientTemp {
		t.Error("The coolant should heat up as it absorbs heat.")
	}

	if c.Sensors()["flow_rate"] >= 1.0 {
		t.Error("Flow rate should not exceed 1.0 (100%).")
	}

	if c.Sensors()["flow_rate"] <= 0.0 {
		t.Error("Coolant should be flowing to absorb heat.")
	}

	if c.Sensors()["pressure"] >= 80 {
		t.Error("Pressure should passively decrease for the player to patch.")
	}

	var deltas []float64
	for i := range 20 {
		prevTemp := c.Sensors()["temp"]
		c.Tick(startingHeat - float64(10*i))
		delta := c.Sensors()["temp"] - prevTemp
		deltas = append(deltas, delta)
	}

	for _, d := range deltas {
		if d <= 0 {
			t.Error("Temperature of coolant should rise each tick")
		}
	}

}
