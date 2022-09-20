package app

import (
	"context"
	"reflect"
	"time"

	"github.com/director74/system_monitoring/internal/cfg"
	"github.com/director74/system_monitoring/internal/metrics"
)

type Application interface {
	BeginCollect(context.Context)
	ClearOldData(context.Context, int)
}

type App struct {
	conf    cfg.Configurable
	metrics map[string]metrics.Measurable
}

func NewApplication(conf cfg.Configurable) Application {
	return &App{
		conf:    conf,
		metrics: make(map[string]metrics.Measurable),
	}
}

func (a *App) BeginCollect(ctx context.Context) {
	trackParams := a.conf.GetAllowedForTracking()
	refStruct := reflect.ValueOf(trackParams)
	refStructType := refStruct.Type()
	for i := 0; i < refStruct.NumField(); i++ {
		fieldDescr := refStructType.Field(i)
		fieldVal := refStruct.Field(i)
		if fieldVal.Bool() == true {
			switch fieldDescr.Name {
			case "LoadAverage":
				a.metrics[fieldDescr.Name] = metrics.NewLoadAverage()
			}
		}
	}

	for _, param := range a.metrics {
		param.Run(ctx, param.Measure)
	}
}

func (a *App) ClearOldData(ctx context.Context, minutesAgo int) {
	ticker := time.NewTicker(10 * time.Minute)
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
			for _, parameter := range a.metrics {
				parameter.ClearOldStat(minutesAgo)
			}
		}
	}
}
