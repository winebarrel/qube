package util_test

import (
	"bufio"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/qube/util"
)

func Test_ReadLine(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	f, _ := os.CreateTemp("", "")
	defer os.Remove(f.Name())
	f.WriteString(`{"q":"select 1"}` + "\n") //nolint:errcheck
	f.WriteString(`{"q":"select 2"}` + "\n") //nolint:errcheck
	f.WriteString(`{"q":"select 3"}` + "\n") //nolint:errcheck
	f.Sync()                                 //nolint:errcheck
	f.Seek(0, io.SeekStart)

	buf := bufio.NewReader(f)
	line, err := util.ReadLine(buf)

	require.NoError(err)
	assert.Equal([]byte(`{"q":"select 1"}`), line)
}
