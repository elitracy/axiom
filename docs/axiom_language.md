# Axiom Language

### Keywords
- source - the system
- output - the output for each component of a system
- distribute - how a resource of the system will be distrbuted (output)
- fallback - the next system when the current system fails
- status - the current state of the system

### Config Example

```
source reactor.main {
    output power 2400w
}

distribute reactor.main { // outputs
    hvac socket-1 100%
    lifesupport socket-2 100%
    lifesupport socket-3 100%
    reserve socket-4 100%

    air valve-1 100%
} 

collect reactor.main { // inputs
    cooling.main valve-2 100%
}

source cooling.main {
    output flow 100%
}

distribute cooling.main {
    reactor.main valve-1 100%
}

source hvac.main {
    output temperature 20 // celsius
}

collect hvac.main {
    air valve-1 100%
}

distribute hvac.main {
    air valve-2 100%
}

```
