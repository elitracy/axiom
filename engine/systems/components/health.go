package components

import "github.com/elias/axiom/engine/systems"

const (
	healthOfflineThreshold  = 0.0
	healthCriticalThreshold = 0.3
	healthDegradedThreshold = 0.7
)

type Health struct {
	*ComponentCore
}

// Creates a new health component
// initial is the initial health of the component,
func NewHealthComponent(initial float64) *Health {
	return &Health{
		ComponentCore: NewComponent(
			"Health (%%)",
			initial,
			0.0,
			1.0,
		),
	}
}

// Returns the Status per the applied value of health
func (c *Health) Status() systems.Status {

	normHealth := c.Norm()
	switch {
	case normHealth <= healthOfflineThreshold:
		return systems.Offline
	case normHealth <= healthCriticalThreshold:
		return systems.Critical
	case normHealth <= healthDegradedThreshold:
		return systems.Degraded
	default:
		return systems.Online
	}
}
