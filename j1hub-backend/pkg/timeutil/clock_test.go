package timeutil_test

import (
	"testing"
	"time"

	"github.com/j1hub/backend/pkg/timeutil"
	"github.com/stretchr/testify/assert"
)

func TestRealClock_Now(t *testing.T) {
	c := timeutil.RealClock{}
	now := c.Now()
	assert.WithinDuration(t, time.Now(), now, 1*time.Second)
}

func TestMockClock(t *testing.T) {
	initial := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	c := &timeutil.MockClock{CurrentTime: initial}

	assert.Equal(t, initial, c.Now())

	c.Set(initial.Add(1 * time.Hour))
	assert.Equal(t, initial.Add(1*time.Hour), c.Now())

	c.Add(30 * time.Minute)
	assert.Equal(t, initial.Add(90*time.Minute), c.Now())
}
