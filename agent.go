package qube

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

const (
	RecIntvl = 1 * time.Second
)

type AgentOptions struct {
	Force bool `kong:"negatable,default='false',help='Do not abort test on error. (default: disabled)'"`
}

type Agent struct {
	*AgentOptions
	ID      string
	db      DBIface
	data    *Data
	rec     *Recorder
	limiter *rate.Limiter
}

func NewAgent(taskID string, agentNum uint, options *Options, rec *Recorder, limiter *rate.Limiter) (*Agent, error) {
	autoCommit := options.CommitRate == 0
	db, err := options.OpenDBWithPing(autoCommit)

	if err != nil {
		return nil, err
	}

	data, err := NewData(options)

	if err != nil {
		return nil, err
	}

	agent := &Agent{
		AgentOptions: &options.AgentOptions,
		ID:           fmt.Sprintf("%s/%d", taskID, agentNum),
		db:           db,
		data:         data,
		rec:          rec,
		limiter:      limiter,
	}

	return agent, nil
}

func (agent *Agent) Start(ctx context.Context) error {
	defer func() {
		agent.db.Close()
		agent.data.Close()
	}()

	_, err := agent.db.Exec("select 'start agent - " + agent.ID + "'")

	if err != nil {
		return fmt.Errorf("failed to execute start query (%w)", err)
	}

	err = agent.start0(ctx)

	if err != nil && err != EOD {
		return err
	}

	_, err = agent.db.Exec("select 'exit agent - " + agent.ID + "'")

	if err != nil {
		return fmt.Errorf("failed to execute exit query (%w)", err)
	}

	return nil
}

func (agent *Agent) start0(ctx context.Context) error {
	tkrec := time.NewTicker(RecIntvl)
	defer tkrec.Stop()
	dps := []DataPoint{}

	defer func() {
		if len(dps) > 0 {
			agent.rec.Add(dps)
		}
	}()

L:
	for {
		if agent.limiter != nil {
			agent.limiter.Wait(ctx) //nolint:errcheck
		}

		select {
		case <-ctx.Done():
			break L
		case <-tkrec.C:
			agent.rec.Add(dps)
			// Create new slices to avoid race conditions
			dps = []DataPoint{}
		default:
			// Nothing to do
		}

		q, err := agent.data.Next()

		if err != nil {
			return err
		}

		dur, err := agent.execQuery(ctx, q)

		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			continue
		} else if err != nil && !agent.Force {
			return fmt.Errorf("failed to execute query - %s (%w)", q, err)
		}

		dps = append(dps, DataPoint{
			Time:     time.Now(),
			Duration: dur,
			IsError:  err != nil,
		})
	}

	return nil
}

func (agent *Agent) execQuery(ctx context.Context, q string) (time.Duration, error) {
	start := time.Now()
	_, err := agent.db.ExecContext(ctx, q)
	end := time.Now()

	if err != nil {
		return 0, err
	}

	return end.Sub(start), nil
}
