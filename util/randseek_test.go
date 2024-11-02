package util_test

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/qube/util"
)

func Test_RandSeek(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	f, _ := os.CreateTemp("", "")
	defer os.Remove(f.Name())
	f.WriteString(`{"q":"select 1"}` + "\n")
	f.WriteString(`{"q":"select 2"}` + "\n")
	f.WriteString(`{"q":"select 3"}` + "\n")
	f.Sync()
	f.Seek(0, io.SeekStart)

	offsets := map[int64]struct{}{}

	for i := 0; i < 100; i++ {
		err := util.RandSeek(f)
		require.NoError(err)
		offset, _ := f.Seek(0, io.SeekCurrent)
		offsets[offset] = struct{}{}
	}

	assert.True(len(offsets) > 0)
}
