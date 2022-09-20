package app

import (
	"context"
	"reflect"

	"github.com/director74/system_monitoring/internal/cfg"
	"github.com/director74/system_monitoring/internal/metrics"
)

type Application interface {
	BeginCollect(context.Context)
	ClearOldData()
}

type App struct {
	conf    cfg.Configurable
	metrics map[string]metrics.Measurable
}

func NewApplication(conf cfg.Configurable) Application {
	return &App{
		conf: conf,
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
}

func (a *App) ClearOldData() {
	for _, parameter := range a.metrics {
		parameter.ClearOldStat(2)
	}
}
