package qube

import (
	"sync"
	"time"
)

type Recorder struct {
	mu sync.Mutex
	*Options
	ID              string
	dataPoints      []DataPoint
	ErrorQueryCount int
	StartedAt       time.Time
	FinishedAt      time.Time
	ch              chan []DataPoint
	closed          chan struct{}
}

type DataPoint struct {
	Time     int64
	Duration time.Duration
	IsError  bool
}

func NewRecorder(id string, options *Options) *Recorder {
	rec := &Recorder{
		Options:    options,
		ID:         id,
		dataPoints: []DataPoint{},
		ch:         make(chan []DataPoint, options.Nagents*3),
	}

	return rec
}

func (rec *Recorder) Start() {
	rec.closed = make(chan struct{})

	push := func(dps []DataPoint) {
		rec.mu.Lock()
		defer rec.mu.Unlock()

		for _, v := range dps {
			if !v.IsError {
				rec.dataPoints = append(rec.dataPoints, v)
			} else {
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

func (rec *Recorder) CountAll() int {
	// Lock to avoid race conditions
	rec.mu.Lock()
	defer rec.mu.Unlock()
	return len(rec.dataPoints) + rec.ErrorQueryCount
}

func (rec *Recorder) CountSuccess() int {
	// Lock to avoid race conditions
	rec.mu.Lock()
	defer rec.mu.Unlock()
	return len(rec.dataPoints)
}

func (rec *Recorder) DataPoints() []DataPoint {
	return rec.dataPoints
}
