package external

import (
	"time"

	"github.com/einride/clock-go/pkg/clock"
)

// The Timer type represents a single event.
// When the Timer expires, the current time will be sent on C,
// unless the Timer was created by AfterFunc.
// A Timer must be created with NewTimer or AfterFunc.
type Timer struct {
	clock.Ticker
}

// NewTimer creates a new Timer that will send
// the current time on its channel after at least duration d.
func (g *Clock) NewTimer(d time.Duration) *Timer {
	return &Timer{
		Ticker: g.newTickerInternal("not set", nil, d, false),
	}
}

// Stop prevents the Timer from firing.
// It returns true if the call stops the timer, false if the timer has already
// expired or been stopped.
func (t *Timer) Stop() bool {
	t.Ticker.Stop()
	return true
}
