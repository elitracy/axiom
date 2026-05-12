package game

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/elias/axiom/engine/commands"
	"github.com/elias/axiom/engine/filesystem"
	"github.com/elias/axiom/engine/logging"
	"github.com/elias/axiom/engine/state"
	"github.com/elias/axiom/engine/telemetry"
	"github.com/elias/axiom/engine/utils"
)

type Game struct {
	world  *state.State
	shell  *filesystem.Shell
	engine *commands.CommandEngine
	log    *state.GameLogger
	tick   *utils.Tick
	writer *telemetry.TelemetryWriter
}

func NewGame(tick *utils.Tick) *Game {
	world := state.NewState()
	shell := filesystem.NewShell()

	shell.Populate(world)
	cmds := commands.NewCommandEngine(world, shell)
	log := state.NewGameLogger(512)

	logNode := shell.GetChild("sys/logs/station.log")

	logNode.SetReader(func() string {
		return strings.Join(log.Read(), "\n")
	})

	return &Game{world, shell, cmds, log, tick, nil}
}

func (g *Game) Init(configPath string) error {
	file, err := os.Create("telemetry.csv")
	if err != nil {
		return fmt.Errorf("Could not read file: telemetry.csv")
	}

	g.writer = telemetry.NewTelemetryWriter(file, g.tick)

	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("Could not read file: %s", configPath)
	}

	g.Cmd("write", "/usr/conf/station.ax", string(configFile))
	g.Cmd("reload")

	logging.Ok("STARTING AXIOM")

	return nil
}

func (g *Game) Cmd(cmd string, args ...string) string {
	val, err := g.engine.Execute(cmd, args...)
	if err != nil {
		logging.Error(err.Error())
		g.log.Print(err.Error())
		return fmt.Sprintf("error: %s", err.Error())
	}

	return val
}

func (g *Game) StartTickLoop() {
	ticker := time.NewTicker(time.Second)

	go func() {
		for range ticker.C {
			RunSimulation(g.world, g.tick)
			for _, system := range g.world.Subsystems() {
				g.writer.Write(system.ExportFields())
			}
		}
	}()
}

func (g *Game) Log() *state.GameLogger { return g.log }
