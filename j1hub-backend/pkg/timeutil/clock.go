package timeutil

import (
	"log"
	"time"
)

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time {
	log.Println("debugprint: entering (RealClock).Now")
	return time.Now()
}

type MockClock struct {
	CurrentTime time.Time
}

func (m *MockClock) Now() time.Time {
	log.Println("debugprint: entering (*MockClock).Now")
	return m.CurrentTime
}

func (m *MockClock) Set(t time.Time) {
	log.Println("debugprint: entering (*MockClock).Set")
	m.CurrentTime = t
}

func (m *MockClock) Add(d time.Duration) {
	log.Println("debugprint: entering (*MockClock).Add")
	m.CurrentTime = m.CurrentTime.Add(d)
}
