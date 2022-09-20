package metrics

import (
	"context"
	"log"
	"time"
)

type Measurable interface {
	Run(context.Context, measureFunc)
	Measure() error
	ClearOldStat(int)
	GetAverageByPeriod(beginTime time.Time, durationM int32) (interface{}, error)
}

type measureFunc func() error

type Metric struct {
}

func (m *Metric) Run(ctx context.Context, measureMethod measureFunc) {
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
				err := measureMethod()
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()
}
