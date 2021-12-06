package suslik

import (
	"time"
)

type Animation struct {
	TotalFrames int
	Current     int
	Delay       time.Duration
	Playing     bool
	Timer       Timer
}

func NewAnimation(frames int, delay time.Duration) *Animation {
	return &Animation{
		TotalFrames: frames,
		Current:     0,
		Delay:       delay,
		Playing:     false,
		Timer:       MakeTimer(),
	}
}

func (anim *Animation) Play(num int) {
	anim.Timer.Start()
	anim.Current = num
	anim.Playing = true
}

func (anim *Animation) Stop() {
	anim.Timer.Stop()
	anim.Playing = false
}

func (anim *Animation) NextFrame() int {
	if anim.Playing && anim.Timer.Elapsed() > anim.Delay {
		anim.Timer.Reset()
		anim.Current = (anim.Current + 1) % anim.TotalFrames
	}
	return anim.Current
}

func (anim *Animation) GetFrame() int {
	return anim.Current
}

func (anim *Animation) SetFrame(n int) {
	if n > 0 && n < anim.TotalFrames {
		anim.Current = n
	}
}

func (anim *Animation) ResetTime() {
	anim.Timer.Reset()
}

func (anim *Animation) GetDelay() time.Duration {
	return anim.Delay
}

func (anim *Animation) SetDelay(d time.Duration) {
	anim.Delay = d
}

func (anim *Animation) GetTotalFrames() int {
	return anim.TotalFrames
}

func (anim *Animation) SetTotalFrames(total int) {
	anim.TotalFrames = total
}
