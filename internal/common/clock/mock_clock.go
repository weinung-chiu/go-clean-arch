package clock

import (
	"sync"
	"time"
)

type timerRequest struct {
	duration time.Duration
	ch       chan time.Time
	deadline time.Time
}

// MockClock 可用於測試：時間不會自動流動，必須由使用者主動推進
// After 行為模擬 Scheduler，會在 Advance 時推送值到對應 channel
type MockClock struct {
	mu     sync.Mutex
	now    time.Time
	timers []timerRequest
}

func NewMockClock(start time.Time) *MockClock {
	return &MockClock{
		now:    start,
		timers: nil,
	}
}

func (m *MockClock) Now() time.Time {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.now
}

func (m *MockClock) Sleep(d time.Duration) {
	m.Advance(d)
}

func (m *MockClock) After(d time.Duration) <-chan time.Time {
	m.mu.Lock()
	defer m.mu.Unlock()
	ch := make(chan time.Time, 1)
	m.timers = append(m.timers, timerRequest{
		duration: d,
		ch:       ch,
		deadline: m.now.Add(d),
	})
	return ch
}

// Advance 手動推進時間，觸發對應的 After channel
func (m *MockClock) Advance(d time.Duration) {
	m.mu.Lock()
	m.now = m.now.Add(d)
	var ready []timerRequest
	remaining := m.timers[:0] // reuse slice
	for _, t := range m.timers {
		if !m.now.Before(t.deadline) {
			ready = append(ready, t)
		} else {
			remaining = append(remaining, t)
		}
	}
	m.timers = remaining
	m.mu.Unlock()

	for _, t := range ready {
		t.ch <- t.deadline
		close(t.ch)
	}
}
