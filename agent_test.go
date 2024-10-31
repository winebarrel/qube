package qube_test

import (
	"bytes"
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/qube"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

func Test_Agent(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	f, _ := os.CreateTemp("", "")
	defer os.Remove(f.Name())
	f.WriteString(`{"q":"select 1"}` + "\n") //nolint:errcheck
	f.Sync()                                 //nolint:errcheck

	var buf bytes.Buffer

	options := &qube.Options{
		AgentOptions: qube.AgentOptions{
			Force: false,
		},
		DataOptions: qube.DataOptions{
			DataFiles:  []string{f.Name()},
			Key:        "q",
			Loop:       true,
			Random:     false,
			CommitRate: 0,
		},
		DBConfig: qube.DBConfig{
			DSN:       testDSN_MySQL,
			Driver:    qube.DBDriverMySQL,
			Noop:      true,
			NullDBOut: &buf,
		},
		Nagents:  1,
		Rate:     0,
		Time:     1 * time.Second,
		Progress: false,
	}

	rec := qube.NewRecorder(testUUID, options)
	limiter := rate.NewLimiter(rate.Limit(1), 1)
	agent, err := qube.NewAgent(testUUID, 0, options, rec, limiter)
	require.NoError(err)

	eg, ctx := errgroup.WithContext(context.Background())
	ctx, cancel := context.WithCancel(ctx)
	eg.Go(func() error { return agent.Start(ctx) })
	time.Sleep(2 * time.Second)
	cancel()
	err = eg.Wait()

	require.True(err == nil || errors.Is(err, context.Canceled))
	assert.Regexp(`select 1`, buf.String())
}
