package utils

type SubsystemType int

//go:generate stringer -type=SubsystemType
const (
	Power SubsystemType = iota
	Cooling
	Machine
)

type Status int

//go:generate stringer -type=Status
const (
	Healthy Status = iota
	Warning
	Critical
	Offline
)
