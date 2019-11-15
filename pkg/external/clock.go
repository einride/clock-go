package external

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/einride/clock-go/pkg/clock"
	"github.com/einride/clock-go/pkg/external/ticker"
	perceptionv1 "github.com/einride/proto/gen/go/perception/v1"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"go.uber.org/zap"
)

type Clock struct {
	Logger        *zap.Logger
	timestampChan chan time.Time
	timeMutex     sync.RWMutex
	currentTime   time.Time
	tickerMutex   sync.RWMutex
	tickers       map[string]*ticker.Ticker
}

func NewClock(logger *zap.Logger) *Clock {
	c := &Clock{
		Logger:        logger,
		timestampChan: make(chan time.Time),
		tickers:       map[string]*ticker.Ticker{},
	}
	c.currentTime = time.Unix(0, 0)
	return c
}

func (g *Clock) SetEgoState(egoState *perceptionv1.EgoState) {
	g.timestampChan <- time.Unix(egoState.Time.Seconds, int64(egoState.Time.Nanos)).UTC()
}

func (g *Clock) Run(ctx context.Context) error {
	ctxDone := ctx.Done()
	g.Logger.Info("clock started")
	for {
		select {
		case <-ctxDone:
			return nil
		case recvTime := <-g.timestampChan:
			g.tickerMutex.RLock()
			for _, tickerInstance := range g.tickers {
				if !tickerInstance.IsDurationReached(recvTime) {
					continue
				}
				tickerInstance.SetLastTimestamp(recvTime)
				if !tickerInstance.IsPeriodic {
					g.tickerMutex.RUnlock()
					tickerInstance.Stop()
					g.tickerMutex.RLock()
				}
				select {
				case tickerInstance.TimeChan <- recvTime:
				case <-time.After(20 * time.Millisecond):
					g.Logger.Warn("ticker dropped message", zap.String("called from", tickerInstance.Caller))
				}
			}
			g.tickerMutex.RUnlock()
			g.setTime(recvTime)
		}
	}
}

func (g *Clock) getTime() time.Time {
	g.timeMutex.RLock()
	defer g.timeMutex.RUnlock()
	return g.currentTime
}

func (g *Clock) setTime(newTime time.Time) {
	g.timeMutex.Lock()
	defer g.timeMutex.Unlock()
	g.currentTime = newTime
}

func (g *Clock) After(duration time.Duration) <-chan time.Time {
	_, file, no, ok := runtime.Caller(1)
	var calledFrom string
	if ok {
		calledFrom = fmt.Sprintf("called from %s#%d\n", file, no)
	}
	tickerInstance := g.newTickerInternal(calledFrom, duration, false)
	return tickerInstance.C()
}

func (g *Clock) NewTicker(d time.Duration) clock.Ticker {
	_, file, no, ok := runtime.Caller(1)
	var calledFrom string
	if ok {
		calledFrom = fmt.Sprintf("called from %s#%d\n", file, no)
	}
	g.Logger.Info("added new ticker", zap.String("called from", calledFrom))
	return g.newTickerInternal(calledFrom, d, true)
}

func (g *Clock) newTickerInternal(caller string, d time.Duration, periodic bool) clock.Ticker {
	c := make(chan time.Time)
	uuid := makeUUID()
	intervalTicker := &ticker.Ticker{
		Caller:   caller,
		TimeChan: c,
		Duration: d,
		StopFunc: func() {
			g.tickerMutex.Lock()
			delete(g.tickers, uuid)
			g.tickerMutex.Unlock()
		},
		IsPeriodic: periodic,
	}
	intervalTicker.SetLastTimestamp(g.getTime())
	g.tickerMutex.Lock()
	g.tickers[uuid] = intervalTicker
	g.tickerMutex.Unlock()
	return intervalTicker
}

func (g *Clock) Now() time.Time {
	return g.getTime()
}

func (g *Clock) NowProto() *timestamp.Timestamp {
	t, err := ptypes.TimestampProto(g.getTime())
	if err != nil {
		g.Logger.Error("Convert time from external", zap.Error(err))
	}
	return t
}

func (g *Clock) Since(t time.Time) time.Duration {
	now := g.getTime()
	return now.Sub(t)
}

func makeUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}
