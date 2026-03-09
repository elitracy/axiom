package systems_test

import (
	"testing"

	"github.com/elias/axiom/engine/systems"
)

func TestCoolantTick(t *testing.T) {
	const ambientTemp = 20
	p := systems.NewCoolantSystem(ambientTemp)

	p.Tick(350)

	if p.Health() != 99 {
		t.Error("Health should degrade as the system ticks.")
	}

	if p.Sensors()["temp"] <= ambientTemp {
		t.Error("The coolant should heat up as it absorbs heat.")
	}

	if p.Sensors()["flow_rate"] >= 1.0 {
		t.Error("Flow rate should not exceed 1.0 (100%).")
	}

	if p.Sensors()["flow_rate"] <= 0.0 {
		t.Error("Coolant should be flowing to absorb heat.")
	}
	t.Log(p.Sensors()["flow_rate"])
	t.Log(p.Sensors()["temp"])
	t.Log(p.Sensors()["pressure"])

	if p.Sensors()["pressure"] >= 80 {
		t.Error("Pressure should passively decrease for the player to patch.")
	}

}
