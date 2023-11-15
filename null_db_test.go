package qube

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NullDB(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	var db DBIface
	var buf bytes.Buffer
	db = &NullDB{&buf}

	_, err := db.Exec("select 1")
	require.NoError(err)

	_, err = db.ExecContext(context.Background(), "select 2")
	require.NoError(err)

	err = db.Close()
	require.NoError(err)

	assert.Equal("select 1\nselect 2\n", buf.String())
}
