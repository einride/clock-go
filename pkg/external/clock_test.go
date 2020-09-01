package external

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestExternalClock_NewTicker(t *testing.T) {
	externalClock := newTestFixture(t)
	// Given a ticker with a tick time
	tickTime := 100 * time.Millisecond

	// Then add a ticker
	loopTicker := externalClock.NewTicker(tickTime)
	defer loopTicker.Stop()

	// should find a new one.
	assert.Assert(t, is.Len(externalClock.tickers, 1))
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
	assert.Assert(t, is.Len(externalClock.tickers, 0), "should be empty %+v", externalClock.tickers)
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
		assert.Assert(t, is.Len(externalClock.tickers, 0))
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
	Clock        *Clock
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

func newTestFixture(t *testing.T) *Clock {
	t.Helper()
	c := NewClock(zap.NewExample())
	var g errgroup.Group
	g.Go(func() error {
		return c.Run(context.Background())
	})
	return c
}
