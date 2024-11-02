package util_test

import (
	"testing"

	"github.com/creack/pty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/qube/util"
)

func Test_MustGetTermSize(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	ptmx, tty, err := pty.Open()
	require.NoError(err)
	pty.Setsize(tty, &pty.Winsize{Rows: 30, Cols: 123})

	w := util.MustGetTermSize(ptmx.Fd())
	assert.Equal(123, w)
}
