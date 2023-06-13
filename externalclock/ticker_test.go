package externalclock_test

import (
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestExternalClock_TickerReset(t *testing.T) {
	externalClock := newTestFixture(t)
	// Given a ticker ticking every 3rd ms
	externalTime := time.Unix(0, 0)
	count := 10
	delta := time.Millisecond
	tickerDelta := 3 * time.Millisecond

	loopTicker := externalClock.NewTicker(tickerDelta)
	cancel := make(chan struct{})
	receivedTime := make([]time.Time, 0, 2)
	go func() {
		for {
			select {
			case <-cancel:
				return
			case ts := <-loopTicker.C():
				receivedTime = append(receivedTime, ts)
			}
		}
	}()
	// we let the clock go up to 5, giving 1 tick at 3.
	for i := 0; i < count/2; i++ {
		externalTime = externalTime.Add(delta)
		externalClock.SetTimestamp(externalTime)
		// add sleep to let the channels clear between the calls.
		time.Sleep(time.Millisecond)
	}
	// reset the ticker just before the next tick at 6, but reset so next tick will be at 5 + 3 => 8
	assert.Equal(t, time.UnixMilli(5), externalTime)
	loopTicker.Reset(tickerDelta)
	for i := count / 2; i < count; i++ {
		externalTime = externalTime.Add(delta)
		externalClock.SetTimestamp(externalTime)
		time.Sleep(time.Millisecond)
	}
	cancel <- struct{}{}
	assert.DeepEqual(t, receivedTime, []time.Time{time.UnixMilli(3), time.UnixMilli(8)})
}
