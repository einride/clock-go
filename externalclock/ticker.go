package externalclock

import (
	"fmt"
	"log/slog"
	"runtime"
	"sync"
	"time"

	"go.einride.tech/clock"
)

type ticker struct {
	mutex         sync.Mutex
	caller        string
	lastTimeStamp time.Time
	duration      time.Duration
	timeChan      chan time.Time
	stopFunc      func()
	isPeriodic    bool
	getTimeFunc   func() time.Time
}

func (t *ticker) C() <-chan time.Time {
	return t.timeChan
}

func (t *ticker) Stop() {
	t.stopFunc()
}

func (t *ticker) Reset(duration time.Duration) {
	now := t.getTimeFunc()
	t.mutex.Lock()
	t.duration = duration
	t.lastTimeStamp = now
	t.mutex.Unlock()
}

func (t *ticker) IsDurationReached(currentTime time.Time) bool {
	t.mutex.Lock()
	dur := t.duration
	ts := t.lastTimeStamp
	t.mutex.Unlock()
	return dur <= currentTime.Sub(ts)
}

func (t *ticker) GetLastTimestamp() time.Time {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.lastTimeStamp
}

func (t *ticker) SetLastTimestamp(lastTimestamp time.Time) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.lastTimeStamp = lastTimestamp
}

func (g *Clock) NewTicker(d time.Duration) clock.Ticker {
	_, file, no, ok := runtime.Caller(1)
	var caller string
	if ok {
		caller = fmt.Sprintf("called from %s#%d\n", file, no)
	}
	slog.Debug("added new ticker", slog.String("caller", caller))
	return g.newTickerInternal(caller, nil, d, true)
}

func (g *Clock) newTickerInternal(caller string, endFunc func(), d time.Duration, periodic bool) clock.Ticker {
	// Give the channel a 1-element time buffer.
	// If the client falls behind while reading, we drop ticks
	// on the floor until the client catches up.
	c := make(chan time.Time, 1)
	uuid := makeUUID()
	intervalTicker := &ticker{
		caller:   caller,
		timeChan: c,
		duration: d,
		stopFunc: func() {
			g.tickerMutex.Lock()
			delete(g.tickers, uuid)
			g.tickerMutex.Unlock()
			if endFunc != nil {
				endFunc()
			}
		},
		isPeriodic:  periodic,
		getTimeFunc: g.getTime,
	}
	intervalTicker.SetLastTimestamp(g.getTime())
	g.tickerMutex.Lock()
	g.tickers[uuid] = intervalTicker
	g.tickerMutex.Unlock()
	return intervalTicker
}
