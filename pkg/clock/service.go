package clock

import "context"

type Service struct {
	Clock
}

// Run function for compatibility with the supervisor-go Service interface.
func (receiver *Service) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}
