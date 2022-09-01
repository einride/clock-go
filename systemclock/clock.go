package systemclock

import (
	"time"

	"go.einride.tech/clock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// New returns a clock.Clock implementation that delegates to the time package.
func New() clock.Clock {
	return &systemClock{}
}

type systemClock struct{}

var _ clock.Clock = &systemClock{}

func (c systemClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func (c systemClock) NowProto() *timestamppb.Timestamp {
	return timestamppb.Now()
}

func (c systemClock) NewTicker(d time.Duration) clock.Ticker {
	return &systemTicker{ticker: time.NewTicker(d)}
}

func (c systemClock) Now() time.Time {
	return time.Now()
}

func (c systemClock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

func (c systemClock) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (c systemClock) AfterFunc(d time.Duration, f func()) clock.Timer {
	return &systemTimer{Timer: time.AfterFunc(d, f)}
}

type systemTicker struct {
	ticker *time.Ticker
}

func (t systemTicker) Stop() {
	t.ticker.Stop()
}

func (t systemTicker) Reset(duration time.Duration) {
	t.ticker.Reset(duration)
}

func (t systemTicker) C() <-chan time.Time {
	return t.ticker.C
}

type systemTimer struct {
	*time.Timer
}

func (t systemTimer) C() <-chan time.Time {
	return t.Timer.C
}
