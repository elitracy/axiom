package systems_test

import (
	"testing"

	"github.com/elias/axiom/engine/systems"
)

func TestLifeSupportTick(t *testing.T) {

	const powerAvailable = 1.0
	const lethalO2 = 0.12
	const lethalCO2 = 0.08

	l := systems.NewLifeSupportSystem()

	startingO2 := l.Sensors()["o2"]
	startingCO2 := l.Sensors()["co2"]

	for range 25 {
		l.Tick(powerAvailable)
	}

	if l.Sensors()["o2"] >= startingO2 {
		t.Error("o2 should slowly decrease as the health of the life support decreases")
	}

	if l.Sensors()["co2"] <= startingCO2 {
		t.Error("co2 should slowly increase as the health of the life support decreases")
	}

	for range 100 {
		l.Tick(powerAvailable)
	}

	if l.Sensors()["o2"] > lethalO2 {
		t.Error("o2 should be below the lethal level.")
	}

	if l.Sensors()["co2"] <= lethalCO2 {
		t.Error("co2 should be above the lethal level.")
	}

}
