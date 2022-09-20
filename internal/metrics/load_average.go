package metrics

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

type values struct {
	minute1  float32
	minute5  float32
	minute15 float32
}

type LoadAverage struct {
	*Metric
	storage map[time.Time]values
}

func (l *LoadAverage) GetIndicators(everyN int, durationM int) (interface{}, error) {
	l.storage[time.Now()] = values{}

	return values{}, nil
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

	resultVals := values{}

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
			resultVals.minute1 = float32(value)
		}
		if vv == "minute5" {
			resultVals.minute5 = float32(value)
		}
		if vv == "minute15" {
			resultVals.minute15 = float32(value)
		}
	}
	l.storage[time.Now()] = resultVals

	return nil
}

func (l *LoadAverage) ClearOldStat(minutesAgo int) {

}

func NewLoadAverage() Measurable {
	return &LoadAverage{
		storage: make(map[time.Time]values),
	}
}
