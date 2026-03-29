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
	throughput utils.Norm
}

func NewConnection(src, dest *Port, throughput utils.Norm) *Connection {
	return &Connection{
		id:         newConnectionID(),
		src:        src,
		dest:       dest,
		throughput: throughput,
	}
}
