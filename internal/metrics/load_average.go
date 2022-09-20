package metrics

import (
	"context"
	"time"
)

//$ uptime
type values struct {
	minute1  float32
	minute5  float32
	minute15 float32
}

type LoadAverage struct {
	storage map[time.Time]values
}

func (l *LoadAverage) GetIndicators(everyN int, durationM int) interface{} {
	//exec.Command("uptime")
	l.storage[time.Now()] = values{}

	return values{}
}

func (l *LoadAverage) Measure() {
	//exec.Command("uptime")
	l.storage[time.Now()] = values{}
}

func (l *LoadAverage) Run(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(1 * time.Second)
				l.Measure()
			}
		}
	}()
}

func (l *LoadAverage) ClearOldStat(hoursAgo int) {

}

func NewLoadAverage() Measurable {
	return &LoadAverage{}
}
