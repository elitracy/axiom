# AXIOM Go Engine — MVP Todo

> **Goal:** Build the Go engine that powers the prototype's core loop:
> *see problem → investigate → fix config/script → watch recovery → something else breaks*
>
> **Difficulty:** 1 = straightforward, 2 = moderate, 3 = requires design decisions & iteration
>
> Items are ordered. Later items depend on earlier ones.

---

## Phase 1: Foundation

### 1. Project scaffolding
**Difficulty: 1**

Set up the Go module with packages for simulation, systems, filesystem, and commands. Verify `go build` succeeds. No game logic yet — just the skeleton that compiles.

- [x] `go.mod` initialized
- [x] Package structure: `simulation/`, `systems/`, `filesystem/`, `commands/`
- [x] `utils/` package with generic helpers (e.g., `Clamp`)
- [x] `go build` succeeds

---

### 2. World state & tick loop
**Difficulty: 2**

Define the core `WorldState` struct that holds all subsystems and a tick counter. Implement `Tick()` which advances the simulation by one step — tick each subsystem with its dependencies, update sensor values. The world tick explicitly wires subsystem inputs/outputs (no generic loop).

- [x] `WorldState` struct (tick count, named subsystem fields)
- [x] `Subsystem` interface: ID, Name, Health, Status, Sensors, DegradationRate
- [x] `SubsystemCore` embedded struct with shared fields and accessors
- [x] Status derived from health: >70 Online, 30-70 Degraded, 1-29 Critical, 0 Offline
- [x] `Tick()` method that updates all subsystems with explicit dependency wiring
- [ ] Unit test: create a world, tick it N times, verify health decreases

---

### 3. Three subsystems: power, coolant, life support
**Difficulty: 2**

Implement the three concrete subsystems with their unique behaviors and per-system `Tick()` signatures. Each receives its dependencies as arguments from the world tick.

- [x] **Power**: `Tick(coolantFlow, ambientTemp)` — consumes fuel scaled by output level, temperature rises from output heat minus coolant dissipation, clamped to bounds
- [x] **Coolant**: `Tick(heatLoad)` — flow rate derived from pump health and backpressure, pressure accumulates from flow and bleeds passively, temperature uses diminishing returns as coolant temp approaches heat source
- [ ] **Life Support**: `Tick(powerAvailable)` — O2 consumed per tick, scrubber restores O2 and removes CO2 scaled by power available
- [x] Each subsystem updates its own sensors during tick
- [ ] Unit tests for each subsystem's tick behavior

---

### 4. Dependency graph
**Difficulty: 3**

Wire subsystems together in WorldState.Tick() so failures cascade. The wiring is explicit — power's coolant flow comes from coolant's sensor, coolant's heat load comes from power's sensor, life support's power available comes from power's sensor.

- [x] Explicit dependency wiring in WorldState.Tick() (no generic subsystem loop)
- [ ] Power output affects life support and coolant effectiveness
- [ ] Coolant failure causes power temperature to rise
- [ ] Power overheating triggers safety throttle (reduced output) then shutdown
- [ ] Life support effectiveness scales with available power
- [ ] Integration test: kill coolant → power overheats → life support degrades → O2 drops (the cascade)

---

## Phase 2: Player Interaction

### 5. Virtual filesystem
**Difficulty: 2**

Build an in-memory filesystem tree that represents game state as browsable paths. Subsystem status, sensor readings, and config files all live at paths like `/systems/power/status` and `/sensors/temp/reactor`. Support `ls` (list directory) and `cat` (read file) operations that return strings. Sensor value files should return live data from the simulation.

- [ ] VFS tree structure (directories and files)
- [ ] Populate VFS from world state each tick (or read live from state on access)
- [ ] `ls(path)` → list of entries
- [ ] `cat(path)` → file contents as string
- [ ] Paths: `/systems/{name}/status`, `/systems/{name}/config.ax`, `/sensors/{category}/{name}`
- [ ] Unit test: build VFS from world state, verify ls and cat return expected data

---

### 6. Editable config files
**Difficulty: 2**

Each subsystem has a config file (`/systems/power/config.ax`, etc.) with key-value pairs the simulation reads each tick. The player can write to these files to change behavior (e.g., adjust pump speed, fuel burn rate, power output). The simulation picks up changes on the next tick. Simple `key = value` format is fine.

- [ ] Config file format: `key = value` (one per line, `#` comments)
- [ ] Parser that reads config string into `map[string]float64`
- [ ] Each subsystem reads its config values during tick (with defaults if missing)
- [ ] `write(path, contents)` operation on the VFS
- [ ] Test: change a config value → tick → verify subsystem behavior changed

---

