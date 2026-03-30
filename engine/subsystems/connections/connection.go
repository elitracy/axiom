package connections

import (
	"github.com/elias/axiom/engine/subsystems"
	"github.com/elias/axiom/engine/utils"
)

type ConnectionID int64

var (
	currentConnectionID = 0
)

func newConnectionID() ConnectionID {
	id := currentConnectionID
	currentConnectionID++
	return ConnectionID(id)
}

type Connection struct {
	id         ConnectionID
	src        *Port
	dest       subsystems.SubsystemID
	throughput utils.Norm
}

func (c Connection) ID() ConnectionID             { return c.id }
func (c Connection) Src() *Port                   { return c.src }
func (c Connection) Dest() subsystems.SubsystemID { return c.dest }
func (c Connection) Throughput() utils.Norm       { return c.throughput }

func NewConnection(src *Port, dest subsystems.SubsystemID, throughput utils.Norm) *Connection {
	return &Connection{
		id:         newConnectionID(),
		src:        src,
		dest:       dest,
		throughput: throughput,
	}
}
