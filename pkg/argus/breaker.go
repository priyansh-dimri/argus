package argus

import (
	"time"

	"github.com/sony/gobreaker/v2"
)

type Breaker struct {
	cb *gobreaker.CircuitBreaker[any]
}

func NewBreaker(name string) *Breaker {
	settings := gobreaker.Settings{
		Name:        name,
		MaxRequests: 1,                // Allowing 1 request to pass in Half Open state
		Interval:    60 * time.Second, // Clear failure counts every 60s if it is not tripped again
		Timeout:     30 * time.Second, // Half Open for 30s before retrying again

		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 3
		},
	}

	return &Breaker{
		cb: gobreaker.NewCircuitBreaker[any](settings),
	}
}

func (b *Breaker) Execute(req func() (any, error)) (any, error) {
	return b.cb.Execute(req)
}
