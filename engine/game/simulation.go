package game

import "github.com/elias/axiom/engine/utils"

type Simulation interface {
	Update(tick *utils.Tick)
}

func RunSimulation(sim Simulation, startTick *utils.Tick) {
	sim.Update(startTick)
	startTick.Val.Add(1)

}
