# AXIOM

## A Terminal-Based Station Management & Programming Game

**Genre:** Programming / Simulation / Survival / MMO
**Platform:** Terminal (TUI) → Godot dashboard
**Engine:** Go (simulation) + Godot (frontend)
**Multiplayer:** Single-player → Persistent MMO transition

---

## 1. Vision

You wake up alone in a failing underground installation. The world above is silent. The systems keeping you alive are degraded, undocumented, and held together by the work of engineers long gone. Your only tools are a terminal and AXIOM OS — the station's retro operating system.

Through engineering, scripting, and problem-solving, you repair what's broken, expand what's possible, and eventually reach beyond the walls to discover a network of other survivors doing the same thing. What begins as a solitary fight for survival becomes a massively multiplayer civilization built entirely on code.

AXIOM is a game where **programming is the gameplay**. You monitor systems through a TUI dashboard, diagnose failures by reading logs and tracing signal paths, and fix problems by editing config files and writing scripts in a clean domain-specific language. As your skills grow, so does your station — and eventually, so does your reach.

The game is split into two interconnected parts:

- **Part 1: The Bunker** — A complete single-player experience. Rebuild your station, learn AXIOM, survive.
- **Part 2: The Network** — The same game, expanded into a persistent MMO. Discover other players, build shared infrastructure, form alliances, wage code wars, and engineer a civilization.

---

## 2. Core Design Pillars

### 2.1 Programming IS the Gameplay

No combat system. No crafting grid. The game is engineering. You fix systems by editing real text files, writing real scripts, and reasoning about real (simulated) system behavior. Mastery is rewarded with speed, efficiency, and automation.

### 2.2 Discoverable Complexity

Inspired by modded Minecraft (Feed The Beast, Create). Information is always accessible — never hidden behind external wikis or manuals. Tab completion shows every command. The filesystem is the documentation. `help`, `list`, and `info` are always one keystroke away. The complexity is opt-in: simple problems have simple solutions, and depth reveals itself as you explore.

### 2.3 Earned Growth

Your station grows because you made it grow. Every green panel on your dashboard represents a system you fixed, a script you wrote, or an automation you deployed. The dashboard *is* your portfolio. Growth isn't gated by XP or levels — it's gated by your engineering ability.

### 2.4 Living World

The simulation runs whether you're watching or not. Systems degrade. Events unfold. In multiplayer, other players act. Your automation scripts are your proxy — the quality of your code determines what you come back to.

### 2.5 Infrastructure as Power (Part 2)

In multiplayer, there are no artificial power systems, currencies, or governance mechanics. Power emerges from what you build. Control a critical relay node? You have leverage. Build the power grid that feeds a sector? You can shape policy. Engineering is geopolitics.

---

## 3. The World

### 3.1 Setting

An unspecified future. Something happened to the surface — the details are fragmented, discovered through old logs, corrupted archives, and eventually, direct observation. The underground is a vast network of automated installations built to sustain life independently. They've been running on autopilot for a long time. Too long.

The technology is advanced but aged. Systems are a mix of sophisticated engineering and desperate jury-rigging by previous occupants. Comments in old scripts tell stories. Error logs hint at what went wrong. The lore is embedded in the infrastructure — you discover the world by repairing it.

### 3.2 Geography

The underground network spans a large region. Installations vary in size, purpose, and condition. Some are small relay stations. Others are massive industrial complexes. They're connected by physical tunnel infrastructure (cable runs, pipes, rail lines) and, eventually, wireless links and a dormant satellite constellation overhead.

Geography matters for multiplayer — physical proximity between stations allows hardline connections (fast, reliable). Distance requires satellite relay (higher latency, more fragile, shared resource).

### 3.3 The Satellite Constellation

A network of satellites orbits above, placed there before the collapse. They're dormant but functional. Bringing them online requires significant engineering: power, antenna arrays, signal processing, orbital tracking. Each satellite a player activates extends the network's reach. The constellation is a shared resource the community must collectively maintain — if satellites go unmaintained, regions lose connectivity.

---

## 4. The Operator Fantasy (Interaction Model)

### 4.1 You Are at the Control Terminal

The core fiction: you are sitting at the central operations console of an automated facility. This is the nerve center the station was designed to be operated from. The original engineers sat in this same chair. Everything in the facility was built to be remotely operated from this terminal — this is real industrial control system (SCADA) design. Modern power plants, water treatment facilities, and factories work exactly this way.

### 4.2 The Control Bus

Every valve has a motor on it. Every pump has a controller. Every circuit breaker is electronically switchable. Your terminal connects to subsystem controllers through a wired control network running through the station. When you type `open valve.coolant-main`, you're sending a command down the control bus to a motor that physically turns the valve.

The control bus IS the game's internal network. When that network is damaged, you lose the ability to command things remotely. A severed cable to sublevel 3 means your commands have nowhere to go. This creates a fundamental failure mode: hardware is fine, but the control link is broken. The reactor is physically intact but you can't start it because the command bus to its controller is down.

This also defines discovery — you can't see or interact with any section of the station until you have a data link to it. Establishing control network connections to new areas is how the game world expands.

### 4.3 Four Layers of Interaction

**Software problems — fix directly from the terminal.**
Corrupted scripts, bad configs, logic errors, miscalibrated values. Edit files, run diagnostics, write automation. Instant feedback. This is the most common and most frequent type of work.

**Remotely operable hardware — command from the terminal.**
Valves, pumps, switches, breakers, motors. Send a command, the actuator responds. Fast, but depends on the control bus being intact. If a motor is burnt out, your command returns an error — now it's a physical problem.

**Physical repairs — dispatch drones.**
Replacing a component, patching a pipe, splicing a cable, clearing debris. Takes time (ticks), costs resources (drone battery, spare parts). You monitor progress from your terminal and may need to adjust the drone's approach if it hits obstacles.

**Major construction — fabrication plus drone deployment.**
Building new infrastructure, expanding sections, installing new subsystems. The biggest time and resource investment. Queue parts in the fabrication bay, wait for them to build, then dispatch drones to install.

### 4.4 A Typical Repair Sequence

A complete example of how these layers interact in a single problem:

1. Dashboard shows temperature anomaly in coolant system (monitoring)
2. `diagnose coolant` traces it to low flow in coolant loop 3 (investigation)
3. `inspect pump.coolant-3` — controller is running but output is low (software check)
4. `cat /systems/coolant/pump-3.ax` — the speed ramp curve has a bad value, probably from data corruption. Fix the value, save. Pump speeds up. (software fix)
5. Flow improves but is still below nominal. `inspect pipe.coolant-7` shows 60% corrosion (physical problem identified)
6. `fabricate pipe-segment model:coolant-7` — queue a replacement part (resource cost, wait time)
7. `dispatch drone-2 repair pipe.coolant-7` — send drone with the new part once fabricated (physical repair, more wait time)
8. Monitor drone progress from terminal. Pipe replaced. Flow returns to nominal. (resolution)
9. Write a monitoring script to alert if flow drops again (automation / prevention)

