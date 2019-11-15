package ticker

import (
	"sync"
	"time"
)

type Ticker struct {
	Mutex         sync.Mutex
	Caller        string
	lastTimeStamp time.Time
	Duration      time.Duration
	TimeChan      chan time.Time
	StopFunc      func()
	IsPeriodic    bool
}

func (t *Ticker) C() <-chan time.Time {
	return t.TimeChan
}

func (t *Ticker) Stop() {
	t.StopFunc()
}

func (t *Ticker) IsDurationReached(currentTime time.Time) bool {
	return t.Duration <= currentTime.Sub(t.GetLastTimestamp())
}

func (t *Ticker) GetLastTimestamp() time.Time {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	return t.lastTimeStamp
}

func (t *Ticker) SetLastTimestamp(lastTimestamp time.Time) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	t.lastTimeStamp = lastTimestamp
}
