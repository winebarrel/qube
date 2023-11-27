package qube_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/creack/pty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/qube"
)

func Test_Progress(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	pty, tty, err := pty.Open()
	require.NoError(err)
	defer pty.Close()
	defer tty.Close()

	rec := qube.NewRecorder(testUUID, &qube.Options{})
	progress := qube.NewProgress(tty, false)
	ctx, cancel := context.WithCancel(context.Background())

	progress.Start(ctx, rec)
	time.Sleep(1 * time.Second)
	cancel()
	progress.Close()

	buf := make([]byte, 1024)
	_, err = pty.Read(buf)
	require.NoError(err)
	assert.Equal("00:01 | 0 agents / exec 0 queries, 0 errors (0 qps)", strings.Trim(string(buf), "\r\x00"))
}
