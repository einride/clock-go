// Package clock provides primitives for mocking time.
package clock

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// Clock provides capabilities from the time standard library package.
type Clock interface {
	// After waits for the duration to elapse and then sends the current time on the returned channel.
	After(duration time.Duration) <-chan time.Time

	// AfterFunc waits for the duration to elapse and then calls the function f given to it.
	AfterFunc(d time.Duration, f func()) Timer

	// NewTicker returns a new Ticker.
	NewTicker(d time.Duration) Ticker

	// Now returns the current local time.
	Now() time.Time

	// NowProto returns a new Protobuf timestamp representing the current local time.
	NowProto() *timestamppb.Timestamp

	// Since returns the time elapsed since t.
	Since(t time.Time) time.Duration

	// Sleep pauses the current goroutine for at least the duration d.
	// A negative or zero duration causes Sleep to return immediately.
	Sleep(d time.Duration)
}

// Ticker wraps the time.Ticker class.
type Ticker interface {
	// C returns the channel on which the ticks are delivered.
	C() <-chan time.Time

	// Stop the Ticker.
	Stop()
}

type Timer interface {
	// C returns the channel on which the timer is going to be triggered.
	C() <-chan time.Time

	// Stop the Timer.
	Stop() bool
}
