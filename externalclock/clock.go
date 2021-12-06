package externalclock

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"go.einride.tech/clock"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Clock struct {
	Logger        *zap.Logger
	timestampChan chan time.Time
	timeMutex     sync.Mutex
	currentTime   time.Time
	tickerMutex   sync.RWMutex
	tickers       map[string]*ticker
}

func New(logger *zap.Logger, initialTime time.Time) *Clock {
	c := &Clock{
		Logger:        logger,
		timestampChan: make(chan time.Time),
		tickers:       map[string]*ticker{},
	}
	c.currentTime = initialTime
	return c
}

func (g *Clock) SetTimestamp(t time.Time) {
	g.timestampChan <- t
}

func (g *Clock) NumberOfTriggers() int {
	return len(g.tickers)
}

func (g *Clock) Run(ctx context.Context) error {
	ctxDone := ctx.Done()
	g.Logger.Info("clock started")
	for {
		select {
		case <-ctxDone:
			return nil
		case recvTime := <-g.timestampChan:
			g.timeMutex.Lock()
			g.tickerMutex.RLock()
			for _, tickerInstance := range g.tickers {
				if !tickerInstance.IsDurationReached(recvTime) {
					continue
				}
				tickerInstance.SetLastTimestamp(recvTime)
				if !tickerInstance.isPeriodic {
					g.tickerMutex.RUnlock()
					tickerInstance.Stop()
					g.tickerMutex.RLock()
				}
				select {
				case tickerInstance.timeChan <- recvTime:
				case <-time.After(20 * time.Millisecond):
					g.Logger.Warn("ticker dropped message", zap.String("called from", tickerInstance.caller))
				}
			}
			g.setTime(recvTime)
			g.timeMutex.Unlock()
			g.tickerMutex.RUnlock()
		}
	}
}

func (g *Clock) getTime() time.Time {
	g.timeMutex.Lock()
	defer g.timeMutex.Unlock()
	return g.currentTime
}

func (g *Clock) setTime(newTime time.Time) {
	// note: mutex should be locked when arriving here.
	g.currentTime = newTime
}

func (g *Clock) After(duration time.Duration) <-chan time.Time {
	_, file, no, ok := runtime.Caller(1)
	var calledFrom string
	if ok {
		calledFrom = fmt.Sprintf("called from %s#%d\n", file, no)
	}
	tickerInstance := g.newTickerInternal(calledFrom, nil, duration, false)
	return tickerInstance.C()
}

func (g *Clock) AfterFunc(d time.Duration, f func()) clock.Timer {
	_, file, no, ok := runtime.Caller(1)
	var calledFrom string
	if ok {
		calledFrom = fmt.Sprintf("called from %s#%d\n", file, no)
	}
	return &Timer{
		Ticker: g.newTickerInternal(calledFrom, f, d, false),
	}
}

func (g *Clock) Now() time.Time {
	return g.getTime()
}

func (g *Clock) NowProto() *timestamppb.Timestamp {
	return timestamppb.New(g.getTime())
}

func (g *Clock) Since(t time.Time) time.Duration {
	now := g.getTime()
	return now.Sub(t)
}

func (g *Clock) Sleep(d time.Duration) {
	<-g.After(d)
}

func makeUUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}
