# AXIOM Go Engine — MVP Todo

> **Goal:** Build a playable vertical slice of the core loop:
> *see problem → diagnose via commands → fix config → monitor recovery → something else breaks*
>
> Terminal-only. No Godot integration yet. Validate the fun before building more.
>
> **Difficulty:** 1 = straightforward, 2 = moderate, 3 = requires design decisions & iteration
>
> Items are ordered. Later items depend on earlier ones.

---

## Phase 1: Foundation

### 1. Project scaffolding
**Difficulty: 1**

- [x] `go.mod` initialized
- [x] Package structure: `simulation/`, `subsystems/`, `filesystem/`, `commands/`
- [x] `utils/` package with generic helpers (e.g., `Clamp`)
- [x] `go build` succeeds

---

### 2. World state & tick loop
**Difficulty: 2**

- [x] `WorldState` struct (tick count, named subsystem fields)
- [x] `Subsystem` interface: ID, Name, Effort, Components, Tick
- [x] `SubsystemCore` embedded struct with shared fields and accessors
- [x] `Tick()` method that updates all subsystems with explicit dependency wiring
- [x] Unit test: create a world, tick it N times, verify state changes

---

### 3. Three subsystems: power, coolant, HVAC
**Difficulty: 2**

- [x] **Power**: temperature rises from output heat, cooled by incoming coolant
- [x] **Coolant**: outputs flow rate and coolant temperature
- [x] **HVAC**: regulates ambient temperature using incoming power, affected by incoming heat
- [x] Each subsystem updates its own components during tick

---

### 4. Dependency graph & connections
**Difficulty: 3**

- [x] Connection system: ports on components, connections between subsystems with throughput multiplier
- [x] DFS-based dependency resolution in WorldState.updateSubsystems()
- [x] Power output feeds HVAC; cooling output feeds power
- [ ] Integration test: kill coolant → power overheats → HVAC degrades (the cascade)

---

## Phase 2: Connection Role Refactor & Tuning

### 5. Role-based connections (Option A)
**Difficulty: 2**

Connections currently route inputs by `ComponentType`, which conflates physical quantity with routing role. Change to role-based routing so each connection declares what role it fills on the destination subsystem.

- [ ] Add `destRole string` to `Connection` struct
- [ ] Change `Subsystem.Tick()` signature: `map[ComponentType][]Component` → `map[string][]Component`
- [ ] Update `updateSubsystems()` to key inputs by `conn.DestRole()`
- [ ] Update Power tick: read `inputs["coolant-temp"]` and `inputs["coolant-flow"]`
- [ ] Update HVAC tick: read `inputs["power-in"]` and `inputs["heat-in"]`
- [ ] Update `WorldState.Init()` to pass role names when creating connections

---

### 6. Tuning profiles
**Difficulty: 2**

Replace inline constants and ad-hoc formulas with named response profiles. Each profile describes how an input affects a component: gain (fraction of gap closed per tick), ceiling (max delta), floor (min drift).

- [ ] `ThermalResponse` struct with `Gain`, `Ceiling`, `Floor` fields and `Delta(current, target)` method
- [ ] Refactor Power tick to use profiles for coolant-temp and coolant-flow responses
- [ ] Refactor HVAC tick to use profiles for heat-in and power-in responses
- [ ] Remove `hvacHeatingRate` const, `calcHvacHeatDelta`, `calcPowerTempDelta`
- [ ] All tuning knobs visible in one place per subsystem, not scattered in formulas

---

## Phase 3: Player Interaction

### 7. Virtual filesystem overhaul
**Difficulty: 2**

The VFS exists but can't read/write content. Add the ability for file nodes to serve live data (virtual readers) or store editable content.

- [ ] Add `content string`, `writable bool` to `Node`
- [ ] `Read()` method: if `reader != nil` call it, else return `content`
- [ ] `Write(content)` method: if `writable` set content, else error
- [ ] `Cat(path)` method: resolve path and call `Read()`
- [ ] Fix `Cd()` recursion bug

---

### 8. Config parser
**Difficulty: 2**

A minimal line-oriented config format that maps directly to simulation wiring. Three directives: `system` (declare subsystem), `set` (set component value), `connect` (wire a connection with role and throughput).

```
system power    type=power
set power.effort       0.5
connect cooling.flow-out -> power coolant-flow 1.0
```

- [ ] New `engine/config/` package
- [ ] Parse `system`, `set`, `connect` directives from lines
- [ ] Return `StationConfig` struct with declarations + collected errors (line number + message)
- [ ] Unit test: parse valid config, parse config with errors

---

### 9. WorldState.ApplyConfig()
**Difficulty: 2**

Replace the hardcoded `Init()` body with config-driven setup. A factory creates subsystems by type name, then applies setpoints and wiring from the parsed config.

