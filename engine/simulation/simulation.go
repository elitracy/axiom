package simulation

import (
	"github.com/elias/axiom/engine"
	"github.com/elias/axiom/engine/subsystems"
	"github.com/elias/axiom/engine/subsystems/components"
	"github.com/elias/axiom/engine/utils"
)

type subsystemConnection struct {
	component *components.Component
	id        subsystems.SubsystemID
}

type WorldState struct {
	subsystems   map[subsystems.SubsystemID]subsystems.Subsystem
	dependencies map[subsystems.SubsystemID][]subsystemConnection
}

func (ws *WorldState) addSubsystem(subsystem subsystems.Subsystem) {
	ws.subsystems[subsystem.ID()] = subsystem
	ws.dependencies[subsystem.ID()] = make([]subsystemConnection, 0)
}

func (ws *WorldState) addDependency(subsystem, dep subsystems.Subsystem, compType components.ComponentType) {
	connection := subsystemConnection{
		id:        dep.ID(),
		component: dep.Components()[compType],
	}

	ws.dependencies[subsystem.ID()] = append(ws.dependencies[subsystem.ID()], connection)
}

func (ws *WorldState) Init() {
	power := subsystems.NewPower(.5)
	cooling := subsystems.NewCooling(.5)
	hvac := subsystems.NewHvac()

	ws.addSubsystem(power)
	ws.addSubsystem(cooling)
	ws.addSubsystem(hvac)

	ws.addDependency(hvac, power, components.Power)
	ws.addDependency(hvac, power, components.Temperature)
	ws.addDependency(power, cooling, components.Temperature)
}

// updates the world state
func (ws *WorldState) Update(tick *engine.Tick) {

	ws.updateSubsystems()
}

// iterates through the dependency tree for subsystems using DFS
func (ws *WorldState) updateSubsystems() {
	visited := make(map[subsystems.SubsystemID]struct{})

	depStack := utils.NewStack[subsystems.Subsystem]()

	for _, system := range ws.subsystems {
		depStack.Push(system)
		for depStack.Len() > 0 {
			subsystem := depStack.Pop()
			if _, seen := visited[subsystem.ID()]; seen {
				continue
			}

			visited[subsystem.ID()] = struct{}{}
			if len(ws.dependencies[subsystem.ID()]) <= 0 {
				subsystem.Tick(nil)
			}

			for _, dep := range ws.dependencies[subsystem.ID()] {
				if _, seen := visited[dep.id]; !seen {
					subsystem := ws.subsystems[dep.id]
					depStack.Push(subsystem)

				}
			}
		}

		inputs := make(map[components.ComponentType]*components.Component)
		for _, dep := range ws.dependencies[system.ID()] {
			inputs[dep.component.Type()] = dep.component
		}
		system.Tick(inputs)

	}

}
