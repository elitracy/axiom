package connections

import (
	"fmt"

	"github.com/elias/axiom/engine/subsystems"
	"github.com/elias/axiom/engine/utils"
)

type ConnectionID int64

type Connection struct {
	id         ConnectionID
	srcPort    *subsystems.OutputPort
	destPort   *subsystems.InputPort
	srcSystem  utils.SubsystemName
	destSystem utils.SubsystemName
	throughput utils.Unit
}

func (c *Connection) ID() ConnectionID                { return c.id }
func (c *Connection) SrcPort() *subsystems.OutputPort { return c.srcPort }
func (c *Connection) DestPort() *subsystems.InputPort { return c.destPort }
func (c *Connection) SrcSystem() utils.SubsystemName  { return c.srcSystem }
func (c *Connection) DestSystem() utils.SubsystemName { return c.destSystem }
func (c *Connection) Throughput() utils.Unit          { return c.throughput }

func NewConnection(id ConnectionID, src *subsystems.OutputPort, dest *subsystems.InputPort, srcSystem, destSystem utils.SubsystemName, throughput utils.Unit) *Connection {
	return &Connection{id, src, dest, srcSystem, destSystem, throughput}
}

func (c *Connection) String() string {
	return fmt.Sprintf("%s[%s] -> %s[%s] @ %.2f", c.srcSystem, c.srcPort.Name(), c.destSystem, c.destPort.Name(), c.throughput)

}
