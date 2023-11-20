package qube

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/term"
)

const (
	InterimReportIntvl = 1 * time.Second
)

type Progress struct {
	w           io.Writer
	noop        bool
	prevDPLen   int
	nDeadAgents atomic.Uint64
	closed      chan struct{}
}

func NewProgress(w io.Writer, noop bool) *Progress {
	progress := &Progress{
		w:    w,
		noop: noop,
	}

	return progress
}

func (progress *Progress) Start(ctx context.Context, rec *Recorder) {
	if progress.noop {
		return
	}

	progress.closed = make(chan struct{})
	tk := time.NewTicker(InterimReportIntvl)

	go func() {
	L:
		for {
			select {
			case <-ctx.Done():
				tk.Stop()
				break L
			case <-tk.C:
				progress.report(rec)
			}
		}

		close(progress.closed)
	}()
}

func (progress *Progress) IncrDead() {
	if progress.noop {
		return
	}

	progress.nDeadAgents.Add(1)
}

func (progress *Progress) report(rec *Recorder) {
	if progress.noop {
		return
	}

	dpLen := rec.CountWithoutError()
	delta := dpLen - progress.prevDPLen
	progress.prevDPLen = dpLen
	qps := float64(time.Duration(delta) * time.Second / InterimReportIntvl)
	elapsed := time.Since(rec.StartedAt)
	running := uint64(rec.Nagents) - progress.nDeadAgents.Load()
	width, _, err := term.GetSize(0)

	if err != nil {
		panic(err)
	}

	elapsed = elapsed.Round(time.Second)
	minute := elapsed / time.Minute
	second := (elapsed - minute*time.Minute) / time.Second
	line := fmt.Sprintf("%02d:%02d | %d agents / exec %d queries, %d errors (%.0f qps)", minute, second, running, dpLen, rec.ErrorQueryCount, qps)
	fmt.Fprintf(progress.w, "\r%-*s", width, line)
}

func (progress *Progress) Close() {
	if progress.noop {
		return
	}

	<-progress.closed
	width, _, err := term.GetSize(0)

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(progress.w, "\r"+strings.Repeat(" ", width)+"\r")
}
