package contextctrl_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/einride/clock-go/pkg/contextctrl"
	"github.com/einride/clock-go/pkg/external"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"gotest.tools/v3/assert"
)

func TestExternalClock_ContextWithTimeout(t *testing.T) {
	externalClock := newTestFixture(t)
	// given a external clock
	eContext := contextctrl.ExternalContext{C: externalClock}
	ctx := context.Background()

	// create a context with a timeout and tick it
	newCtx, _ := eContext.WithTimeout(ctx, time.Millisecond)
	externalClock.SetTimestamp(time.Unix(0, time.Millisecond.Nanoseconds()+1))
	d := newCtx.Done()
	select {
	case <-time.After(time.Millisecond):
		t.FailNow()
	case <-d:
	}

	// make sure there is an deadline error
	assert.Error(t, newCtx.Err(), "context deadline exceeded")
}

func TestExternalClock_ContextWithTimeoutCancel(t *testing.T) {
	externalClock := newTestFixture(t)
	// given a external clock
	eContext := contextctrl.ExternalContext{C: externalClock}

	// create a context with a timeout
	ctx := context.Background()
	newCtx, cancel := eContext.WithTimeout(ctx, time.Millisecond)

	// then cancel it
	cancel()
	d := newCtx.Done()
	select {
	case <-time.After(time.Millisecond):
		t.FailNow()
	case <-d:
	}

	// make sure there is a canceled error
	assert.Error(t, newCtx.Err(), "context canceled")
}

func TestExternalClock_ContextWithTimeoutSuccess(t *testing.T) {
	externalClock := newTestFixture(t)
	// given a external clock
	eContext := &contextctrl.ExternalContext{C: externalClock}

	// create a context with a timeout and tick it
	ctx := context.Background()
	newCtx, _ := eContext.WithTimeout(ctx, time.Millisecond)
	d := newCtx.Done()
	externalClock.SetTimestamp(time.Unix(0, time.Millisecond.Nanoseconds()-1))
	select {
	case <-d:
		t.FailNow()
	case <-time.After(time.Millisecond):
		// all is well
	}
	assert.NilError(t, newCtx.Err())
}

func newTestFixture(t *testing.T) *external.Clock {
	t.Helper()
	c := external.NewClock(zap.NewExample(), time.Unix(0, 0))
	var g errgroup.Group
	g.Go(func() error {
		if err := c.Run(context.Background()); err != nil {
			return fmt.Errorf("new fixture: %w", err)
		}
		return nil
	})
	return c
}
