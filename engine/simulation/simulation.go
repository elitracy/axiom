package simulation

import (
	"fmt"

	"github.com/elias/axiom/engine/systems"
)

const (
	ambientTemp = 25.0
)

type WorldState struct {
	tick        int64
	power       *systems.PowerSystem
	coolant     *systems.CoolantSystem
	lifeSupport *systems.LifeSupportSystem

	ambientTemp float64
}

func (s *WorldState) Power() *systems.PowerSystem             { return s.power }
func (s *WorldState) Coolant() *systems.CoolantSystem         { return s.coolant }
func (s *WorldState) LifeSupport() *systems.LifeSupportSystem { return s.lifeSupport }

func NewWorldState() *WorldState {
	return &WorldState{
		tick:        0,
		power:       systems.NewPowerSystem(ambientTemp),
		coolant:     systems.NewCoolantSystem(ambientTemp),
		lifeSupport: systems.NewLifeSupportSystem(),
		ambientTemp: ambientTemp,
	}
}

func (w *WorldState) Tick() {
	powerTemp := w.power.Sensors()["temp"]
	coolantFlowRate := w.coolant.Sensors()["flow_rate"]
	powerLevel := w.power.Sensors()["output_level"]

	w.coolant.Tick(powerTemp)
	w.power.Tick(coolantFlowRate, w.ambientTemp)
	w.lifeSupport.Tick(powerLevel)
	w.tick++
}

func (w *WorldState) String() string {
	output := fmt.Sprintf("[tick:%v] %v", w.tick, w.power)
	output += fmt.Sprintf("\n[tick:%v] %v", w.tick, w.coolant)
	output += fmt.Sprintf("\n[tick:%v] %v", w.tick, w.lifeSupport)
	return output
}
