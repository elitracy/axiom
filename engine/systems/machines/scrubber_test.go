package machines

import (
	"testing"

	"github.com/elias/axiom/engine/systems"
	"github.com/stretchr/testify/assert"
)

func TestScrubber_Status(t *testing.T) {
	tests := []struct {
		name     string
		o2Norm   float64
		co2Norm  float64
		expected systems.Status
	}{
		{"O2 online", 0.55, 0.375, systems.Online},
		{"O2 degraded", 0.4, 0.375, systems.Degraded},
		{"O2 critical", 0.2, 0.375, systems.Critical},
		{"O2 offline", 0.0, 0.375, systems.Offline},

		{"CO2 degraded", 0.55, 0.5, systems.Degraded},
		{"CO2 critical", 0.55, 0.75, systems.Critical},
		{"CO2 offline", 0.55, 1.0, systems.Offline},

		{"Both degraded", 0.4, 0.5, systems.Degraded},
		{"O2 critical CO2 degraded", 0.2, 0.5, systems.Critical},
		{"O2 degraded CO2 critical", 0.4, 0.75, systems.Critical},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewScrubber()
			s.o2.SetNorm(tt.o2Norm)
			s.co2.SetNorm(tt.co2Norm)
			assert.Equal(t, tt.expected, s.Status())
		})
	}
}

func TestScrubber_TickSteadyState(t *testing.T) {
	tests := []struct {
		name           string
		power          float64
		ticks          int
		expectedO2Dir  string
		expectedCO2Dir string
		expectedStatus systems.Status
	}{

		{"Full power - stable", 600, 100, "stable", "stable", systems.Online},

		{"No power - still online", 0, 15, "down", "up", systems.Online},
		{"No power - degraded", 0, 25, "down", "up", systems.Degraded},
		{"No power - critical", 0, 50, "down", "up", systems.Critical},
		{"No power - offline", 0, 75, "down", "up", systems.Offline},

		{"Quarter power - online", 150, 20, "down", "up", systems.Online},
		{"Quarter power - degraded", 150, 30, "down", "up", systems.Degraded},
		{"Quarter power - critical", 150, 65, "down", "up", systems.Critical},

		{"Half power - online", 300, 30, "down", "up", systems.Online},
		{"Half power - degraded", 300, 45, "down", "up", systems.Degraded},
		{"Half power - critical", 300, 95, "down", "up", systems.Critical}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewScrubber()
			input := ScrubberInput{PowerAvailable: tt.power}

			first := s.Tick(input)
			for range tt.ticks - 2 {
				s.Tick(input)
			}
			last := s.Tick(input)

			switch tt.expectedO2Dir {
			case "up":
				assert.Greater(t, last.O2, first.O2)
			case "down":
				assert.Less(t, last.O2, first.O2)
			case "stable":
				assert.InDelta(t, first.O2, last.O2, 0.001)
			}

			switch tt.expectedCO2Dir {
			case "up":
				assert.Greater(t, last.CO2, first.CO2)
			case "down":
				assert.Less(t, last.CO2, first.CO2)
			case "stable":
				assert.InDelta(t, first.CO2, last.CO2, 0.001)
			}

			assert.Equal(t, tt.expectedStatus, s.Status())
		})
	}
}