- [ ] Subsystem factory: `type=power` → `NewPower()`, etc.
- [ ] `nameIndex map[string]SubsystemID` for name-based lookup
- [ ] Apply `set` directives to component values
- [ ] Create ports and connections from `connect` directives
- [ ] Return errors for the player to see via `diagnose`
- [ ] Exported accessors: `Subsystems()`, `GetSubsystem(name)`
- [ ] Test: apply config → tick → verify subsystem behavior matches config

---

### 10. VFS population & live readers
**Difficulty: 2**

Wire the VFS to WorldState so the filesystem reflects live game state.

```
/station/config.ax       # writable config file
/systems/power/status    # virtual: live subsystem state
/systems/power/components # virtual: component values
/logs/system.log         # virtual: last N log lines
```

- [ ] Population function that builds the directory tree from WorldState
- [ ] Virtual readers that close over subsystem references for live data
- [ ] `/station/config.ax` initialized with starting config text (writable)
- [ ] Log ring buffer in logging package for `/logs/system.log`

---

### 11. Command engine
**Difficulty: 2**

The primary gameplay interface. Player types commands to inspect, diagnose, and manipulate the station.

- [ ] `CommandEngine` struct with `Execute(input string) string`
- [ ] Dispatch via `map[string]handler`
- [ ] `status` — table of all subsystems + component values + OK/WARN/CRIT
- [ ] `inspect <system>` — detailed view with components, connections, input values
- [ ] `diagnose <system>` — config errors, out-of-range values, hints
- [ ] `ls [path]` / `cat <path>` — delegate to VFS
- [ ] `write <path>` — multi-line input, write to VFS
- [ ] `apply` — re-parse config from VFS, call ApplyConfig(), print errors
- [ ] `set <sys>.<comp> <value>` — shortcut to modify and re-apply
- [ ] `help [cmd]` — command list and usage

---

## Phase 4: Wire It Together

### 12. Telemetry CSV export
**Difficulty: 1**

Write tick snapshots to a CSV file that Godot can read for rendering graphs. One row per component per tick. Flush after each tick so Godot can tail the file.

- [ ] `TelemetryWriter` that opens/creates CSV on init
- [ ] Header row: `tick,system,component,value`
- [ ] Append rows after each `Update()`
- [ ] Flush per tick

---

### 13. REPL + simulation goroutine
**Difficulty: 2**

Refactor `main.go` from a blocking game loop to a concurrent design: simulation ticks in a background goroutine, player interacts via REPL on the main goroutine.

- [ ] Simulation goroutine ticking once per second
- [ ] `sync.RWMutex` on WorldState (write-lock during tick, read-lock for commands)
- [ ] REPL: `bufio.Scanner` on stdin → `Execute()` → print result
- [ ] Boot message with station warning and `help` prompt

---

## Phase 5: The Playable Scenario

### 14. MVP scenario
**Difficulty: 2**

A broken starting config that creates an obvious problem the player must diagnose and fix. The fix reveals a second emergent problem. Then the player adds a new subsystem via config.

Starting config bug: HVAC `power-in` throughput is `0.0` — no power reaches HVAC, ambient temp rises. After fixing, power runs hot at effort `0.7`, creating a second problem the player solves by tuning effort or adding a second cooling unit.

- [ ] Write the broken `config.ax` starting file
- [ ] Boot message warns about rising temperature
- [ ] `diagnose hvac` hints at the zero-throughput connection
- [ ] After fix, verify HVAC temp converges toward target in telemetry
- [ ] Power overheating emerges naturally from the physics
- [ ] Adding `system cooling2 type=cooling` + connections works via `apply`

---

### 15. Tuning & playtesting
**Difficulty: 3**

Tune ThermalResponse profiles until the scenario feels right. The broken state should be obviously wrong. Recovery should be visible within ~10 ticks. Power overheating should be a gradual pressure, not instant.

- [ ] Profile values that produce satisfying convergence curves
- [ ] Verify the cascade: broken HVAC → rising temp is legible from `status`
- [ ] Verify the fix: `apply` → `status` shows improvement within a few ticks
- [ ] Verify expansion: adding a subsystem via config works without restart
- [ ] Full end-to-end walkthrough of the player flow

---

## Done Checklist

When complete, you should be able to:

- [ ] Boot the engine and see a terminal with a station warning
- [ ] Run `status` and `diagnose` to find the problem
- [ ] Read and edit the config via `cat` and `write`
- [ ] Run `apply` and watch the simulation respond
- [ ] Add a new subsystem by editing the config
- [ ] See tick-by-tick telemetry in a CSV file for Godot to render
- [ ] Feel the core loop: diagnose → fix → monitor → new problem emerges
