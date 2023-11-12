package qube

import (
	"context"
	"errors"
	"os"
	"os/signal"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

type Task struct {
	*Options
	ID string
}

func NewTask(options *Options) *Task {
	task := &Task{
		Options: options,
		ID:      uuid.NewString(),
	}

	return task
}

func (task *Task) makeAgents() ([]*Agent, *Recorder, error) {
	agents := make([]*Agent, task.Nagents)
	rec := NewRecorder(task.ID, task.Options)
	limiter := rate.NewLimiter(rate.Limit(task.Rate), 1)

	for i := 0; i < task.Nagents; i++ {
		var err error
		agents[i], err = NewAgent(task.ID, i, task.Options, rec, limiter)

		if err != nil {
			return nil, nil, err
		}
	}

	return agents, rec, nil
}

func (task *Task) Run() (*Report, error) {
	agents, rec, err := task.makeAgents()

	if err != nil {
		return nil, err
	}

	eg, ctx := errgroup.WithContext(context.Background())
	ctx, cancel := context.WithCancel(ctx)

	if task.Time > 0 {
		ctx, cancel = context.WithTimeout(ctx, task.Time)
	}

	// Trap SIGINT
	{
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)

		go func() {
			select {
			case <-ctx.Done():
				// Nothing to do
			case <-sigint:
				// Stop query on interrupt
				cancel()
				eg.Wait() //nolint:errcheck
				os.Exit(130)
			}
		}()
	}

	// Timeout
	if task.Time > 0 {
		go func() {
			<-ctx.Done()

			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				cancel()
			}
		}()
	}

	var progress = NewProgress(os.Stderr, !task.Progress || task.Noop)
	rec.Start()

	for _, v := range agents {
		agent := v

		eg.Go(func() error {
			err := agent.Start(ctx)
			progress.IncrDead()
			return err
		})
	}

	progress.Start(ctx, rec)
	err = eg.Wait()
	cancel()
	progress.Close() // Wait for ticker to stop
	rec.Close()      // Wait for buffer flush

	if err != nil && !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
		return nil, err
	}

	return rec.Report(), nil
}