Every step makes sense within the fiction. The player is never wondering *how* their commands affect the world.

---

## 5. Learning Design (How the Game Teaches)

### 5.1 The Station Is the Tutorial

There is no tutorial screen, no popup tips, no separate training mode. The station's failing systems are the tutorial. The game teaches by giving you an urgent problem and making the next step always visible.

**Boot sequence:**

```
[AXIOM OS v2.7.1 — EMERGENCY BOOT]

WARNING: Multiple system failures detected.
WARNING: Atmospheric O2 at 18.2% — declining.
WARNING: Primary power unstable — backup generator active.

Type 'status' to assess station condition.
```

One instruction. One clear action. The player types `status`, sees a dashboard with critical alerts, and the most urgent one says: `SCRUBBER UNIT OFFLINE — CO2 rising`. Below it: `Run 'inspect life-support' for details`.

### 5.2 Every Problem Teaches One New Skill

The first problem (fix the scrubber) teaches the core diagnostic loop through necessity:

- `status` → see what's wrong (learn the dashboard)
- `inspect life-support` → see which component is broken (learn inspection)
- `diagnose atmo-scrubber` → see the specific fault and its location (learn diagnostics)
- `cat /systems/life-support/atmo-scrubber/startup.ax` → read the broken script (learn the filesystem)
- Edit the file, fix the sensor reference typo → (learn the editor)
- Scrubber starts, CO2 drops, panel goes from red to yellow → (learn that fixes have immediate feedback)

The player learned six core skills and it felt like solving an emergency, not reading a manual.

### 5.3 Problem Complexity Progression

**First hour — heavily guided.**
Fix a typo. Change a value. Uncomment a line. The diagnostic system tells you exactly where to look and suggests what to do. Learning basic commands and filesystem layout.

**Next few hours — lighter guidance.**
Diagnostics tell you "fault in coolant controller" but not which line. You read the code and find the issue yourself. Problems have two to three steps. Learning to investigate.

**Mid game — symptoms, not causes.**
"Coolant pressure dropping." Is it a pump failure? A leak? A bad config? A sensor giving wrong readings? You trace through the system. Maybe everything looks fine until you realize the pressure set point was silently corrupted. Learning to think systemically.

**Late game — minimal hand-holding.**
The system tells you something is wrong. You know the tools. You know the patterns. You diagnose from experience and muscle memory. The reward is speed and confidence — you *know* this station.

### 5.4 Discovery Through Breadcrumbs

Players never randomly stumble on new areas. The station's existing systems create natural leads:

- The backup generator's diagnostic says it's a temporary unit — there's a primary power source elsewhere
- Archived facility maps reference rooms and sections you haven't accessed
- Network topology shows dead zones where the control bus doesn't reach
- Sensors near sealed areas occasionally pick up anomalous readings (heat, vibration, atmosphere)
- Old logs reference systems by name that don't exist in your current filesystem

Each breadcrumb gives you a goal ("find the reactor on sublevel 3") and achieving that goal requires a chain of engineering steps (establish network link → power the corridor → repair the access route → diagnose the new hardware → bring it online). Each step is clear and achievable individually.

### 5.5 Design Rules for Keeping It Fun

**Always one clear next action.** The player should never stare at the screen without a lead. Every diagnostic suggests a follow-up. Every alert points somewhere. Every dead end has a "try this instead."

**Fast feedback loops.** When you fix something, the result is visible immediately. Dashboard updates. Alert clears. Values change. Never make the player wonder if their fix worked.

**Dangerous problems are loud.** Atmosphere failing, power going out — these are blinking, screaming alerts. The player is never blindsided by a silent catastrophe. Subtle problems exist (slow efficiency decline, gradual resource drain) but they're never immediately lethal. You have time to notice.

**Mistakes are recoverable.** Bad config? Error message, old config stays active. Botched reactor startup? Safety interlocks engage. Drone stuck? Recall it. The penalty for experimenting is time and resources, never game over. The game should encourage tinkering, not punish it.

**Help is always one command away.** `help` shows command categories. `help diagnose` explains with examples. `info reactor-01` gives plain-language description. `list sensors` shows every accessible sensor. `list commands` shows everything you can type. This is the FTB/JEI equivalent — the information is never hidden.

**The filesystem is the documentation.** Sensor names are discoverable by browsing (`ls /sensors/`). Component IDs are discoverable by inspecting systems. Command targets are tab-completable. The game world *is* the reference manual.

---

## 6. AXIOM OS

AXIOM OS is the in-game operating system. It is the player's entire interface to the game world. It boots with scan lines and a blinking cursor. Everything happens through AXIOM.

### 6.1 Shell

The AXIOM shell is a retro terminal interface with the following core interactions:

- **Command input** with tab completion on everything — commands, paths, sensor names, component IDs
- **Filesystem navigation** — the game state is the filesystem. `ls /systems/` shows subsystems. `cat /sensors/temp/reactor` shows a live value. `cd /scripts/` and you're in your automation directory.
- **Built-in help** — `help` for commands, `info <component>` for detailed specs, `list <category>` for discovery. Never guessing.
- **Error messages that teach** — contextual, helpful, suggesting corrections. "Line 3: `valve.coolant-maim` — did you mean `valve.coolant-main`?"

### 6.2 Dashboard

The dashboard is a multi-pane TUI display showing live system state. It is the player's primary monitoring interface.

Key properties:
- **Dynamic** — panels are added as you bring systems online. Early game: sparse, a few flickering readouts. Late game: mission control.
- **Customizable layout** — arrange panels to your preference. Your station, your ops center.
- **Drill-down** — select any panel to inspect the underlying subsystem in detail: component health, sensor readings, active scripts, recent events.
- **Color-coded status** — green (nominal), yellow (degraded), red (critical), dark (offline). The goal is a sea of green. The reality is managing the red.
- **Network view (Part 2)** — a top-level view showing all connected stations, link health, and shared infrastructure status.

Running `axiom dashboard` or `axiom status` shows the current state at a glance.

### 6.3 In-Game Editor

A built-in text editor with syntax highlighting for AXIOM Script. Minimal but functional — the goal is fast iteration, not IDE features. Open a file, edit, save, and the system hot-reloads.

Navigation is vim-inspired (given the target audience), with a command palette for AXIOM-specific actions: insert a component reference, look up a sensor path, validate syntax.

### 6.4 The Filesystem as Game State

This is a core design principle. The filesystem *is* the game.

