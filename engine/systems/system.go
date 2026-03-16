package systems

const (
	TICKS_TILL_DEATH_DEBUG = 10
)

type Status int8

const (
	Offline Status = iota
	Critical
	Degraded
	Online
)

type SystemID struct {
	id int
}

type System interface {
	ID() SystemID
	Name() string
	String() string
	Status() Status
}

type SystemCore struct {
	id   SystemID
	name string
}

// Creates a System
// name is the name of the system
func NewSystemCore(name string) *SystemCore {
	return &SystemCore{
		// TODO: generate system IDs dynamically
		id:   SystemID{id: 0},
		name: name,
	}
}

func (s *SystemCore) ID() SystemID { return s.id }
func (s *SystemCore) Name() string { return s.name }
