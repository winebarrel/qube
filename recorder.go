package qube

import (
	"sync"
	"time"
)

type Recorder struct {
	sync.Mutex
	*Options
	ID              string
	DataPoints      []DataPoint
	ErrorQueryCount int
	StartedAt       time.Time
	FinishedAt      time.Time
	ch              chan []DataPoint
	closed          chan struct{}
}

type DataPoint struct {
	Time     time.Time
	Duration time.Duration
	IsError  bool
}

func NewRecorder(id string, options *Options) *Recorder {
	rec := &Recorder{
		Options:    options,
		ID:         id,
		DataPoints: []DataPoint{},
		ch:         make(chan []DataPoint, options.Nagents*3),
	}

	return rec
}

func (rec *Recorder) Start() {
	rec.closed = make(chan struct{})

	push := func(dps []DataPoint) {
		rec.Lock()
		defer rec.Unlock()
		rec.DataPoints = append(rec.DataPoints, dps...)

		for _, v := range dps {
			if v.IsError {
				rec.ErrorQueryCount++
			}
		}
	}

	go func() {
		for dps := range rec.ch {
			push(dps)
		}

		close(rec.closed)
	}()

	rec.StartedAt = time.Now()
}

func (rec *Recorder) Close() {
	rec.FinishedAt = time.Now()
	close(rec.ch)
	<-rec.closed
}

func (rec *Recorder) Add(dps []DataPoint) {
	rec.ch <- dps
}

func (rec *Recorder) Report() *Report {
	return NewReport(rec)
}

func (rec *Recorder) Count() int {
	// Lock to avoid race conditions
	rec.Lock()
	defer rec.Unlock()
	return len(rec.DataPoints)
}
