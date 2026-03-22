package systems

type Status int8

const (
	Offline Status = iota
	Critical
	Degraded
	Online
)

type SystemID struct {
	ID int
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
		id:   SystemID{ID: 0},
		name: name,
	}
}

func (s *SystemCore) ID() SystemID { return s.id }
func (s *SystemCore) Name() string { return s.name }