```
/
├── systems/           # Subsystem definitions and state
│   ├── power/
│   │   ├── grid.ax          # Power distribution config
│   │   ├── reactor.ax       # Reactor controller
│   │   └── status           # Live status (read-only)
│   ├── coolant/
│   ├── comms/
│   ├── life-support/
│   └── sensors/
├── scripts/           # Player-written automation
│   ├── watchdogs/
│   ├── macros/
│   └── tests/
├── sensors/           # Live sensor data (read-only)
│   ├── temp/
│   ├── pressure/
│   ├── power/
│   └── atmo/
├── logs/              # System and event logs
│   ├── system.log
│   ├── events.log
│   └── errors.log
├── archives/          # Lore: old logs, notes, corrupted files
├── network/           # (Part 2) Shared multiplayer namespace
│   ├── stations/
│   ├── shared-systems/
│   └── routes/
└── docs/              # In-game documentation, auto-generated
```

Players interact with the game by navigating this tree, reading files, editing configs, and writing scripts. The filesystem is always the source of truth.

---

## 7. AXIOM Script

### 7.1 Design Philosophy

AXIOM Script is the domain-specific language players use to configure, automate, and control their station. It sits at a careful intersection:

- **Feels like real programming** — text files edited in a real editor, real syntax, rewards muscle memory and speed
- **Plays like a puzzle game** — constrained vocabulary, discoverable through tab completion, contextual errors guide you
- **Scales with the player** — simple configs early, complex logic late, same language throughout

It is NOT a visual/block language. It is NOT a general-purpose programming language. It is a clean, readable, domain-specific language designed for one thing: managing systems.

### 7.2 Language Overview

**Configuration mode** — declaring how systems should behave:

```
# /systems/power/grid.ax
# Power distribution manifest

source reactor.main {
  output 2400
  priority critical
}

distribute {
  life-support  -> 600  [priority: critical]
  comms-array   -> 400  [priority: high]
  sensor-grid   -> 300  [priority: high]
  fabrication   -> 500  [priority: normal]
  drone-bay     -> 400  [priority: normal]
  reserve       -> 200  [priority: low]
}

fallback reactor.backup {
  trigger source.reactor.main offline
  output 1200
  shed priority below normal
}
```

**Trigger mode** — reactive automation:

```
# /scripts/watchdogs/coolant-monitor.ax
# Coolant system watchdog

trigger temp.reactor > 85 {
  open valve.coolant-main
  set pump.coolant-2 speed 80
  wait until temp.reactor < 70
  set pump.coolant-2 speed 40
  log "Reactor temp stabilized"
}

trigger pressure.coolant < 30 {
  alert "Coolant pressure critical"
  close valve.coolant-main
  run diagnostics pump.coolant-2
}
```

**Test mode** — validation scripts:

```
# /scripts/tests/coolant-failover.ax
# Verify coolant failover procedure

test "coolant pump failure triggers backup" {
  simulate pump.coolant-1 offline
  wait 10 ticks
  assert pump.coolant-2 speed > 0
  assert temp.reactor < 90
  log "PASS: backup pump activated"
}

test "pressure loss triggers alert" {
  simulate pressure.coolant drop-to 20
  wait 5 ticks
  assert alerts contains "Coolant pressure critical"
  log "PASS: alert triggered"
}
```

**Macro mode** — recorded command sequences:

```
# /scripts/macros/emergency-power.ax
# Emergency power shedding procedure

macro emergency-power {
  set reactor.main output max
  shed priority below normal
  alert "Emergency power mode — non-critical systems offline"
  wait until power.reserve > 80
  restore priority normal
  alert "Power reserves restored — resuming normal operation"
}
```

### 7.3 Language Features

- **Namespaced references** — `temp.reactor`, `pump.coolant-2`, `valve.coolant-main`. Tab-completable, browsable via the filesystem.
- **Triggers** — conditional blocks that fire when sensor values cross thresholds or system events occur.
- **Control flow** — `if/else`, `wait until`, `wait <n> ticks`, `repeat`, `while`.
- **Commands** — `open`, `close`, `set`, `run`, `alert`, `log`, `simulate`, `assert`.
- **Composition** — scripts can call other scripts. Macros can invoke macros. Build layered automation.
- **Hot reload** — save a file and the system picks up changes immediately. No compile step, no restart.
- **Sandboxed execution** — scripts have resource limits. Infinite loops get killed with a clear error. No crashing the station (unless that's the bug you're debugging).

### 7.4 Network Extensions (Part 2)

When multiplayer unlocks, AXIOM Script gains new capabilities:

```
# Remote operations
query station.4417 power.status
send station.4417 resource.fuel amount 200

# Network triggers
trigger network.link.station-4417 latency > 500 {
  alert "Link to 4417 degrading"
  reroute traffic via station.2201
}

# Shared system definitions
shared service power-exchange {
  export power.surplus to network.grid
  import power.deficit from network.grid
  balance every 30 ticks
}
```

---

## 8. Simulation Engine

### 8.1 Tick-Based World Model

The simulation advances in discrete ticks. Each tick:

1. All subsystem states update (health degrades, resources deplete, temperatures shift)
2. Player scripts and triggers evaluate
3. Events check preconditions and fire
4. Dependencies cascade (failures propagate through the graph)
5. Dashboard and sensor values refresh

Ticks run in real-time (approximately 1 per second) during gameplay. When a player is offline (Part 2), ticks continue on the server at the same rate — your automation keeps running.

### 8.2 Subsystems

Each subsystem is a simulation entity with:

- **Health** — 0-100, degrades over time and through damage. Below thresholds, behavior changes (reduced output, intermittent failures, full shutdown).
- **Components** — individual parts that can fail independently. A pump, a valve, a circuit board. Each has its own health and failure modes.
- **Dependencies** — what this system needs from other systems to function (power, coolant, data inputs).
- **Outputs** — what this system provides to others (power, cooled fluid, sensor data, processed signals).
- **Configuration** — the `.ax` files that define its behavior. Editable by the player.
- **Sensors** — data points exposed to the player and to scripts. Temperature, pressure, throughput, error rates.

### 8.3 Core Subsystems

**Reactor / Power Plant**
The foundation. Generates power, consumes fuel, produces heat. Everything depends on it. When it struggles, everything struggles. Has a backup generator with limited output for emergencies.

**Coolant System**
Manages heat across the station. Pumps, valves, pipes, heat exchangers. A coolant failure leads to thermal cascades — components overheat, degrade faster, fail. Pipe leaks are common and require diagnosis (which section is losing pressure?).

**Life Support**
Atmospheric scrubbers, O2 generation, CO2 removal, water recycling, temperature regulation. Gradual degradation is scarier than sudden failure — you might not notice air quality dropping until it's a crisis. Monitoring atmospheric sensor trends is key.

**Communications Array**
Internal comms between station sections, and eventually external comms to the surface and satellites. Needs power and functioning hardware. When internal comms degrade, your visibility into remote sections decreases. When external comms go down, you're isolated.

**Sensor Grid**
Perimeter sensors, internal diagnostics, environmental monitoring. When sensors fail, you lose information — your dashboard goes dark in sections, diagnostics become unreliable, threats go undetected. Arguably the most important system to maintain because it's how you know about everything else.

