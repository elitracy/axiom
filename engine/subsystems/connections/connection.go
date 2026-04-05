package connections

import (
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
	dest       *Port
	throughput utils.Unit
}

func (c Connection) ID() ConnectionID       { return c.id }
func (c Connection) Src() *Port             { return c.src }
func (c Connection) Dest() *Port            { return c.dest }
func (c Connection) Throughput() utils.Unit { return c.throughput }

func NewConnection(src *Port, dest *Port, throughput utils.Unit) *Connection {
	return &Connection{
		id:         newConnectionID(),
		src:        src,
		dest:       dest,
		throughput: throughput,
	}
}
