package util_test

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/qube/util"
)

func Test_RandSeek(t *testing.T) {
	assert := assert.New(t)

	f, _ := os.CreateTemp("", "")
	defer os.Remove(f.Name())
	f.WriteString(`{"q":"select 1"}` + "\n") //nolint:errcheck
	f.WriteString(`{"q":"select 2"}` + "\n") //nolint:errcheck
	f.WriteString(`{"q":"select 3"}` + "\n") //nolint:errcheck
	f.Sync()                                 //nolint:errcheck
	f.Seek(0, io.SeekStart)

	offsets := map[int64]struct{}{}

	for i := 0; i < 100; i++ {
		util.RandSeek(f)
		offset, _ := f.Seek(0, io.SeekCurrent)
		offsets[offset] = struct{}{}
	}

	assert.True(len(offsets) > 0)
}