### 7. Command parser
**Difficulty: 2**

Accept a command string, parse it, execute it against the world state, return a result string. This is what the Godot terminal will call into. Start with the core commands from the design doc. Each command should return helpful, in-fiction output.

- [ ] Parse command string into command + arguments
- [ ] `status` — overview of all subsystems (name, health, status, key readings)
- [ ] `inspect <system>` — detailed view of one subsystem (all sensors, config values, health)
- [ ] `diagnose <component>` — identify specific faults (what's wrong and where to look)
- [ ] `help` / `help <command>` — command list and usage
- [ ] `ls <path>` / `cat <path>` — filesystem commands (delegate to VFS)
- [ ] `set <component> <param> <value>` — modify a config value
- [ ] `restart <component>` — attempt to restart an offline subsystem
- [ ] Unknown command → helpful error message
- [ ] Unit tests for each command

---

### 8. Tab completion
**Difficulty: 1**

Given a partial input string, return a list of valid completions. Source completions from the command list, filesystem paths, subsystem names, and component IDs. Doesn't need to be fancy — prefix matching is fine.

- [ ] `GetCompletions(partial string) []string`
- [ ] Complete command names when input has no space
- [ ] Complete paths/subsystem names after a command
- [ ] Test: partial input → expected completions

---

## Phase 3: The Fun Part

### 9. The fixable script (first puzzle)
**Difficulty: 2**

The life support scrubber is offline because its startup script (`/systems/life-support/scrubber.ax`) has a bug — a bad sensor reference. The player must find the file, read it, spot the error, fix it via the VFS write operation, and the scrubber comes online on the next tick. This is the game's first "aha" moment.

- [ ] Scrubber startup script exists in VFS with an intentional bug (e.g., references `sensor.o2-main` instead of `sensors.o2-main`)
- [ ] Simple script validator that checks sensor references against known valid paths
- [ ] On tick: if script is valid → scrubber online, if invalid → scrubber stays offline with error in diagnose output
- [ ] `diagnose atmo-scrubber` hints at the problem ("startup script error on line 3: unknown reference")
- [ ] Test: fix the script → tick → scrubber comes online → CO2 starts dropping

---

### 10. Entropy engine
**Difficulty: 2**

Random events that keep the station from ever being "solved." Components degrade at varying rates. Occasionally a config value gets corrupted, or a subsystem faults. Frequency should be tunable (a knob you can turn up for testing or down for a chill session). This is what makes the core loop repeat.

- [ ] Entropy source with configurable frequency (events per N ticks)
- [ ] Event types: accelerated wear (health drop), config corruption (value changed), component fault (subsystem forced to degraded/critical)
- [ ] Events logged so the player can discover what happened (`/logs/events.log` in VFS)
- [ ] Test: run simulation with high entropy → verify events fire and state changes

---

## Phase 4: Serialization & Integration

### 11. State serialization
**Difficulty: 1**

Add JSON tags to all game state structs so the entire world state can be exported as JSON via `encoding/json`. This is what Godot will consume to render the dashboard and terminal.

- [ ] JSON tags on `WorldState`, subsystem structs, sensor maps, status enum
- [ ] `GetState() string` returns full JSON
- [ ] `GetDashboard() string` returns summary JSON (just what the dashboard needs)
- [ ] Verify JSON output is clean and parseable

---

### 12. CGO exports (deferred)
**Difficulty: 2**

Wire up CGO `//export` functions so the engine can be built as a C shared library (`-buildmode=c-shared`) for Godot. Not a priority until the engine is playable standalone.

- [ ] `axiom_init()` — create and return an engine instance
- [ ] `axiom_tick(engine)` — advance one tick
- [ ] `axiom_execute_command(engine, cmd)` — run a command, return result as C string
- [ ] `axiom_get_state(engine)` — return world state as JSON
- [ ] `axiom_get_completions(engine, partial)` — return completions as JSON
- [ ] `axiom_free_string(ptr)` — free a Go-allocated C string
- [ ] `axiom_save_file(engine, path, contents)` — write to VFS
- [ ] Smoke test: init → tick 10 times → get_state → verify JSON parses

---

## Done Checklist

When all 12 items are complete, you should be able to:

- [ ] Init an engine, tick it, and watch subsystems degrade over time
- [ ] Send commands and get back meaningful text responses
- [ ] Browse the virtual filesystem and read live sensor data
- [ ] Edit a config file and see the simulation respond next tick
- [ ] Fix the scrubber script and watch life support recover
- [ ] Watch cascading failures when you neglect a system
- [ ] See random entropy events create new problems
- [ ] Get all of the above through CGO calls returning JSON

That's your Go MVP. Godot just needs to render what this engine computes.
