package clock

import (
	"testing"
	"time"
)

func TestMockClock_NowAndAdvance(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	mc := NewMockClock(start)

	if mc.Now() != start {
		t.Fatalf("expected start time %v, got %v", start, mc.Now())
	}

	mc.Advance(2 * time.Hour)
	expected := start.Add(2 * time.Hour)
	if mc.Now() != expected {
		t.Errorf("expected time after advance %v, got %v", expected, mc.Now())
	}
}

func TestMockClock_After(t *testing.T) {
	start := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	mc := NewMockClock(start)

	ch := mc.After(1 * time.Hour)
	mc.Advance(30 * time.Minute)

	select {
	case <-ch:
		t.Error("timer fired too early")
	default:
		// ok
	}

	mc.Advance(31 * time.Minute)

	select {
	case v := <-ch:
		expected := start.Add(1 * time.Hour)
		if !v.Equal(expected) {
			t.Errorf("expected timer at %v, got %v", expected, v)
		}
	default:
		t.Error("timer should have fired")
	}
}

func TestMockClock_Sleep(t *testing.T) {
	start := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	mc := NewMockClock(start)

	mc.Sleep(3 * time.Second)
	expected := start.Add(3 * time.Second)
	if mc.Now() != expected {
		t.Errorf("expected %v after Sleep, got %v", expected, mc.Now())
	}
}
