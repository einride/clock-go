package externalclock_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	"go.einride.tech/clock/externalclock"
	"golang.org/x/sync/errgroup"
	"gotest.tools/v3/assert"
)

func TestExternalClock_NewTicker(t *testing.T) {
	externalClock := newTestFixture(t)
	// Given a ticker with a tick time
	tickTime := 100 * time.Millisecond

	// Then add a ticker
	loopTicker := externalClock.NewTicker(tickTime)
	defer loopTicker.Stop()

	// should find a new one.
	assert.Equal(t, externalClock.NumberOfTriggers(), 1)
}

func TestExternalClock_Now(t *testing.T) {
	// TestExternalClock_Now verifies that calling Now() method straight after a SetTimestamp(ts)
	// will return ts as expected.

	// Given
	externalClock := newTestFixture(t)
	externalClock.SetTimestamp(time.Unix(0, 0))

	// make sure timestamp is setup before we continue
	<-time.After(100 * time.Millisecond) // timeout
	if externalClock.Now() != time.Unix(0, 0) {
		t.Fatalf("Could not initialize time before timeout.")
	}

	// Feed the clock with a few different timestamps
	t0 := time.Unix(10, 0)
	timeList := []time.Time{t0, t0.Add(time.Minute), t0.Add(2 * time.Minute)}

	for i := 0; i < 10; i++ { // repeat a few times to try to trigger racing
		for _, ts := range timeList {
			// Given
			externalClock.SetTimestamp(ts)
			// Expect
			cnow := externalClock.Now()
			if ts.UnixNano() != cnow.UnixNano() {
				t.Errorf("Expected: %v, Got: %v", ts.UnixNano(), cnow.UnixNano())
			}
		}
	}
}

func TestExternalClock_Stop(t *testing.T) {
	externalClock := newTestFixture(t)

	// Given a ticker with a tick time
	tickTime := 100 * time.Millisecond
	loopTicker := externalClock.NewTicker(tickTime)
	// then stopping it
	loopTicker.Stop()

	// should not be able to tick
	loopTicks := loopTicker.C()
	externalClock.SetTimestamp(time.Unix(0, tickTime.Nanoseconds()+1))
	didSet := false
	select {
	case <-time.After(1 * time.Millisecond):
	case <-loopTicks:
		didSet = true
	}
	assert.Assert(t, !didSet)
}

func TestExternalClock_RemoveTicker(t *testing.T) {
	externalClock := newTestFixture(t)

	// Given a ticker with a tick time
	tickTime := 100 * time.Millisecond
	loopTicker := externalClock.NewTicker(tickTime)

	// then stopping it
	loopTicker.Stop()

	// should be empty
	assert.Equal(t, externalClock.NumberOfTriggers(), 0, "should be empty %+v", externalClock.NumberOfTriggers())
}

func TestExternalClock_Tick(t *testing.T) {
	externalClock := newTestFixture(t)
	// Given a ticker with a tick time
	tickTime := 100 * time.Millisecond
	loopTicker := externalClock.NewTicker(tickTime)
	loopTicks := loopTicker.C()

	// Send a tick
	externalClock.SetTimestamp(time.Unix(0, tickTime.Nanoseconds()))

	// exect didSet to be true
	didSet := false
	select {
	case <-time.After(1 * time.Millisecond):
		t.FailNow()
	case <-loopTicks:
		didSet = true
	}
	assert.Assert(t, didSet)
}

func TestExternalClock_After(t *testing.T) {
	externalClock := newTestFixture(t)

	// Given a ticker with a tick time
	tickTime := 1 * time.Millisecond
	afterChan := externalClock.After(tickTime)

	// Send a tick
	externalClock.SetTimestamp(time.Unix(0, tickTime.Nanoseconds()+1))

	// exect didSet to be true
	didSet := false
	select {
	case <-time.After(2 * time.Millisecond):
		t.FailNow()
	case <-afterChan:
		didSet = true
	}
	assert.Assert(t, didSet)
}

func TestExternalClock_AfterFunc(t *testing.T) {
	externalClock := newTestFixture(t)

	// Given an external clock and a trigger time
	triggerTime := 1 * time.Millisecond

	// then trigger func after trigger time
	didSet := false
	afterTimer := externalClock.AfterFunc(triggerTime, func() {
		didSet = true
	})
	externalClock.SetTimestamp(time.Unix(0, triggerTime.Nanoseconds()+1))
	select {
	case <-time.After(2 * time.Millisecond):
		t.FailNow()
	case <-afterTimer.C():
	}
	assert.Assert(t, didSet)
}

func TestExternalClock_Removed(t *testing.T) {
	externalClock := newTestFixture(t)

	// Given a ticker with a tick time
	tickTime := 1 * time.Millisecond
	afterChan := externalClock.After(tickTime)

	// Send a tick
	externalClock.SetTimestamp(time.Unix(0, tickTime.Nanoseconds()))

	// exect didSet to be true
	select {
	case <-time.After(2 * time.Millisecond):
		t.FailNow()
	case <-afterChan:
		assert.Equal(t, externalClock.NumberOfTriggers(), 0)
	}
}

func TestExternalClock_AfterFailToShortTime(t *testing.T) {
	externalClock := newTestFixture(t)

	// Given a ticker with a tick time
	tickTime := 500 * time.Microsecond
	afterChan := externalClock.After(tickTime)

	// Send a tick
	externalClock.SetTimestamp(time.Unix(0, 0))

	// exect didSet to be true
	select {
	case <-time.After(1 * time.Millisecond):
	case <-afterChan:
		t.FailNow()
	}
}

