package qube_test

import (
	"io"
	"os"
	"testing"

	seekable "github.com/SaveTheRbtz/zstd-seekable-format-go/pkg"
	"github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/qube"
)

func Test_ZstdFile(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	f, _ := os.CreateTemp("", "")
	defer os.Remove(f.Name())

	{
		enc, _ := zstd.NewWriter(nil)
		defer enc.Close()
		w, _ := seekable.NewWriter(f, enc)
		w.Write([]byte(`{"q":"select 1"}`))
		w.Close()
		f.Sync()
		f.Seek(0, io.SeekStart)
	}

	zf, err := qube.NewZstdFile(f)
	require.NoError(err)

	_, err = zf.Seek(5, io.SeekStart)
	require.NoError(err)

	buf := make([]byte, 10)
	_, err = zf.Read(buf)
	require.NoError(err)
	assert.Equal(`"select 1"`, string(buf))

	err = zf.Close()
	require.NoError(err)
}
