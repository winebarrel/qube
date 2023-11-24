package qube

import (
	"sync"
	"sync/atomic"
	"time"
)

type Recorder struct {
	mu sync.Mutex
	*Options
	ID              string
	dataPoints      []DataPoint
	errorQueryCount atomic.Int64 // Use atomic to avoid race conditions
	StartedAt       time.Time
	FinishedAt      time.Time
	ch              chan []DataPointWithErr
	closed          chan struct{}
}

type DataPoint struct {
	Time     int64
	Duration time.Duration
}

type DataPointWithErr struct {
	DataPoint
	IsError bool
}

func NewRecorder(id string, options *Options) *Recorder {
	rec := &Recorder{
		Options:    options,
		ID:         id,
		dataPoints: []DataPoint{},
		ch:         make(chan []DataPointWithErr, options.Nagents*3),
	}

	return rec
}

func (rec *Recorder) Start() {
	rec.closed = make(chan struct{})

	push := func(dpes []DataPointWithErr) {
		rec.mu.Lock()
		defer rec.mu.Unlock()
		var errCnt int64

		for _, v := range dpes {
			if !v.IsError {
				rec.dataPoints = append(rec.dataPoints, v.DataPoint)
			} else {
				errCnt++
			}
		}

		rec.errorQueryCount.Add(errCnt)
	}

	go func() {
		for dpes := range rec.ch {
			push(dpes)
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

func (rec *Recorder) Add(dpes []DataPointWithErr) {
	rec.ch <- dpes
}

func (rec *Recorder) Report() *Report {
	return NewReport(rec)
}

func (rec *Recorder) CountAll() int {
	// Lock to avoid race conditions
	rec.mu.Lock()
	defer rec.mu.Unlock()
	return len(rec.dataPoints) + int(rec.errorQueryCount.Load())
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

func (rec *Recorder) ErrorQueryCount() int {
	return int(rec.errorQueryCount.Load())
}
