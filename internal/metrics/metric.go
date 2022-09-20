package metrics

import "context"

type Measurable interface {
	Run(ctx context.Context)
	Measure()
	ClearOldStat(hoursAgo int)
	GetIndicators(everyN int, durationM int) interface{}
}
