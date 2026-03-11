package simulation_test

import (
	"testing"

	"github.com/elias/axiom/engine/simulation"
)

func TestSimulationTick(t *testing.T) {
	state := simulation.NewWorldState()

	for range 100 {
		state.Tick()
	}
	t.Log(state)

	if state.Power().Health() >= 1 {
		t.Error("THe power should die with no user intervention.")
	}

	if state.Coolant().Health() >= 1 {
		t.Error("The coolant system should die with no user intervention.")
	}

	if state.LifeSupport().Health() >= 1 {
		t.Error("The life support system should die with no user intervention.")
	}

}
