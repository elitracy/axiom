package simulation

import "github.com/elias/axiom/engine/systems"

type WorldState struct {
	tick        int64
	power       systems.PowerSystem
	coolant     systems.CoolantSystem
	lifeSupport systems.LifeSupportSystem

	ambientTemp float64
}

func (w *WorldState) Tick() {
	w.coolant.Tick(w.power.Sensors()["temp"])
	w.power.Tick(w.power.Sensors()["flowRate"], w.ambientTemp)
	w.lifeSupport.Tick(w.power.Sensors()["outputLevel"])
	w.tick++
}
