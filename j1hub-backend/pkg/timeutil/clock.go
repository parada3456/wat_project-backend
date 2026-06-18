package timeutil

import "time"

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now()
}

type MockClock struct {
	CurrentTime time.Time
}

func (m *MockClock) Now() time.Time {
	return m.CurrentTime
}

func (m *MockClock) Set(t time.Time) {
	m.CurrentTime = t
}

func (m *MockClock) Add(d time.Duration) {
	m.CurrentTime = m.CurrentTime.Add(d)
}
