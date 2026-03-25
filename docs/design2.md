# Design 2 
Because the first one sucked

## SubSystems

| Subsystem | Connection Inputs | Connection Outputs | Params | Sensors |
| --- | --- | --- | --- | --- |
| Power | temperature | temperature, electriciy | effort level | temperature, effort level |
| Cooling | None | temperature | effort level | effort level |
| CO2 Scrubber | electricity | None | effort level | CO2, O2, effort level |
| HVAC | electricity, temperature | None | effort level, target temp| temperature, effort level |


```go

type SubSystem interface {
    ID() SubSystemID
    Name() string
    Effort() float64 // rate of generation
    Inputs() []ComponentType // temperature, power, etc
    Outputs() []ComponentType // tempertuare, power, etc
    Sensors() map[componentType]Component

    Tick(inputs map[ComponentType]Component)
}

type Component interface {
    Type() ComponentType
    Value() float64
}

type ComponentConnection struct {
    ComponentID
    SubSystemID
}

type System struct {
    subsystems map[SubSystemID]SubSystem
    dependencies map[SubSystemID][]ComponentConnection 
}


HVAC[0] -> Power[1] -> Cooling[2]
    (1, elec)    (2, temp)
    (1, temp)    (2, temp)
```

### Power 

coolTemp = 0
coolRate = -0.05
maxHeatDelta = 10
effort = 0.5
```

cooling = (coolTemp - currentTemp) * coolRate
heating = maxHeatDelta * effort
```

With an effort of .5 you can balance heat and cooling. Above that you start heating more every tick.
