package connections

import (
	"github.com/elias/axiom/engine/subsystems"
	"github.com/elias/axiom/engine/utils"
)

type ConnectionID int64

type Connection struct {
	id         ConnectionID
	src        *subsystems.OutputPort
	dest       *subsystems.InputPort
	srcSystem  string
	destSystem string
	throughput utils.Unit
}

func (c Connection) ID() ConnectionID            { return c.id }
func (c Connection) Src() *subsystems.OutputPort { return c.src }
func (c Connection) Dest() *subsystems.InputPort { return c.dest }
func (c Connection) SrcSystem() string           { return c.srcSystem }
func (c Connection) DestSystem() string          { return c.destSystem }
func (c Connection) Throughput() utils.Unit      { return c.throughput }

func NewConnection(id ConnectionID, src *subsystems.OutputPort, dest *subsystems.InputPort, srcSystem string, destSystem string, throughput utils.Unit) *Connection {
	return &Connection{id, src, dest, srcSystem, destSystem, throughput}
}
