package utils

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
