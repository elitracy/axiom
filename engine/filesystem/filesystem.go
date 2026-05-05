package filesystem

import (
	"fmt"
	"strings"

	"github.com/elias/axiom/engine/state"
	"github.com/elias/axiom/engine/subsystems/connections"
	"github.com/elias/axiom/engine/utils"
)

type Shell struct {
	cwd  *Node
	root *Node
}

type worldState interface {
	Subsystems() []state.Subsystem
	Connections() map[utils.SubsystemName]map[utils.PortType][]*connections.Connection
}

func NewShell() *Shell {
	return &Shell{}
}

func (s *Shell) Populate(ws worldState) {
	if ws == nil {
		return
	}

	root := NewDir("/")
	sys := NewDir("sys")
	usr := NewDir("usr")
	root.AddChild(sys)
	root.AddChild(usr)

	conf := NewDir("conf")
	bin := NewDir("bin")
	usr.AddChild(conf)
	conf.AddChild(NewFile("station.ax"))
	usr.AddChild(bin)

	logs := NewDir("logs")
	systems := NewDir("systems")
	sys.AddChild(systems)
	sys.AddChild(logs)

	stationLog := NewFile("station.log")
	logs.AddChild(stationLog)

	power := NewDir("power")
	cooling := NewDir("cooling")
	machines := NewDir("machines")

	systems.AddChild(power)
	systems.AddChild(cooling)
	systems.AddChild(machines)

	s.root = root
	s.cwd = root

}

func (s *Shell) ReloadSubsystems(ws worldState) {
	power := s.GetChild("sys/systems/power")
	cooling := s.GetChild("sys/systems/cooling")
	machines := s.GetChild("sys/systems/machines")

	for _, subsystem := range ws.Subsystems() {
		dir := NewDir(string(subsystem.Name()))
		status := NewFile("status")
		components := NewDir("components")
		ports := NewDir("ports")
		dir.AddChild(status)
		dir.AddChild(components)
		dir.AddChild(ports)

		status.SetReader(func() string {
			return subsystem.Status().String()
		})

		for _, component := range subsystem.Components() {
			dir := NewDir(component.Name())
			components.AddChild(dir)

			compType := NewFile("type")
			compType.SetReader(func() string {
				return component.Type().String()
			})

			compValue := NewFile("value")
			compValue.SetReader(func() string {
				return fmt.Sprintf("%.2f", component.Value())
			})

			dir.AddChild(compType)
			dir.AddChild(compValue)

			ports := NewDir("ports")
			dir.AddChild(ports)

			for _, port := range subsystem.InputPorts() {
				if port.Component().Name() == component.Name() {
					portDir := NewDir(port.Name())
					connectionFile := NewFile("connection")
					portDir.AddChild(connectionFile)

					connectionFile.SetReader(func() string {
						conns := ws.Connections()[subsystem.Name()]

						for _, conn := range conns[utils.PortInput] {
							if conn.DestPort().ID() == port.ID() {
								return fmt.Sprintf("← %s.%s @ %.2f", conn.SrcSystem(), conn.SrcPort().Name(), conn.Throughput())
							}
						}

						return "<no-connection>"
					})

					ports.AddChild(portDir)
				}
			}

			for _, port := range subsystem.OutputPorts() {
				if port.Component().Name() == component.Name() {
					portDir := NewDir(port.Name())
					connectionFile := NewFile("connection")
					portDir.AddChild(connectionFile)

					connectionFile.SetReader(func() string {
						conns := ws.Connections()[subsystem.Name()]

						for _, conn := range conns[utils.PortOutput] {
							if conn.SrcPort().ID() == port.ID() {
								return fmt.Sprintf("→ %s.%s @ %.2f", conn.DestSystem(), conn.DestPort().Name(), conn.Throughput())
							}
						}

						return "<no-connection>"
					})

					ports.AddChild(portDir)
				}
			}

		}

		switch subsystem.Type() {
		case utils.Power:
			power.AddChild(dir)
		case utils.Cooling:
			cooling.AddChild(dir)
		case utils.Hvac:
			machines.AddChild(dir)
		}
	}
}

func (s *Shell) Ls(path string) string {
	if path == "" {
		return s.cwd.ls(path)
	}

	if path[0] == '/' {
		s.root.ls(path)
	}

	return s.cwd.ls(path)
}

func (s *Shell) Cd(path string) {

	if path == "." {
		return
	}

	if path == ".." && s.cwd.Parent() != nil {
		s.cwd = s.cwd.Parent()
	}

	node := s.cwd.GetChild(path)

	if node != nil {
		s.cwd = node
	}

}

func (s Shell) Cat(path string) string {
	path = strings.Trim(path, "/")
	node := s.cwd.GetChild(path)

	if node == nil {
		return ""
	}

	return node.Read()
}

func (s Shell) Pwd() string {
	return s.cwd.pwd()
}

func (s Shell) Tree(path string, depth int) string {
	node := s.cwd.GetChild(path)

	if node == nil {
		return ""
	}

	return node.tree("", true, depth)
}

func (s Shell) GetChild(path string) *Node {
	path = strings.Trim(path, "/")
	node := s.cwd.GetChild(path)

	if node == nil {
		return nil
	}

	return node
}

func (s Shell) Find(path string) *Node {
	path = strings.Trim(path, "/")

	node := s.cwd.FindChild(path)

	if node == nil {
		return nil
	}

	return node
}
