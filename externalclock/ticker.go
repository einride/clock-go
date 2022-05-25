package externalclock

import (
	"fmt"
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
	t.duration = duration
	t.SetLastTimestamp(now)
}

func (t *ticker) IsDurationReached(currentTime time.Time) bool {
	return t.duration <= currentTime.Sub(t.GetLastTimestamp())
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
	g.Logger.V(1).Info("added new ticker", "caller", caller)
	return g.newTickerInternal(caller, nil, d, true)
}

func (g *Clock) newTickerInternal(caller string, endFunc func(), d time.Duration, periodic bool) clock.Ticker {
	c := make(chan time.Time)
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
