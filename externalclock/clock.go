package externalclock

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/go-logr/logr"
	"go.einride.tech/clock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Clock struct {
	Logger      logr.Logger
	timeMutex   sync.Mutex
	currentTime time.Time
	tickerMutex sync.RWMutex
	tickers     map[string]*ticker
}

func New(logger logr.Logger, initialTime time.Time) *Clock {
	c := &Clock{
		Logger:  logger,
		tickers: map[string]*ticker{},
	}
	c.currentTime = initialTime
	return c
}

func (g *Clock) SetTimestamp(t time.Time) {
	g.timeMutex.Lock()
	g.currentTime = t
	g.timeMutex.Unlock()

	g.signalTickers(t)
}

func (g *Clock) signalTickers(t time.Time) {
	g.tickerMutex.RLock()
	for _, tickerInstance := range g.tickers {
		if !tickerInstance.IsDurationReached(t) {
			continue
		}
		tickerInstance.SetLastTimestamp(t)
		if !tickerInstance.isPeriodic {
			g.tickerMutex.RUnlock()
			tickerInstance.Stop()
			g.tickerMutex.RLock()
		}
		select {
		case tickerInstance.timeChan <- t:
		case <-time.After(20 * time.Millisecond):
			g.Logger.V(1).Info("ticker dropped message", "caller", tickerInstance.caller)
		}
	}
	g.tickerMutex.RUnlock()
}

func (g *Clock) NumberOfTriggers() int {
	g.tickerMutex.RLock()
	defer g.tickerMutex.RUnlock()
	return len(g.tickers)
}

// Deprecated: Calling Run is not necessary anymore. The method only blocks until
// context is cancelled.
func (g *Clock) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (g *Clock) getTime() time.Time {
	g.timeMutex.Lock()
	defer g.timeMutex.Unlock()
	return g.currentTime
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
