package clock

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// System returns a Clock implementation that delegate to the time package.
func System() Clock {
	return &systemClock{}
}

type systemClock struct{}

var _ Clock = &systemClock{}

func (c systemClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func (c systemClock) NowProto() *timestamppb.Timestamp {
	return timestamppb.Now()
}

func (c systemClock) NewTicker(d time.Duration) Ticker {
	return &systemTicker{Ticker: *time.NewTicker(d)}
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

func (c systemClock) AfterFunc(d time.Duration, f func()) Timer {
	return &systemTimer{Timer: time.AfterFunc(d, f)}
}

type systemTicker struct {
	time.Ticker
}

func (t systemTicker) C() <-chan time.Time {
	return t.Ticker.C
}

type systemTimer struct {
	*time.Timer
}

func (t systemTimer) C() <-chan time.Time {
	return t.Timer.C
}
