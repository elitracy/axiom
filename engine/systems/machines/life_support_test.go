package machines

import (
	"testing"

	"github.com/elias/axiom/engine/systems"
	"github.com/stretchr/testify/assert"
)

func TestLifeSupport_New(t *testing.T) {
	ls := NewLifeSupport()
	assert.Equal(t, ls.Name(), "Life Support")
	assert.Equal(t, ls.Status(), systems.Online)
	assert.Equal(t, ls.powerCapacity, 600.0)
}

func TestLifeSupport_Tick(t *testing.T) {
	tests := []struct {
		name     string
		power    float64
		temp     systems.Status
		oxygen   systems.Status
		expected systems.Status
	}{
		{"No power", 0, systems.Online, systems.Online, systems.Offline},
		{"No power - degraded inputs", 0, systems.Degraded, systems.Critical, systems.Offline},
		{"Both online", 600, systems.Online, systems.Online, systems.Online},

		{"Temp degraded", 600, systems.Degraded, systems.Online, systems.Degraded},
		{"Oxygen degraded", 600, systems.Online, systems.Degraded, systems.Degraded},
		{"Both degraded", 600, systems.Degraded, systems.Degraded, systems.Degraded},
		{"Temp critical", 600, systems.Critical, systems.Online, systems.Critical},
		{"Oxygen critical", 600, systems.Online, systems.Critical, systems.Critical},
		{"Temp critical oxygen degraded", 600, systems.Critical, systems.Degraded, systems.Critical},
		{"Both critical", 600, systems.Critical, systems.Critical, systems.Critical},
		{"Temp offline", 600, systems.Offline, systems.Online, systems.Offline},
		{"Both offline", 600, systems.Offline, systems.Offline, systems.Offline},

		{"Low power", 100, systems.Online, systems.Online, systems.Online},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := NewLifeSupport()
			input := LifeSupportInput{
				PowerAvailable:    tt.power,
				TemperatureStatus: tt.temp,
				OxygenStatus:      tt.oxygen,
			}
			output := ls.Tick(input)
			assert.Equal(t, tt.expected, output.Status)
		})
	}
}
