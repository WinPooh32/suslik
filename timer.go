package suslik

import "time"

type Timer struct {
	Begin          time.Time
	InactiveBegin  time.Time
	ActiveDuration time.Duration
	Run            bool
}

func MakeTimer() Timer {
	return Timer{}
}

func (t *Timer) Start() {
	t.ActiveDuration += t.InactiveBegin.Sub(t.Begin)
	t.Begin = time.Now()
	t.Run = true
}

func (t *Timer) Stop() {
	t.InactiveBegin = time.Now()
	t.Run = false
}

func (t *Timer) Reset() {
	t.ActiveDuration = 0
	t.Begin = time.Now()
	t.InactiveBegin = t.Begin
}

func (t *Timer) Elapsed() time.Duration {
	if t.Run {
		return time.Since(t.Begin) + t.ActiveDuration
	}
	return t.InactiveBegin.Sub(t.Begin) + t.ActiveDuration
}
