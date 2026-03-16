package systems_test

import (
	"testing"

	"github.com/elias/axiom/engine/systems"
	"github.com/stretchr/testify/assert"
)

func TestSystem(t *testing.T) {
	system := systems.NewSystemCore("Test System")

	assert.Equal(t, system.Name(), "Test System")
	assert.Equal(t, system.ID(), systems.SystemID{ID: 0})

}
