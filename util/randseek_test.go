package util_test

import (
	"io"
	"os"
	"testing"

	seekable "github.com/SaveTheRbtz/zstd-seekable-format-go/pkg"
	"github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/qube"
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
	fi, err := f.Stat()
	require.NoError(err)

	offsets := map[int64]struct{}{}

	for i := 0; i < 100; i++ {
		err := util.RandSeek(f, fi.Size())
		require.NoError(err)
		offset, _ := f.Seek(0, io.SeekCurrent)
		offsets[offset] = struct{}{}
	}

	assert.True(len(offsets) > 0)
}

func Test_RandSeek_ZstdFile(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	f, _ := os.CreateTemp("", "")
	defer os.Remove(f.Name())

	{
		enc, _ := zstd.NewWriter(nil)
		defer enc.Close()
		w, _ := seekable.NewWriter(f, enc)
		w.Write([]byte(`{"q":"select 1"}` + "\n"))
		w.Write([]byte(`{"q":"select 2"}` + "\n"))
		w.Write([]byte(`{"q":"select 3"}` + "\n"))
		w.Close()
		f.Sync()
		f.Seek(0, io.SeekStart)
	}

	zf, err := qube.NewZstdFile(f)
	require.NoError(err)
	defer zf.Close()

	offsets := map[int64]struct{}{}

	for i := 0; i < 100; i++ {
		err := util.RandSeek(zf, zf.Size())
		require.NoError(err)
		offset, _ := zf.Seek(0, io.SeekCurrent)
		offsets[offset] = struct{}{}
	}

	assert.True(len(offsets) > 0)
}