**Fabrication Bay**
Builds replacement parts, new components, and expansion modules. Consumes raw materials and power. Queue-based — you schedule jobs and they take time. Having spare parts on hand vs. fabricating on demand is a strategic choice.

**Drone Hangar**
Automated repair and survey drones. Need firmware (scripts you write), power, and maintenance. Can be dispatched to repair external systems, survey new areas, or (Part 2) interact with the surface.

**Storage & Inventory**
Finite resources. Fuel, raw materials, spare parts, consumables. Everything you use to repair and expand comes from here. Resupply is a challenge — initially scavenging from sealed sections, later trading with other players or running surface expeditions.

**Archive / Data Core**
Corrupted databases containing lore, old blueprints, previous engineers' notes and scripts. Restoring and decrypting archives reveals game narrative, station history, and sometimes useful technical information (a blueprint for a more efficient reactor config, an old script that solves a problem you're facing).

**Satellite Uplink (Late Game)**
Massive antenna arrays, signal processing systems, orbital tracking computers. Bringing satellites online is a multi-stage engineering challenge: physical hardware repair, software configuration, signal calibration. Each satellite extends your communications range and unlocks new network capabilities.

### 8.4 Dependency Graph & Cascades

Systems are connected through a dependency graph. Failures propagate:

```
Coolant leak
  → Reactor overheats
    → Safety throttle reduces power output
      → Sensor grid browns out (low priority)
        → Lost visibility into sections 3-7
          → Missed hull breach in section 5
            → Atmosphere venting
              → Life support stressed
                → O2 dropping in adjacent sections
```

The cascade is the core tension generator. Small problems become big problems if ignored. The player's job is to intervene before cascades spiral — or to write automation robust enough to handle the early stages automatically.

### 8.5 Procedural Problem Generation

Problems are generated by the simulation, not pre-scripted. The engine introduces entropy:

- **Wear** — components degrade over time at varying rates. Older, more-used components fail more often.
- **Environmental stress** — external events (storms, seismic activity, temperature extremes) stress specific systems.
- **Corruption** — script files and configs can become corrupted. A threshold value changes. A routing entry gets garbled. The player must find and fix the corruption.
- **Emergent bugs** — the simulation can introduce subtle logic errors into system controller scripts. An off-by-one error, a race condition, a miscalibrated sensor offset. The player debugs real (simulated) code.
- **Resource depletion** — consumables run out. Fuel, filters, spare parts. Requires forward planning and resource management.

Because problems emerge from simulation state, every playthrough is different. The same system might fail in different ways depending on what else is happening, what the player has built, and what they've neglected.

---

## 9. Gameplay Progression

### 9.1 Phase 1 — Awakening (Early Game)

The station boots. Systems are failing. The player learns AXIOM basics through immediate necessity.

**State:** 1 sector accessible, 3-4 subsystems online (barely). Dashboard is sparse — a few panels, mostly yellow and red. Power is unstable, air is thin, sensors are spotty.

**Gameplay:** Learn to navigate the filesystem. Read logs to understand what's broken. Make first edits to config files (change a threshold, uncomment a line, swap a sensor reference). Run `help` constantly. Tab-complete everything. Each fix is a small victory — a panel goes from red to yellow.

**Narrative:** You find logs from the previous engineer. They were here alone too. Their notes are practical — how to restart the coolant pump, where the spare fuses are. But they stopped writing. The last log is dated months ago.

**Skills learned:** Filesystem navigation, reading logs, basic config editing, understanding the dashboard.

### 9.2 Phase 2 — Stabilization (Early-Mid Game)

Core systems are functional but fragile. The player begins writing their first scripts.

**State:** Core sector stable, adjacent sectors accessible but unpowered. Dashboard is functional with 5-8 panels. Some green, some yellow. Occasional red alerts.

**Gameplay:** Write first trigger scripts — simple watchdogs that monitor a value and take action. Deploy first macro for a common repair sequence. Start to feel the cascade mechanic as multiple things fail at once and you triage. Begin exploring sealed sections — each one requires power hookup and system activation.

**Narrative:** Deeper into the archives. The installation is bigger than you thought. References to other installations, a surface observation post, a "network" that went dark. Blueprints for systems you haven't found yet.

**Skills learned:** AXIOM Script basics (triggers, commands), writing and deploying scripts, multi-system triage, resource management.

### 9.3 Phase 3 — Expansion (Mid Game)

The player actively grows their station by restoring dormant sections and bringing new subsystems online.

**State:** Multiple sectors online. Fabrication bay operational. Drone fleet coming online. Dashboard has 12-20+ panels across multiple views. Resource supply chains established.

**Gameplay:** Each new section is a discovery — new systems, new challenges, sometimes new types of problems you haven't seen before. The fabrication bay lets you build replacement parts instead of scavenging. Drones extend your reach. Writing more sophisticated automation — scripts that call scripts, conditional logic, scheduled tasks. Start writing test scripts to validate your automation.

**Narrative:** You find the surface access corridor. Sealed, but the cameras might work. You restore a camera feed and see daylight for the first time. The surface is quiet, overgrown, but intact. Weather data starts flowing in. References to the satellite uplink facility become more frequent in the archives.

**Skills learned:** Complex scripting, test-driven automation, system architecture (designing how new systems integrate), resource optimization, drone fleet management.

### 9.4 Phase 4 — The Uplink (Late Game / Transition)

The player discovers and restores the satellite communications facility.

**State:** Station is substantial and largely automated. Dashboard is mission control. The player is now tackling the biggest engineering challenge yet — the satellite uplink.

**Gameplay:** The uplink facility is a major undertaking. Massive power requirements (may need to redesign the grid). Complex signal processing code. Antenna hardware that needs physical repair (drone missions). Orbital mechanics calculations to track and contact dormant satellites. Each satellite brought online extends range and capability. The first successful ping to a satellite and the response back is a landmark moment.

**Narrative:** The first satellite you contact returns telemetry. From orbit, you can see the region. Other installations are visible — some dark, some showing faint heat signatures. Then, on a frequency you weren't monitoring, a signal. It's structured. It's AXIOM protocol. Someone else is out there.

**Skills learned:** Large-scale system integration, advanced scripting (signal processing, orbital calculations), infrastructure architecture, performance optimization.

### 9.5 Phase 5 — Contact (Part 2 Begins)

The game transitions from single-player to multiplayer. This is not a hard boundary — the player's station, scripts, and progress carry over completely.

**The moment:** Your comms terminal prints a handshake from another station. A real player. You can see their station ID, their uptime, their signal strength. You can send them a message. You can, if you both agree, establish a network link.

From here, the game opens up. See Part 2 design below.

---

## 10. Part 2: The Network (Multiplayer)

### 10.1 The Transition

When a player brings their satellite uplink online, they join the persistent multiplayer world. Their station becomes a node in a growing network of player-operated installations. Everything from Part 1 persists — station state, scripts, automation, resources. The game simply expands.

