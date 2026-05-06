package utils

import "sync/atomic"

type SubsystemType int

//go:generate stringer -type=SubsystemType
const (
	Power SubsystemType = iota
	Cooling
	Hvac
)

type SubsystemName string

type PortType int

const (
	PortInput PortType = iota
	PortOutput
)

type Status int

//go:generate stringer -type=Status
const (
	Healthy Status = iota
	Warning
	Critical
	Offline
)

type Tick struct {
	Val atomic.Int64
}

func (t *Tick) Tick() int64 { return t.Val.Load() }

func NewTick() *Tick { return &Tick{Val: atomic.Int64{}} }