func TestExternalClock_NewTimer(t *testing.T) {
	externalClock := newTestFixture(t)

	// given duration for Timer
	timerTime := 1 * time.Millisecond
	timer := externalClock.NewTimer(timerTime)
	timerSignal := timer.C()

	externalClock.SetTimestamp(time.Unix(0, int64(1)*timerTime.Nanoseconds()+1))

	select {
	case <-time.After(1 * time.Second):
		t.FailNow()
	case <-timerSignal:
	}
}

func TestExternalClock_NewTicker_Tick_Periodically(t *testing.T) {
	externalClock := newTestFixture(t)
	// Given a ticker with a tick time
	tickTime := 1 * time.Millisecond
	loopTicker := externalClock.NewTicker(tickTime)
	loopTicks := loopTicker.C()

	// Send a tick
	for i := range make([]int64, 1000) {
		externalClock.SetTimestamp(time.Unix(0, int64(i+1)*tickTime.Nanoseconds()+1))

		// exect didSet to be true
		didSet := false
		select {
		case <-time.After(1 * time.Second):
			t.FailNow()
		case <-loopTicks:
			didSet = true
		}
		assert.Assert(t, didSet)
	}
}

func TestExternalClock_SendBeforeRun(t *testing.T) {
	// test verifies that sending time on an unstarted clock does not deadlock
	c := externalclock.New(testr.New(t), time.Unix(0, 0))
	c.SetTimestamp(time.Unix(1, 0))
}

func TestExternalClock_SendAfterRun(t *testing.T) {
	// test verifies that sending time on a cancelled clock does not deadlock
	c := externalclock.New(testr.New(t), time.Unix(0, 0))
	// start clock with a deadline
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	t.Cleanup(cancel)
	assert.NilError(t, c.Run(ctx))

	// sending time
	c.SetTimestamp(time.Unix(1, 0))
}

func TestExternalClock_TestLooper(t *testing.T) {
	externalClock := newTestFixture(t)
	const target = 1000
	ctx, cancel := context.WithCancel(context.Background())
	looper := testLooper{
		init:         make(chan struct{}),
		Target:       target,
		Clock:        externalClock,
		LoopInterval: 1 * time.Millisecond,
		Callback: func(tick int64) {
			if tick == target {
				cancel()
			} else {
				externalClock.SetTimestamp(time.Unix(0, (tick+1)*time.Millisecond.Nanoseconds()))
			}
		},
	}
	var g errgroup.Group
	g.Go(func() error {
		return looper.Run(ctx)
	})
	<-looper.init
	externalClock.SetTimestamp(time.Unix(0, 1*time.Millisecond.Nanoseconds()))
	assert.NilError(t, g.Wait())
	assert.Equal(t, looper.Target, looper.counter)
}

func TestTicker_ResetRace(t *testing.T) {
	// This test does not have any assertions, and instead is intended to catch
	// data races with `-race` flag set.
	const (
		ticks        = 100
		tickInterval = time.Millisecond
	)
	clock := newTestFixture(t)
	done := make(chan struct{}, 2)
	go func() {
		initialTime := time.Unix(0, 0)
		for i := 0; i < ticks; i++ {
			clock.SetTimestamp(initialTime.Add(time.Second))
			time.Sleep(tickInterval)
		}
		done <- struct{}{}
	}()
	go func() {
		ticker := clock.NewTicker(time.Second)
		for i := 0; i < ticks; i++ {
			ticker.Reset(time.Second)
			time.Sleep(tickInterval)
		}
		ticker.Stop()
		done <- struct{}{}
	}()

	<-done
	<-done
}

func TestExternalClock_TestLooper_AddTicker(t *testing.T) {
	externalClock := newTestFixture(t)
	ctx, cancel := context.WithCancel(context.Background())
	looper := testLooper{
		init:         make(chan struct{}),
		Clock:        externalClock,
		LoopInterval: 1 * time.Millisecond,
		Callback: func(tick int64) {
			cancel()
		},
	}
	var g errgroup.Group
	g.Go(func() error {
		return looper.Run(ctx)
	})
	<-looper.init
	externalClock.SetTimestamp(time.Unix(0, 1*time.Millisecond.Nanoseconds()+1))
	assert.NilError(t, g.Wait())
}

type testLooper struct {
	init         chan struct{}
	Clock        *externalclock.Clock
	LoopInterval time.Duration
	Target       int64
	counter      int64
	Callback     func(tick int64)
}

func (t *testLooper) Run(ctx context.Context) error {
	ctxDone := ctx.Done()
	loopTicker := t.Clock.NewTicker(t.LoopInterval)
	defer loopTicker.Stop()
	loopTicks := loopTicker.C()
	close(t.init)
	for {
		select {
		case <-ctxDone:
			return nil
		case <-loopTicks:
			t.counter++
			t.Callback(t.counter)
		}
	}
}

func newTestFixture(t *testing.T) *externalclock.Clock {
	t.Helper()
	c := externalclock.New(testr.New(t), time.Unix(0, 0))
	var g errgroup.Group
	g.Go(func() error {
		if err := c.Run(context.Background()); err != nil {
			return fmt.Errorf("new fixture: %w", err)
		}
		return nil
	})
	return c
}