The `/network/` namespace appears in the filesystem. New AXIOM Script keywords unlock. The dashboard gains a network overview panel. The world gets bigger.

Players who haven't reached the uplink yet are still playing Part 1 — they exist in the world but can't be contacted. As more players come online, the network grows organically.

### 10.2 Persistence & Offline Automation

Stations run 24/7 on the server. When a player logs off, their automation scripts continue executing. This is the critical design element that makes the MMO work:

- **Good automation = resilience.** If your scripts handle common failures, you come back to a healthy station.
- **Poor automation = decay.** If your scripts can't handle what happens overnight, you come back to cascading failures and a full log.
- **No automation = disaster.** An unattended station with no scripts will degrade quickly. This is the incentive to get good at AXIOM.

Players receive async notifications (if they opt in) for critical events — "Your reactor went offline at 3:42 AM. Backup generator active. Estimated reserve: 6 hours."

### 10.3 Network Infrastructure

Players build connections between stations:

**Hardline links** — physical connections through tunnel infrastructure. Fast, reliable, but only possible between geographically close stations. Require engineering at both ends.

**Wireless links** — radio connections through surface repeaters. Medium range, affected by weather and terrain. Require surface hardware.

**Satellite relay** — long-range connections through the shared constellation. Higher latency, subject to orbital availability and constellation health. The backbone of the network.

**The network itself is infrastructure that needs maintenance.** Links degrade. Satellites drift. Relay nodes need power. The network is not a given — it's something the community builds and maintains through engineering.

### 10.4 Shared Systems

Players can build infrastructure that spans stations:

- **Distributed power grids** — surplus from one station feeds another. Requires load balancing, failover logic, and trust.
- **Shared sensor networks** — pooled environmental data gives broader coverage. Requires data aggregation and conflict resolution.
- **Communications relay chains** — extend reach beyond direct satellite range. Requires routing protocols.
- **Resource trade networks** — supply chain infrastructure connecting stations with complementary resources.
- **Shared fabrication queues** — distributed manufacturing across multiple stations.

All shared systems are defined in AXIOM Script, stored in the `/network/` namespace, and require active maintenance from participating players. They introduce distributed systems challenges: consensus, synchronization, latency handling, partition tolerance.

### 10.5 Player Dynamics

The full spectrum of player interaction emerges from the engineering layer:

**Cooperation**
Players collaborate on infrastructure too large or complex for one station. Building a regional power grid, maintaining the satellite constellation, constructing a shared fabrication network. Trust is built through reliability — your uptime, your code quality, your contribution to shared systems.

**Competition**
Resources are finite. Territory (geographic proximity to valuable infrastructure or resources) matters. Players compete for strategic positions, resource deposits, and satellite access windows. Competition is economic and infrastructural, not violent.

**Alliances & Factions**
Groups of players share code, resources, and infrastructure. A faction might develop a proprietary toolkit of AXIOM scripts — their competitive advantage. Internal trust, shared standards, mutual defense. Factions form organically around shared infrastructure dependencies.

**Sabotage & Warfare**
Conflict is conducted through engineering. Attack vectors include:

- Probing exposed network services for vulnerabilities
- Injecting corrupted data into shared sensor feeds
- Overloading shared power grids to cause cascading failures
- Deploying scripts that subtly siphon resources from shared systems
- Social engineering — gaining access to a faction's shared scripts through trust, then leaking or corrupting them

**Defense** is also engineering:

- Firewalls and access control on network services
- Anomaly detection scripts monitoring for unusual patterns
- Integrity checking on shared system configs
- Audit logs and intrusion detection
- Redundancy and failover so sabotage doesn't cascade

**Espionage**
A player's codebase is their intellectual property. Well-designed automation scripts, efficient resource algorithms, clever security measures — all valuable. Infiltrating a rival faction to access their `/scripts/` directory is a real play. Defecting with stolen code is a real betrayal.

### 10.6 Emergent Economy

No artificial currency. No auction house. Trade emerges from infrastructure:

- Players with surplus resources negotiate exchanges with players who need them
- The trade infrastructure itself must be built — routing, inventory sync, transport (drone convoys, pipeline networks)
- Specialization emerges naturally: some stations optimize for power, others for fabrication, others for network infrastructure
- "Backbone operators" emerge — players who specialize in maintaining satellite and relay infrastructure that everyone depends on, extracting value from their position

### 10.7 Civilization Scale

As hundreds of players come online:

- Regional clusters form around geographic proximity
- Inter-regional links depend on satellite infrastructure
- The satellite constellation becomes a critical shared commons — who maintains it, who controls access, who pays the resource cost
- Governance emerges from infrastructure dependency — not elected governments, but *engineering councils* where the people who maintain critical systems have influence proportional to their contributions
- The "health" of the civilization is literally visible: a network map where green means functioning infrastructure and red means failures. The state of the world is the state of the code

---

## 11. Narrative Design

### 11.1 Environmental Storytelling

There are no cutscenes, no dialogue trees, no quest markers. All narrative is embedded in the infrastructure:

- **Old log files** — `cat /archives/log-2847.txt` reveals an engineer's notes from decades ago. Personal, practical, sometimes desperate.
- **Corrupted archives** — decrypt and restore data fragments to piece together what happened. The decryption process itself is a programming challenge.
- **System behavior** — some systems behave unexpectedly. Investigating why reveals their history. A coolant loop with a strange routing wasn't a bug — it was a workaround for a section that flooded years ago.
- **The surface** — camera feeds, weather data, and eventually direct exploration (via drones) reveal the state of the world above. It tells its own story through what's there and what isn't.
- **Other players' stations (Part 2)** — every station has its own archaeology. When you connect with another player, their station tells a parallel story through different archives, different damage, different engineering decisions by previous occupants.

### 11.2 Mystery Arc

The overarching mystery — what happened, why are you here, what is the purpose of these installations — is never fully explained. Fragments accumulate. Players on the network share discoveries, piece together a larger picture. The truth, if it exists, is a collaborative discovery.

### 11.3 Event System

Narrative events are driven by preconditions on world state, not scripted timelines:

- Low sensor coverage + late game cycle → "Anomalous reading on perimeter sensors"
- Restored surface cameras → first visual of the world above
- Multiple communication failures → "Pattern detected in signal interference — not random"
- Network scale reaches threshold → discovery of a dormant mega-installation that requires massive collaborative effort to activate

Events create pressure and story without interrupting gameplay. They appear in the event log and on the dashboard. How (and whether) you respond is up to you.

---

## 12. Technical Architecture

### 12.1 Technology Stack

