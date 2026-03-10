package systems_test

import (
	"testing"

	"github.com/elias/axiom/engine/systems"
)

func TestLifeSupportTick(t *testing.T) {

	const powerAvailable = 1.0

	l := systems.NewLifeSupportSystem()

	l.Tick(powerAvailable)
	l.Tick(powerAvailable)
	l.Tick(powerAvailable)
	l.Tick(powerAvailable)
	l.Tick(powerAvailable)
	l.Tick(powerAvailable)

	if l.Health() != 99 {
		t.Error("Health should degrade as the system ticks.")
	}

	t.Log(l.Sensors())

}
