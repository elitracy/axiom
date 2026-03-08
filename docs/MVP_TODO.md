# AXIOM Rust Engine — MVP Todo

> **Goal:** Build the Rust engine that powers the prototype's core loop:
> *see problem → investigate → fix config/script → watch recovery → something else breaks*
>
> **Difficulty:** 1 = straightforward, 2 = moderate, 3 = requires design decisions & iteration
>
> Items are ordered. Later items depend on earlier ones.

---

## Phase 1: Foundation

### 1. Project scaffolding
**Difficulty: 1**

Set up the Rust workspace with `axiom-core` (lib) and `axiom-ffi` (cdylib). Add `serde` and `serde_json` as dependencies to `axiom-core`. Verify `cargo build` produces a shared library in `target/`. No game logic yet — just the skeleton that compiles.

- [x] Workspace `Cargo.toml` with both members
- [x] `axiom-core` lib crate with empty module declarations (`simulation`, `systems`, `filesystem`, `commands`)
- [x] `axiom-ffi` cdylib crate that depends on `axiom-core`
- [x] `serde`/`serde_json` in `axiom-core` deps
- [x] `cargo build` succeeds and outputs `.dylib`

---

### 2. World state & tick loop
**Difficulty: 2**

Define the core `WorldState` struct that holds all subsystems and a tick counter. Implement `axiom_tick()` which advances the simulation by one step — iterate over subsystems, apply degradation, update sensor values. For now subsystems can be stubs that just decrement health each tick. The important thing is the loop exists and state changes over time.

- [ ] `WorldState` struct (tick count, list of subsystems)
- [ ] `Subsystem` trait or struct: id, name, health (0-100), status enum (Online/Degraded/Critical/Offline), sensor values (`HashMap<String, f64>`), degradation rate
- [ ] Status derived from health: >70 Online, 30-70 Degraded, 1-29 Critical, 0 Offline
- [ ] `tick()` method that updates all subsystems
- [ ] Unit test: create a world, tick it N times, verify health decreases

---

### 3. Three subsystems: power, coolant, life support
**Difficulty: 2**

Implement the three concrete subsystems with their unique behaviors. Each ticks differently — power consumes fuel and produces energy, coolant manages flow rate and temperature, life support manages O2/CO2. Don't worry about dependencies between them yet, just get each one simulating independently.

- [ ] **Power**: sensors (fuel_level, output, temperature), consumes fuel per tick, output level configurable
- [ ] **Coolant**: sensors (flow_rate, pressure, pump_speed), pressure slowly leaks, pump_speed affects flow
- [ ] **Life Support**: sensors (o2_level, co2_level, scrubber_status), O2 depletes / CO2 rises each tick, scrubber reverses this when online
- [ ] Each subsystem updates its own sensors during tick
- [ ] Unit tests for each subsystem's tick behavior

---

### 4. Dependency graph
**Difficulty: 3**

Wire subsystems together so failures cascade. Power feeds life support and coolant (they degrade faster without it). Coolant keeps power temperature down (without it, power overheats and throttles/shuts down). This is the core tension engine — get it right and the game creates its own emergencies.

- [ ] Define dependencies: which subsystems feed which
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
- [ ] Parser that reads config string into `HashMap<String, f64>`
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

- [ ] `get_completions(partial: &str) -> Vec<String>`
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

## Phase 4: FFI & Integration

### 11. State serialization
**Difficulty: 1**

Implement `serde::Serialize` on all game state structs so the entire world state can be exported as JSON. This is what Godot will consume to render the dashboard and terminal.

- [ ] Derive `Serialize` on `WorldState`, `Subsystem`, sensor maps, status enums
- [ ] `get_state() -> String` returns full JSON
- [ ] `get_dashboard() -> String` returns summary JSON (just what the dashboard needs)
- [ ] Verify JSON output is clean and parseable

---

### 12. FFI exports
**Difficulty: 2**

Wire up the C FFI layer in `axiom-ffi`. Expose `extern "C"` functions that Godot can call via P/Invoke. Handle string passing across the boundary (Rust allocates, caller frees). This is the final piece — after this, a frontend can drive the engine.

- [ ] `axiom_init() -> *mut Engine` — create and return an engine instance
- [ ] `axiom_tick(engine)` — advance one tick
- [ ] `axiom_execute_command(engine, cmd) -> *const c_char` — run a command, return result
- [ ] `axiom_get_state(engine) -> *const c_char` — return world state as JSON
- [ ] `axiom_get_completions(engine, partial) -> *const c_char` — return completions as JSON
- [ ] `axiom_free_string(ptr)` — free a Rust-allocated string
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
- [ ] Get all of the above through C FFI calls returning JSON

That's your Rust MVP. Godot just needs to render what this engine computes.
