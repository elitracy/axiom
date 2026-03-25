package subsystems

import (
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type Hvac struct {
	*subsystemCore
	targetTemp utils.Norm
}

func NewHvac() *Hvac {
	hvac := &Hvac{
		subsystemCore: newSubsystemCore("Cooling"),
		targetTemp:    utils.Norm(.2),
	}
	hvac.AddComponent(components.Temperature, hvac.targetTemp)

	return hvac
}

func (s *Hvac) Effort() utils.Norm { return s.components[components.Effort].Value() }

func (s *Hvac) Tick(inputs map[components.ComponentType]*components.Component) {
	power := utils.Norm(0.0)
	if effort, ok := inputs[components.Effort]; ok {
		power = effort.Value()
	}

	heat := utils.Norm(0.0)
	if temp, ok := inputs[components.Temperature]; ok {
		heat = temp.Value()
	}

	currentTemp := s.components[components.Temperature].Value()
	delta := calcHvacHeatDelta(s.targetTemp, heat, power)

	if currentTemp != s.targetTemp {
		s.components[components.Temperature].AddValue(delta)
	}

}

func calcHvacHeatDelta(targetTemp, heat, rate utils.Norm) utils.Norm {
	cooling := (heat - targetTemp) * -rate
	return cooling
}
