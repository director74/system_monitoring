package app

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/director74/system_monitoring/internal/cfg"
	"github.com/director74/system_monitoring/internal/metrics"
)

type Application interface {
	BeginCollect(context.Context)
	ClearOldData(context.Context, int)
	GetConfig() cfg.Configurable
	GetMetricStat(name string) (metrics.Measurable, error)
	GetAllMetricNames() []string
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

func (a *App) ClearOldData(ctx context.Context, olderMinutes int) {
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
				parameter.ClearOldStat(olderMinutes)
			}
		}
	}
}

func (a *App) GetConfig() cfg.Configurable {
	return a.conf
}

func (a *App) GetMetricStat(name string) (metrics.Measurable, error) {
	metric, ok := a.metrics[name]
	if !ok {
		return nil, fmt.Errorf("metric %s not found", name)
	}
	return metric, nil
}

func (a *App) GetAllMetricNames() []string {
	result := make([]string, len(a.metrics))
	for name, _ := range a.metrics {
		result = append(result, name)
	}

	return result
}
