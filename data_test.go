package qube_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/qube"
)

func Test_Data(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	f, _ := os.CreateTemp("", "")
	defer os.Remove(f.Name())
	f.WriteString(`{"q":"select 1"}` + "\n") //nolint:errcheck
	f.Sync()                                 //nolint:errcheck

	options := &qube.Options{
		DataOptions: qube.DataOptions{
			DataFile: f.Name(),
			Key:      "q",
		},
	}

	data, err := qube.NewData(options)
	require.NoError(err)
	defer data.Close()

	q, err := data.Next()
	require.NoError(err)
	assert.Equal("select 1", q)
	_, err = data.Next()
	assert.ErrorIs(err, qube.EOD)
}

func Test_Data_Loop(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	f, _ := os.CreateTemp("", "")
	defer os.Remove(f.Name())
	f.WriteString(`{"q":"select 1"}` + "\n") //nolint:errcheck
	f.Sync()                                 //nolint:errcheck

	options := &qube.Options{
		DataOptions: qube.DataOptions{
			DataFile: f.Name(),
			Key:      "q",
			Loop:     true,
		},
	}

	data, err := qube.NewData(options)
	require.NoError(err)
	defer data.Close()

	q, err := data.Next()
	require.NoError(err)
	assert.Equal("select 1", q)
	q, err = data.Next()
	require.NoError(err)
	assert.Equal("select 1", q)
}

func Test_Data_Random(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	f, _ := os.CreateTemp("", "")
	defer os.Remove(f.Name())
	f.WriteString(`{"q":"select 1"}` + "\n") //nolint:errcheck
	f.Sync()                                 //nolint:errcheck

	options := &qube.Options{
		DataOptions: qube.DataOptions{
			DataFile: f.Name(),
			Key:      "q",
			Loop:     true,
			Random:   true,
		},
	}

	data, err := qube.NewData(options)
	require.NoError(err)
	defer data.Close()

	q, err := data.Next()
	require.NoError(err)
	assert.Equal("select 1", q)
	q, err = data.Next()
	require.NoError(err)
	assert.Equal("select 1", q)
}

func Test_Data_WithCommitRate(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	f, _ := os.CreateTemp("", "")
	defer os.Remove(f.Name())
	f.WriteString(`{"q":"select 1"}` + "\n") //nolint:errcheck
	f.Sync()                                 //nolint:errcheck

	options := &qube.Options{
		DataOptions: qube.DataOptions{
			DataFile:   f.Name(),
			Key:        "q",
			Loop:       true,
			CommitRate: 2,
		},
	}

	data, err := qube.NewData(options)
	require.NoError(err)
	defer data.Close()

	qs := []string{
		"begin",
		"select 1",
		"select 1",
		"commit",
		"begin",
		"select 1",
		"select 1",
		"commit",
	}

	for _, expected := range qs {
		q, err := data.Next()
		require.NoError(err)
		assert.Equal(expected, q)
	}
}

func Test_Data_WithoutKey(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	f, _ := os.CreateTemp("", "")
	defer os.Remove(f.Name())
	f.WriteString(`{"_q_":"select 1"}` + "\n") //nolint:errcheck
	f.Sync()                                   //nolint:errcheck

	options := &qube.Options{
		DataOptions: qube.DataOptions{
			DataFile: f.Name(),
			Key:      "q",
			Loop:     true,
		},
	}

	data, err := qube.NewData(options)
	require.NoError(err)

	_, err = data.Next()
	assert.ErrorContains(err, `failed to get query field "q" from '{"_q_":"select 1"}'`)
}
