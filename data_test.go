package qube_test

import (
	"fmt"
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
	f.WriteString(`{"q":"select 1"}` + "\n")
	f.Sync()

	options := &qube.Options{
		DataOptions: qube.DataOptions{
			DataFiles: []string{f.Name()},
			Key:       "q",
		},
	}

	data, err := qube.NewData(options, 0)
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
	f.WriteString(`{"q":"select 1"}` + "\n")
	f.Sync()

	options := &qube.Options{
		DataOptions: qube.DataOptions{
			DataFiles: []string{f.Name()},
			Key:       "q",
			Loop:      true,
		},
	}

	data, err := qube.NewData(options, 0)
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
	f.WriteString(`{"q":"select 1"}` + "\n")
	f.Sync()

	options := &qube.Options{
		DataOptions: qube.DataOptions{
			DataFiles: []string{f.Name()},
			Key:       "q",
			Loop:      true,
			Random:    true,
		},
	}

	data, err := qube.NewData(options, 0)
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
	f.WriteString(`{"q":"select 1"}` + "\n")
	f.Sync()

	options := &qube.Options{
		DataOptions: qube.DataOptions{
			DataFiles:  []string{f.Name()},
			Key:        "q",
			Loop:       true,
			CommitRate: 2,
		},
	}

	data, err := qube.NewData(options, 0)
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
	f.WriteString(`{"_q_":"select 1"}` + "\n")
	f.Sync()

	options := &qube.Options{
		DataOptions: qube.DataOptions{
			DataFiles: []string{f.Name()},
			Key:       "q",
			Loop:      true,
		},
	}

	data, err := qube.NewData(options, 0)
	require.NoError(err)

	_, err = data.Next()
	assert.ErrorContains(err, `failed to get query field "q" from '{"_q_":"select 1"}'`)
}

func Test_MultiData(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	f1, _ := os.CreateTemp("", "")
	defer os.Remove(f1.Name())
	f1.WriteString(`{"q":"select 1"}` + "\n")
	f1.Sync()

	f2, _ := os.CreateTemp("", "")
	defer os.Remove(f2.Name())
	f2.WriteString(`{"q":"select 2"}` + "\n")
	f2.Sync()

	f3, _ := os.CreateTemp("", "")
	defer os.Remove(f3.Name())
	f3.WriteString(`{"q":"select 3"}` + "\n")
	f3.Sync()

	options := &qube.Options{
		DataOptions: qube.DataOptions{
			DataFiles: []string{f1.Name(), f2.Name(), f3.Name()},
			Key:       "q",
		},
	}

	for i := range 6 {
		data, err := qube.NewData(options, uint64(i))
		require.NoError(err)
		defer data.Close()

		q, err := data.Next()
		require.NoError(err)
		assert.Equal(fmt.Sprintf("select %d", i%3+1), q)
		_, err = data.Next()
		assert.ErrorIs(err, qube.EOD)
	}
}

func Test_Data_WithEmptyLine(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	f, _ := os.CreateTemp("", "")
	defer os.Remove(f.Name())
	f.WriteString(`{"q":"select 1"}` + "\n\n\n" + `{"q":"select 2"}` + "\n")
	f.Sync()

	options := &qube.Options{
		DataOptions: qube.DataOptions{
			DataFiles: []string{f.Name()},
			Key:       "q",
		},
	}

	data, err := qube.NewData(options, 0)
	require.NoError(err)
	defer data.Close()

	q, err := data.Next()
	require.NoError(err)
	assert.Equal("select 1", q)

	q, err = data.Next()
	require.NoError(err)
	assert.Equal("select 2", q)

	_, err = data.Next()
	assert.ErrorIs(err, qube.EOD)
}

func Test_Data_WithCommentOut(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	f, _ := os.CreateTemp("", "")
	defer os.Remove(f.Name())
	f.WriteString(`{"q":"select 1"}` + "\n" + `//{"q":"select 2"}` + "\n" + `{"q":"select 3"}` + "\n")
	f.Sync()

	options := &qube.Options{
		DataOptions: qube.DataOptions{
			DataFiles: []string{f.Name()},
			Key:       "q",
		},
	}

	data, err := qube.NewData(options, 0)
	require.NoError(err)
	defer data.Close()

	q, err := data.Next()
	require.NoError(err)
	assert.Equal("select 1", q)

	q, err = data.Next()
	require.NoError(err)
	assert.Equal("select 3", q)

	_, err = data.Next()
	assert.ErrorIs(err, qube.EOD)
}

func Test_Data_Empty(t *testing.T) {
	assert := assert.New(t)

	f, _ := os.CreateTemp("", "")
	defer os.Remove(f.Name())
	f.Sync()

	options := &qube.Options{
		DataOptions: qube.DataOptions{
			DataFiles: []string{f.Name()},
			Key:       "q",
		},
	}

	_, err := qube.NewData(options, 0)
	assert.ErrorContains(err, "test data is empty")
}
