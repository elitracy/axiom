# AXIOM

Terminal-based station management game. Programming is the gameplay.

## Architecture

- **Engine:** Go (under `engine/`) — simulation, game logic, command parsing
- **Frontend:** Godot — renders dashboard and terminal UI
- **Integration:** Engine will expose C shared library via CGO when ready. Not a priority yet.

## Engine Packages

- `simulation` — WorldState, tick loop
- `systems` — Subsystem definitions (power, coolant, life support)
- `filesystem` — In-memory VFS (ls, cat, write)
- `commands` — Command parser and execution

## Current Phase

Phase 1: Foundation. Scaffolding the Go engine. See `docs/MVP_TODO.md` for the full roadmap.

## Conventions

- User is building this themselves. Advise and review, don't implement.
- Keep packages decoupled. Systems shouldn't know about commands, commands use simulation as the entry point.
- JSON serialization for all state (stdlib `encoding/json`).
- Test with `go test`. Keep tests next to the code they test.
