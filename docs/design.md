# Axiom 

This document outlines the structure of the axiom game, it's underlying systems, and the interface (domain specific langauge) that is 
used to manipulate these systems.

## Game
The core of this game is built around survival and your ability to solve problems using the domain specific language (AXIOM). The AXIOM language is the user's
way of interfaces with the necessary systems in place which:
 1. keep the base running
 1. ensure your player's survival
 1. help your base expand
 1. ward off external threats

#### Lore
As a player you start the game incenerating your sibling. What was once a thriving underground community has become a barren, barely functional, post-apoctolyptic 
jail cell. Machines that were meant to work for hundreds of years have only proved functional for half that time. Around a decade prior your community began
siloing themselves off from other communities and rationing necesary resources amongst its population. In the time since, your community has shrunk from around 200 
people to 1.

In an extinction prevention effort your community began not only teaching its youngest citizens the necessary skills to survive but also allocated more resources for
them. First the elderly sacrificed their bodies for the community could fuel, water, and food, then followed the latter generations, until all that was left was your
sibling and yourself. At 20 years old they knew they had a responsibility to not only their race but to you. As the youngest member of your colony you are all that 
remains. 

You start the game by recycling your sibling's matter into all the remaining resources you can. As you are well aware by now the systems in your bunker are  
not manipulatable by human touch (prior generations favored automation and computer interfaces for these types of functionalities), so if you want to fix something 
you are going to have to use the terminal interface in your bunker's control center. 

The control center is a single desk scattered with old monitors and wires which have been recklessly unkept. As you sit at the desk your peripheral is filled with
the glow of log lines, metrics, and programs running. Overwhelmed but persistent you begin solving the first issue... power generation.

## Systems

All systems need a way to be interacted with:
 - The power system should have a way to turn up or down the output.
 - The coolant should have a way to increase or decrease the flow.
 - The life support system should have a way to adjust the scrubbing rate
 - The food generation system should have a way to adjust how much food is being made
 - The water system should have a way to adjust how much water is being generated

| Need | Function | 
| --- | --- |
| Unique identifier | ID |
| Human readable identifier | Name |
| Update different inputs (components) | UpdateComponent[T] | 
| Get component info | ComponentInfo[T] |
| Read relavent telemetry/logs (the what) | Telemetry[T]  |
| Get diagnostics (the why) | Diagnostics[T] |
| Enable system | Enable |
| Disable system | Disable |
| Update system | Tick |


I believe I will also need a system manager that can register systems and their components. For instance the power system will have a fuel component, a thermal component, 
and a power component. 

### Power
Power can be generated from several different sources. It can produce energy from organic matter at a highly efficient rate, but organic matter is difficul to 
come by. It can also use fuel (such as oil or gasoline which are stocked in reserves), steam, solar power, or radiation. At the start of the game the user  
will be using organic matter but will quickly transition to oil. 

| Type | Name | Efficiency | Components |
| --- | --- | --- | --- |
| Organic matter | bio-generator | high | fuel, thermal, power |  
| Radiation | reactor | high | fuel, thermal, power, radiation | 
| Oil | generator | moderate | fuel, thermal, power |
| Solar Power | solar-panel | moderate | direction ,power |
| Steam | steam-engine | low | fuel (water, coal), pressure, thermal, power |


### Cooling
Cooling for power generation can be done using a "coolant". The coolant you use will differ per thermal source.

| Thermal Source | Coolant Type |
| --- | --- |
| Oil | propylene-glycol | 
| Steam | water |
| Reactor | water | 


**Coolant Type**
| Type | Thermal Source | Components | 
| --- | --- | --- |
| Propylene-glycol (C3_H8_O2) | oil | thermal, viscosity, pressure |
| Water | steam, reactor | thermal, pressure |


### Life Support
Life support is the mechanism for how the base converts carbon dioxide into oxygen, recycles water, and regulates temperature. The reading form the life support system
will essentially be how the user perceives the temperature in their base. 

| Type | Components |
| --- | --- |
| Basic | scrubber, power, water, thermal |


### Food
*TBD*

## Language