**Frontend (Godot / C#):**
- **Engine:** Godot 4 with C#
- **Role:** All rendering, UI, audio, input handling, visual effects
- **Includes:** Dashboard rendering (including telemetry graphs), retro CRT aesthetic, in-game terminal emulator, text editor, visual effects (scan lines, flicker, glow), sound design, Steam integration
- **Distribution:** Native executables via Godot export (Windows, Linux, Mac)

**Backend Engine (Go):**
- **Role:** Simulation engine, config parser, command engine, game logic
- **Integration:** Compiled as a C shared library via CGO (`go build -buildmode=c-shared`), called from C# via P/Invoke or GDExtension
- **Includes:** Tick-based simulation loop, subsystem models, dependency graph, cascade engine, event/threat system, config parser, virtual filesystem logic, command parser and execution
- **Scripting Engine:** Custom tree-walk interpreter for AXIOM config language. For trigger/automation scripting, either custom or embedded Lua via `gopher-lua` — TBD based on complexity needs.

**Server (Part 2 — Go):**
- **Protocol:** WebSocket or custom TCP for station-to-station communication
- **Persistence:** SQLite per station for state, shared database for network/constellation state
- **Hosting:** Dedicated server binary, horizontally scalable via regional sharding

### 12.2 Architecture Split: What Lives Where

**Go (the brain)** owns all game logic. It is the single source of truth for simulation state. Godot never modifies game state directly — it sends commands to Go and receives state updates back.

- Simulation tick loop and world clock
- All subsystem models, sensors, dependencies
- Dependency graph and cascade propagation
- Config parsing and application (declarative .ax files)
- Command parsing, validation, and execution
- Virtual filesystem (game state as files/directories)
- Telemetry export (CSV) for Godot graph rendering
- Event and threat engine (precondition evaluation, event firing)
- Procedural problem generation (wear, corruption, emergent bugs)
- Network protocol and inter-station communication (Part 2)
- Server-side simulation for offline persistence (Part 2)

**Godot/C# (the face)** owns all presentation and player interaction. It renders the world Go computes.

- AXIOM OS visual shell — retro CRT monitor aesthetic, scan lines, phosphor glow, screen flicker
- Dashboard rendering — dynamic panel layout, color-coded status indicators, drill-down views
- Telemetry graphs — reads CSV telemetry from the Go engine, renders time-series graphs of subsystem values
- In-game text editor with AXIOM config syntax highlighting
- Terminal emulator with tab completion UI, command history, scrollback
- Sound design — ambient hums, alert klaxons, keyboard clatter, system boot sequences
- Visual effects — power fluctuations dimming the screen, alerts flashing, panels lighting up as systems come online
- Station environment — if you want a visual representation of the bunker beyond just terminals
- Network map visualization (Part 2)
- Steam integration, achievements, settings, save management
- Input handling and routing commands to Go

### 12.3 The FFI Boundary

Go compiles to a C-compatible shared library via CGO. C# calls into it via P/Invoke (or optionally via GDExtension for tighter Godot integration).

The API boundary is clean and narrow:

```
// Go exposes these to Godot via CGO:
axiom_init()                          → Initialize simulation
axiom_tick()                          → Advance one simulation tick
axiom_get_state()                     → Full world state as JSON
axiom_get_subsystem(id)               → Detailed subsystem state
axiom_execute_command(cmd_string)      → Player typed a command
axiom_save_file(path, contents)        → Player saved a config file
axiom_get_filesystem(path)             → List directory or read file
axiom_get_completions(partial)         → Tab completion suggestions
axiom_get_dashboard()                  → Dashboard panel data
axiom_get_events()                     → Recent events/alerts since last poll
```

Data crosses the boundary as JSON. Godot deserializes it into C# objects for rendering. Telemetry data is written to CSV files that Godot tails for graph rendering.

This separation means:
- Go can be developed and tested independently (unit tests, headless simulation, terminal REPL)
- Godot can be developed with mock data while Go features are in progress
- The Go engine powers a standalone terminal REPL for MVP playtesting without Godot
- The Go server for Part 2 shares the exact same simulation code as the client

```
┌──────────────────────────────────────────────────────────┐
│                    Godot / C# (Frontend)                  │
│  ┌────────────┐ ┌───────────┐ ┌────────────────────────┐ │
│  │ Dashboard   │ │ Terminal  │ │  Editor / Effects /    │ │
│  │ + Graphs    │ │ Emulator  │ │  Audio / Environment   │ │
│  └─────┬──────┘ └─────┬─────┘ └──────────┬─────────────┘ │
│        │              │                   │               │
│  ══════╪══════════════╪═══════════════════╪═══════════    │
│        │    FFI Boundary (CGO / P/Invoke) │               │
│  ══════╪══════════════╪═══════════════════╪═══════════    │
│        │              │                   │               │
│  ┌─────▼──────────────▼───────────────────▼─────────────┐ │
│  │               Go Engine (Shared Library)               │ │
│  │  ┌──────────┐ ┌───────────┐ ┌──────────────────────┐ │ │
│  │  │Simulation│ │  Config   │ │  Virtual Filesystem   │ │ │
│  │  │  Engine  │ │  Parser   │ │  & Command Engine     │ │ │
│  │  │          │ │           │ │                       │ │ │
│  │  └──────────┘ └───────────┘ └──────────────────────┘ │ │
│  │  ┌──────────┐ ┌───────────┐ ┌──────────────────────┐ │ │
│  │  │Dependency│ │  Event /  │ │  Telemetry Export     │ │ │
│  │  │  Graph   │ │  Threat   │ │  (CSV for Godot)      │ │ │
│  │  └──────────┘ └───────────┘ └──────────────────────┘ │ │
│  └──────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────┘
```

### 12.4 Two Layers of Player Input

**Commands** are the primary gameplay interface. The player types commands in the terminal to inspect, diagnose, and manipulate the station in real time. Commands are imperative — they do something now.

Core commands: `status`, `inspect`, `diagnose`, `help`, `ls`, `cat`, `write`, `apply`, `set`

**Config files** are the deeper interaction layer. When a command reveals a problem ("pump speed set to 40%"), the player navigates the virtual filesystem, opens the relevant `.ax` file, and edits it. Config files are declarative — they describe how systems should be wired and what their setpoints are.

The config language handles:
- Subsystem declarations (`system`)
- Connection topology (`connect` with named roles and throughput)
- Setpoints and thresholds (`set` component values)

The player edits config files via the VFS (`write`), then runs `apply` to re-parse and apply the config to the live simulation. Errors are surfaced through `diagnose`.

The command engine is higher priority than the config parser. Commands enable all gameplay; config editing enables the "fix the broken config" subset.

### 12.5 Simulation Tuning

Subsystem behavior is governed by tuning profiles rather than inline constants. Each subsystem defines its response characteristics as data:

- **Gain** — fraction of the gap between current and target value closed per tick
- **Ceiling** — maximum absolute delta per tick (prevents runaway)
- **Floor** — minimum drift per tick (prevents stalling)

Profiles are defined per input role on each subsystem. The tick body becomes: for each input role, look up the profile, compute `profile.Delta(current, input)`, accumulate deltas. This makes every subsystem's physics consistent and tunable without rewriting formulas.

Loop stability is verifiable by inspection: if the product of gains around any feedback loop is less than 1.0, the system converges. If greater than 1.0, it diverges. This replaces ad-hoc balancing with a predictable framework.

Telemetry is exported to CSV each tick (`tick,system,component,value`) for Godot to render as time-series graphs and for the developer to plot during tuning.

### 12.6 AXIOM Script Execution (Future)

Beyond config files, the full AXIOM Script language (triggers, macros, test scripts) will run in a sandboxed interpreter with:

- **Resource limits** — max ticks per execution, max memory, max concurrent scripts
- **Permissions** — scripts can only access systems they're authorized for. Network scripts need explicit network permissions.
- **Hot reload** — filesystem watcher detects changes and reloads scripts without restarting the simulation
- **Error containment** — a crashing script logs the error and deactivates. It doesn't take down the station.
- **Deterministic execution** — given the same state and inputs, scripts produce the same outputs. Critical for debugging and for server validation in Part 2.

### 12.5 Server Architecture (Part 2)

Each station runs as a simulation instance on the server. The server:

- Ticks all stations continuously (whether players are online or not)
- Validates script execution (prevents cheating)
- Manages inter-station networking (message passing, shared state)
- Handles the satellite constellation as a shared simulation
- Persists all state to disk for durability

Scale considerations:
- Station simulation is lightweight per-instance (mostly state machine updates)
- Script execution is the heaviest workload — sandboxed, resource-limited
- Inter-station traffic is relatively low-bandwidth (state updates, commands, sensor data)
- Horizontal scaling: regional sharding based on geographic proximity in the game world

---

## 13. Development Roadmap

### Milestone 1 — The Heartbeat (Done)
**Goal:** A ticking simulation with basic systems.

- ~~Set up Go module with package structure~~
- ~~Implement tick-based simulation loop with 3 subsystems (power, coolant, HVAC)~~
- ~~Connection system with ports, throughput, and DFS dependency resolution~~
- ~~Logging system with structured output~~
- Systems interact through a dependency graph — the world is alive

### Milestone 2 — The Hands
**Goal:** Player can interact with the simulation through a terminal REPL and config files.

- Role-based connection routing (connections declare destination role, not just component type)
- Tuning profiles replace inline constants (gain/ceiling/floor per input role)
- Config parser: `system`, `set`, `connect` directives from .ax files
- WorldState.ApplyConfig() replaces hardcoded Init()
- VFS wired to live game state (virtual readers for status, writable config files)
- Command engine: `status`, `inspect`, `diagnose`, `ls`, `cat`, `write`, `apply`, `set`, `help`
- REPL on main goroutine, simulation ticking in background goroutine
- Telemetry CSV export for Godot graph rendering
- MVP scenario: broken config → diagnose → fix → monitor recovery → add subsystem
- Prove the core loop is fun before building more

### Milestone 3 — The Dashboard
**Goal:** Godot renders the simulation state as a visual dashboard with telemetry graphs.

- Set up Godot project with C# and P/Invoke bindings to Go CGO shared library
- Build basic dashboard scene — panels showing live values from the Go engine
- Telemetry graph rendering from CSV data
- Build terminal emulator UI in Godot (retro CRT aesthetic, text input, scrollback)
- Godot sends command strings to Go, renders results
- In-game text editor for config files with syntax highlighting

### Milestone 4 — The Brain
**Goal:** AXIOM Script interpreter running player-written automation.

- Implement script parser and interpreter in Go for triggers and macros
- Trigger evaluation on each tick
- Script management: `deploy`, `undeploy`, `list-scripts`, `script-status`
- Hot reload on file save
- Error handling and sandboxing in Go, error display in Godot
- First real gameplay loop: system breaks → player writes fix → script handles it next time

### Milestone 5 — Cascades & Chaos
**Goal:** Dependency graph creates emergent, cascading failures.

- Full dependency graph implementation
- Cascade propagation engine
- Complex multi-system failure scenarios emerge naturally
- Player must triage and prioritize — the core tension loop

### Milestone 6 — Growth
**Goal:** Player can expand their station and watch it grow.

- Sealed sections that can be explored and activated
- Fabrication system for building components
- New subsystem types unlocked through expansion
- Dashboard dynamically grows with the station — Godot adds panels as Go reports new subsystems
- Test scripts and health checks as first-class features
- Visual and audio polish: boot sequences, alert sounds, ambient hums, CRT effects

### Milestone 7 — Narrative
**Goal:** The world tells a story through its infrastructure.

- Archive system with discoverable lore fragments
- Event engine in Go with precondition-based triggers
- Surface access: camera feeds rendered in Godot, sensor data from Go
- Mystery elements seeded throughout the station
- Threat events: environmental, mechanical, and ambiguous
- Environmental art and atmosphere in Godot (station visuals, lighting, mood)

### Milestone 8 — The Uplink
**Goal:** Satellite system and the transition to multiplayer.

- Satellite uplink facility as a late-game engineering challenge
- Orbital mechanics simulation in Go (simplified)
- Signal processing and antenna management
- Godot renders satellite tracking interface, signal visualizations
- First contact: detecting another station's signal
- Steam store page, Godot export pipeline, packaging and distribution

### Milestone 9 — The Network
**Goal:** Persistent multiplayer world.

- Go server binary running persistent station instances (shares simulation code with client)
- Inter-station communication protocol
- Offline automation (Go server ticks stations while players are away)
- Shared systems namespace and network AXIOM extensions
- Basic resource trading between stations
- Network map visualization in Godot

### Milestone 10 — Civilization
**Goal:** Full-scale emergent multiplayer dynamics.

- Satellite constellation as shared infrastructure
- Security, access control, and intrusion detection mechanics
- Faction support (shared codebases, private namespaces)
- Sabotage and defense mechanics
- Scale testing with hundreds of concurrent stations
- Server infrastructure, deployment, and operations

---

## 14. Inspirations

- **Dwarf Fortress** — emergent complexity from simulated systems
- **EVE Online** — player-driven economy, politics, and warfare at civilization scale
- **Zachtronics games (TIS-100, Shenzhen I/O, Exapunks)** — programming as core gameplay
- **FTB / Modded Minecraft** — discoverable complexity, automation as progression, build-and-monitor satisfaction
- **Arc Raiders / Marathon** — post-collapse tech aesthetic, discovering remnants of an advanced world
- **Factorio** — the factory must grow; scaling systems and watching them work
- **Unix philosophy** — everything is a file, small tools composed together, text as universal interface
- **Kubernetes / Docker / Grafana** — the real-world satisfaction of dashboards, green status indicators, and infrastructure-as-code

---

## 15. The Sandbox (Post-Crisis Creative Layer)

*Note: This section describes the long-term vision for AXIOM's creative sandbox. It is NOT part of the prototype or MVP. The prototype focuses exclusively on the core crisis loop. However, the AXIOM Script language and system architecture should be designed with these capabilities in mind to avoid costly rewrites later.*

### 15.1 The Minecraft Transition

The game has two phases of engagement. The first is survival: fix the crises, stabilize the station, learn the systems. The second begins when the station is stable and the question shifts from "how do I survive?" to "what can I build?"

Just as Minecraft transitions from surviving the first night to building whatever you can imagine, AXIOM transitions from firefighting to creative engineering. The tools you learned under pressure become a medium for expression.

### 15.2 What the Sandbox Enables

**Custom tools and commands.** Players compose new commands from existing primitives. A `deep-scan` command that audits all systems and ranks them by failure risk. A `report` command that generates a daily station summary. These become part of the player's personal OS layer.

**Custom dashboard panels.** Players define their own monitoring views — computed values, trend visualizations, risk boards, efficiency gauges. Two players with identical stations will have completely different dashboards because they built different tools.

**Intelligent automation.** Beyond simple triggers — full autonomous systems. An intelligent power manager that predicts demand. A self-healing infrastructure that diagnoses and fixes problems while you sleep. A resource optimizer that balances consumption against reserves. Each is a genuine programming project the player iterates on.

**Data processing and analysis.** Pipe sensor streams through filters, aggregators, and transformers. Track trends, detect anomalies, predict failures. The data layer is where creative programming really shines.

**Non-functional creativity.** ASCII art generators, procedural status broadcasts, visualizations, games-within-the-game. If someone wants to build a cellular automaton in AXIOM Script, that should be possible.

### 15.3 Language Requirements for Sandbox

For the sandbox to work, AXIOM Script eventually needs to grow beyond triggers and configs to include: variables and data structures, iteration and loops, string manipulation, math and statistics, custom command definitions, custom dashboard panel definitions, file I/O for logging and inter-script communication, and piping/composition to chain scripts together. The game world's physics and resource constraints provide the game design boundaries — the language itself should not be the limiting factor.

### 15.4 Social Creativity (Part 2)

In multiplayer, the sandbox becomes social. Visit another player's station and see their tools, their dashboard, their automation philosophy. Trade scripts. Adapt ideas. Develop faction-wide toolkits and standards. A player's codebase becomes their identity and reputation.

---

## 16. Open Questions

- **AXIOM Script implementation:** Custom tree-walk interpreter in Go (full control, fits the constrained language) vs. embedded Lua via `gopher-lua` (faster to prototype triggers/automation, proven sandboxing)? Config parsing is simple either way — the question is whether trigger/automation scripting justifies Lua.
- **FFI strategy:** P/Invoke from C# via CGO shared library (simplest) vs. GDExtension (tighter integration). Start with P/Invoke.
- **Data serialization across FFI:** JSON for development. Optimize later if profiling demands it.
- **Command parser vs config parser priority:** Commands first — they enable all gameplay. Config parser enables the "edit the broken file" puzzle type.
- **Tick rate tuning:** 1/second feels right for tension, but may need adjustment for offline simulation at scale.
- **Subsystem tuning methodology:** Telemetry CSV export + plotting for feedback behavior. Data-driven tuning profiles to replace inline constants. Loop gain product < 1.0 for stability.
- **Satellite mechanics depth:** How realistic should orbital mechanics be? Simplified model vs. actual Keplerian elements?
- **Anti-cheat in Part 2:** Server-side script validation is essential, but how much can be trusted to the client?
- **Onboarding:** The first broken config IS the tutorial. Pacing matters — the first fix should be obvious, the second should require investigation.
- **Visual scope in Godot:** Is the game purely terminal/dashboard screens, or does the player also see a visual representation of their station? Starting terminal-only is safer; visual environments can be added later.
- **Scope management:** Part 1 is a complete game. Part 2 is ambitious. The roadmap should ensure Part 1 ships and is satisfying on its own before Part 2 development begins.

---

## 17. Prototype: Prove the Fun

Before building the full game, build the smallest thing that answers: **is the core loop fun?** The core loop is: see a problem on the dashboard → investigate through the terminal → fix it by editing code/config → watch the dashboard recover → something else breaks. If that loop is satisfying, the game works. Everything else is content and scale.

### What the Prototype Needs

**Go Engine (terminal REPL, no Godot yet):**

- [x] Project setup: Go module with package structure, tick loop, logging
- [x] Three subsystems wired up: **Power** (generates energy, produces heat), **Coolant** (manages temperature via flow), **HVAC** (regulates ambient temperature)
- [x] Dependency graph: coolant feeds power (cooling), power feeds HVAC (power + heat)
- [ ] Role-based connections: connections declare destination role, not just component type
- [ ] Tuning profiles: gain/ceiling/floor per input role, replacing inline constants
- [ ] Config parser: `system`, `set`, `connect` directives from .ax files
- [ ] WorldState.ApplyConfig() replacing hardcoded Init()
- [ ] Virtual filesystem wired to live game state (virtual readers for status, writable config files)
- [ ] Command engine: `status`, `inspect`, `diagnose`, `ls`, `cat`, `write`, `apply`, `set`, `help`
- [ ] REPL: simulation goroutine + stdin command loop
- [ ] Telemetry CSV export: tick-by-tick component values for Godot graph rendering
- [ ] Broken starting config that creates an obvious problem to diagnose and fix

**Godot Frontend (after core loop is validated):**

- [ ] Project setup: Godot 4 C# project, P/Invoke bindings to Go CGO shared library
- [ ] Dashboard scene: panels showing live values from Go state, color coded
- [ ] Telemetry graphs: render time-series from CSV data
- [ ] Terminal scene: text input, scrollback, command history. Send commands to Go, display results
- [ ] In-game file editor for config files
- [ ] Boot sequence and alert system
- [ ] CRT shader (optional but motivating)

**The First Playable Scenario:**

- [ ] Player boots into a station where HVAC is non-functional (power-in connection throughput is 0.0)
- [ ] `status` shows HVAC critical, ambient temp rising
- [ ] `diagnose hvac` reveals: "power-in throughput is 0.0 — system receives no power"
- [ ] Player reads config via `cat /station/config.ax`, spots the zero throughput
- [ ] Player fixes it via `write` + `apply`, watches HVAC temp converge toward target
- [ ] Power is running hot at effort 0.7 — a second problem emerges naturally from the physics
- [ ] Player adds a second cooling unit via config (`system cooling2 type=cooling` + connections)
- [ ] The loop: diagnose → fix config → monitor recovery → new problem emerges

### Success Criteria

The prototype is successful if:

- [ ] A playtester can figure out what to do without external instructions
- [ ] Fixing a problem and watching values recover feels satisfying
- [ ] The player feels a sense of urgency when multiple things go wrong at once
- [ ] The player wants to keep playing after the first crisis is resolved
- [ ] Editing configs feels like "real" engineering, not busywork
- [ ] Adding a new subsystem via config and watching it take effect feels powerful
