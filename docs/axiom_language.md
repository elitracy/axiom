# Axiom Configuration Language

The configuration language is how the player declares subsystem wiring and setpoints. Config files live in the virtual filesystem as `.ax` files and are applied to the simulation via the `apply` command.

## MVP Format

Line-oriented. Three directives. Comments with `#`.

### Directives

**`system`** -- declare a subsystem
```
system <name> type=<type>
```

**`set`** -- set a component value or parameter
```
set <system>.<component> <value>
```

**`connect`** -- wire a source port to a destination subsystem with a named role and throughput
```
connect <system>.<port> -> <dest> <role> <throughput>
```

### Example Config

```
# Subsystem declarations
system power    type=power
system cooling  type=cooling
system hvac     type=hvac

# Component setpoints
set power.effort       0.5
set cooling.effort     0.5
set hvac.target-temp   0.2

# Connections: source.port -> dest  role  throughput
connect cooling.temp-out  -> power  coolant-temp  1.0
connect cooling.flow-out  -> power  coolant-flow  1.0
connect power.power       -> hvac   power-in      0.5
connect power.temp        -> hvac   heat-in       1.0
```

### How it maps to the engine

| Directive | Engine operation |
|-----------|-----------------|
| `system power type=power` | Subsystem factory creates a Power subsystem, registers as "power" |
| `set power.effort 0.5` | Looks up subsystem "power", sets component "effort" to 0.5 |
| `connect cooling.flow-out -> power coolant-flow 1.0` | Creates a port on cooling's "flow-out" component, creates a connection to power with role "coolant-flow" and throughput 1.0 |

### Connection roles

The `role` in a `connect` directive is the name the destination subsystem uses to identify this input. When a subsystem's `Tick()` runs, it receives inputs keyed by role name:

```
connect cooling.temp-out -> power coolant-temp 1.0
```

Power's Tick() reads this as `inputs["coolant-temp"]`. The role name is what makes the connection meaningful to the destination -- it knows "coolant-temp" affects temperature via a cooling profile, while "heat-source" would affect temperature via a heating profile.

### Throughput

A multiplier (0.0 to 1.0) applied to the source value before it reaches the destination. A throughput of 0.5 means the destination receives half the source's value. A throughput of 0.0 means the connection is effectively disconnected.

## Future: Full AXIOM Script

The MVP config language covers static wiring and setpoints. The full AXIOM Script language (triggers, macros, test scripts, control flow) is a separate system built on top of this foundation. See `AXIOM_Design_Document.md` section 7 for the full language vision.

### Keywords (future)
- `source` -- the system (replaces `system` with richer declaration)
- `output` -- component output specification with units
- `distribute` -- how a resource is distributed (output routing)
- `collect` -- how inputs are consumed
- `fallback` -- the next system when the current system fails
- `trigger` -- reactive automation blocks
- `test` -- validation scripts
- `macro` -- recorded command sequences
- `status` -- current state queries
