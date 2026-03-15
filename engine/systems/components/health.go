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
		),
	}
}

// Returns the Status per the applied value of health
func (c *Health) Status() systems.Status {

	appliedHealth := c.ApplyValueCurve()
	switch {
	case appliedHealth <= healthOfflineThreshold:
		return systems.Offline
	case appliedHealth <= healthCriticalThreshold:
		return systems.Critical
	case appliedHealth <= healthDegradedThreshold:
		return systems.Degraded
	default:
		return systems.Online
	}
}
