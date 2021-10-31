package contextctrl

import (
	"context"
	"time"
)

type SystemContext struct{}

func (s *SystemContext) WithTimeout(
	parent context.Context,
	timeout time.Duration,
) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}
