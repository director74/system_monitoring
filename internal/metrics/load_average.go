package metrics

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type LoadAverageResult struct {
	Minute1  float32
	Minute5  float32
	Minute15 float32
}

type LoadAverage struct {
	*Metric
	mu      sync.RWMutex
	storage map[int64]LoadAverageResult
}

func (l *LoadAverage) GetAverageByPeriod(measures chan MeasureResult, beginTimeUnix int64, endTimeUnix int64) {
	var cntr float32

	result := make(MeasureResult)
	calculated := LoadAverageResult{}
	for timeIndex, mark := range l.storage {
		if timeIndex >= beginTimeUnix && timeIndex <= endTimeUnix {
			cntr++
			calculated.Minute1 += mark.Minute1
			calculated.Minute5 += mark.Minute5
			calculated.Minute15 += mark.Minute15
		}
	}

	if cntr > 0 {
		result["LoadAverage"] = LoadAverageResult{Minute1: calculated.Minute1 / cntr, Minute5: calculated.Minute5, Minute15: calculated.Minute15}
	}
	measures <- result
}

func (l *LoadAverage) Measure() error {
	var value float64
	var out bytes.Buffer

	cmd := exec.Command("bash", "-c", "uptime")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("cant measure load average: %w", err)
	}

	resultVals := LoadAverageResult{}

	re, err := regexp.Compile(`load average: (?P<minute1>\d+\.\d+)+,\s*(?P<minute5>\d+\.\d+)+,\s*(?P<minute15>\d+\.\d+)+`)
	if err != nil {
		return fmt.Errorf("regexp problem in load average: %w", err)
	}

	res := re.FindStringSubmatch(out.String())
	for kk, vv := range re.SubexpNames() {
		if vv != "" {
			value, err = strconv.ParseFloat(res[kk], 32)
			if err != nil {
				return fmt.Errorf("cant convert value in load average: %w", err)
			}
		}

		if vv == "minute1" {
			resultVals.Minute1 = float32(value)
		}
		if vv == "minute5" {
			resultVals.Minute5 = float32(value)
		}
		if vv == "minute15" {
			resultVals.Minute15 = float32(value)
		}
	}
	l.mu.Lock()
	l.storage[time.Now().Unix()] = resultVals
	l.mu.Unlock()

	return nil
}

func (l *LoadAverage) ClearOldStat(olderMinutes int) {

}

func NewLoadAverage() Measurable {
	return &LoadAverage{
		storage: make(map[int64]LoadAverageResult),
	}
}
