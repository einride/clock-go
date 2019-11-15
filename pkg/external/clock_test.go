package external

import (
	"context"
	"testing"
	"time"

	perceptionv1 "github.com/einride/proto/gen/go/perception/v1"
	"github.com/golang/protobuf/ptypes"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func TestExternalClock_NewTicker(t *testing.T) {
	externalClock := newTestFixture(t)
	// Given a ticker with a tick time
	tickTime := 100 * time.Millisecond

	// Then add a ticker
	loopTicker := externalClock.NewTicker(tickTime)
	defer loopTicker.Stop()

	// should find a new one.
	require.Len(t, externalClock.tickers, 1)
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
	nanoTime, _ := ptypes.TimestampProto(time.Unix(0, tickTime.Nanoseconds()+1))
	externalClock.SetEgoState(
		&perceptionv1.EgoState{
			Time: nanoTime,
		})
	didSet := false
	select {
	case <-time.After(1 * time.Millisecond):
	case <-loopTicks:
		didSet = true
	}
	require.False(t, didSet)
}

func TestExternalClock_RemoveTicker(t *testing.T) {
	externalClock := newTestFixture(t)

	// Given a ticker with a tick time
	tickTime := 100 * time.Millisecond
	loopTicker := externalClock.NewTicker(tickTime)

	// then stopping it
	loopTicker.Stop()

	// should be empty
	require.Len(t, externalClock.tickers, 0, "should be empty %+v", externalClock.tickers)
}

func TestExternalClock_Tick(t *testing.T) {
	externalClock := newTestFixture(t)
	// Given a ticker with a tick time
	tickTime := 100 * time.Millisecond
	loopTicker := externalClock.NewTicker(tickTime)
	loopTicks := loopTicker.C()

	// Send a tick
	tickTimeProto, _ := ptypes.TimestampProto(time.Unix(0, tickTime.Nanoseconds()))
	externalClock.SetEgoState(
		&perceptionv1.EgoState{
			Time: tickTimeProto,
		})

	// exect didSet to be true
	didSet := false
	select {
	case <-time.After(1 * time.Millisecond):
		t.FailNow()
	case <-loopTicks:
		didSet = true
	}
	require.True(t, didSet)
}

func TestExternalClock_After(t *testing.T) {
	externalClock := newTestFixture(t)

	// Given a ticker with a tick time
	tickTime := 1 * time.Millisecond
	afterChan := externalClock.After(tickTime)

	// Send a tick
	tickTimeProto, _ := ptypes.TimestampProto(time.Unix(0, tickTime.Nanoseconds()+1))
	externalClock.SetEgoState(
		&perceptionv1.EgoState{
			Time: tickTimeProto,
		})

	// exect didSet to be true
	didSet := false
	select {
	case <-time.After(2 * time.Millisecond):
		t.FailNow()
	case <-afterChan:
		didSet = true
	}
	require.True(t, didSet)
}

func TestExternalClock_Removed(t *testing.T) {
	externalClock := newTestFixture(t)

	// Given a ticker with a tick time
	tickTime := 1 * time.Millisecond
	afterChan := externalClock.After(tickTime)

	// Send a tick
	tickTimeProto, _ := ptypes.TimestampProto(time.Unix(0, tickTime.Nanoseconds()))
	externalClock.SetEgoState(
		&perceptionv1.EgoState{
			Time: tickTimeProto,
		})

	// exect didSet to be true
	select {
	case <-time.After(2 * time.Millisecond):
		t.FailNow()
	case <-afterChan:
		require.Len(t, externalClock.tickers, 0)
	}
}

func TestExternalClock_AfterFailToShortTime(t *testing.T) {
	externalClock := newTestFixture(t)

	// Given a ticker with a tick time
	tickTime := 500 * time.Microsecond
	afterChan := externalClock.After(tickTime)

	// Send a tick
	tickTimeProto, _ := ptypes.TimestampProto(time.Unix(0, 0))
	externalClock.SetEgoState(
		&perceptionv1.EgoState{
			Time: tickTimeProto,
		})

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
		tickTimeProto, _ := ptypes.TimestampProto(time.Unix(0, int64(i+1)*tickTime.Nanoseconds()+1))
		externalClock.SetEgoState(
			&perceptionv1.EgoState{
				Time: tickTimeProto,
			})

		// exect didSet to be true
		didSet := false
		select {
		case <-time.After(1 * time.Second):
			t.FailNow()
		case <-loopTicks:
			didSet = true
		}
		require.True(t, didSet)
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
				tickProto, _ := ptypes.TimestampProto(time.Unix(0, (tick+1)*time.Millisecond.Nanoseconds()))
				externalClock.SetEgoState(
					&perceptionv1.EgoState{
						Time: tickProto,
					})
			}
		},
	}
	var g errgroup.Group
	g.Go(func() error {
		return looper.Run(ctx)
	})
	<-looper.init
	tickProto, _ := ptypes.TimestampProto(time.Unix(0, 1*time.Millisecond.Nanoseconds()))
	externalClock.SetEgoState(
		&perceptionv1.EgoState{
			Time: tickProto,
		})
	err := g.Wait()
	require.NoError(t, err)
	require.Equal(t, looper.Target, looper.counter)
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
	tickProto, _ := ptypes.TimestampProto(time.Unix(0, 1*time.Millisecond.Nanoseconds()+1))
	externalClock.SetEgoState(
		&perceptionv1.EgoState{
			Time: tickProto,
		})
	err := g.Wait()
	require.NoError(t, err)
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
