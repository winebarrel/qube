package qube

import (
	"encoding/json"
	"runtime"
	"time"

	"github.com/jamiealquiza/tachymeter"
)

type JSONDuration time.Duration

func (jd JSONDuration) MarshalJSON() (b []byte, err error) {
	d := time.Duration(jd)
	return []byte(`"` + d.String() + `"`), nil
}

type Report struct {
	ID              string
	StartedAt       time.Time
	FinishedAt      time.Time
	ElapsedTime     JSONDuration
	Options         *Options
	GOMAXPROCS      int
	QueryCount      int
	ErrorQueryCount int
	AvgQPS          float64
	MaxQPS          float64
	MinQPS          float64
	MedianQPS       float64
	Duration        *tachymeter.Metrics
}

func NewReport(rec *Recorder) *Report {
	dps := rec.DataPoints()
	dpOkLen := len(dps)

	if rec.CountAll() == 0 {
		return nil
	}

	nanoElapsed := rec.FinishedAt.Sub(rec.StartedAt)

	report := &Report{
		ID:              rec.ID,
		StartedAt:       rec.StartedAt,
		FinishedAt:      rec.FinishedAt,
		ElapsedTime:     JSONDuration(nanoElapsed),
		Options:         rec.Options,
		GOMAXPROCS:      runtime.GOMAXPROCS(0),
		QueryCount:      rec.CountAll(),
		ErrorQueryCount: rec.ErrorQueryCount(),
	}

	if dpOkLen > 0 {
		report.AvgQPS = float64(time.Duration(dpOkLen) * time.Second / nanoElapsed)

		t := tachymeter.New(&tachymeter.Config{
			Size:  dpOkLen,
			HBins: 10,
		})

		for _, v := range dps {
			t.AddTime(v.Duration)
		}

		report.Duration = t.Calc()
		qpsSet := NewQPSSet(dps)
		report.MinQPS, report.MaxQPS, report.MedianQPS = qpsSet.Stats()
	}

	return report
}

func (report *Report) JSON() string {
	rawJson, _ := json.MarshalIndent(report, "", "  ")
	return string(rawJson)
}
