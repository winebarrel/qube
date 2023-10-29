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
	nDeadAgents int32
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
	}()
}

func (progress *Progress) IncrDead() {
	if progress.noop {
		return
	}

	atomic.AddInt32(&progress.nDeadAgents, 1)
}

func (progress *Progress) report(rec *Recorder) {
	if progress.noop {
		return
	}

	dpLen := len(rec.DataPoints)
	delta := dpLen - progress.prevDPLen
	progress.prevDPLen = dpLen
	qps := float64(time.Duration(delta) * time.Second / InterimReportIntvl)
	elapsed := time.Since(rec.StartedAt)
	running := rec.Nagents - int(progress.nDeadAgents)
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

func (progress *Progress) Clear() {
	if progress.noop {
		return
	}

	width, _, err := term.GetSize(0)

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(progress.w, "\r"+strings.Repeat(" ", width)+"\r")
}
