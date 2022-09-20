package metrics

import (
	"context"
	"time"
)

type Measurable interface {
	Run(ctx context.Context)
	Measure() error
	ClearOldStat(hoursAgo int)
	GetIndicators(everyN int, durationM int) (interface{}, error)
}

type Metric struct {
	Measure func() error
}

func (m *Metric) Run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			default:

			}

			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				m.Measure()
			}
		}
	}()
}
