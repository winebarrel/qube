package qube_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/qube"
)

func TestAcc_DBConfig(t *testing.T) {
	if !testAcc {
		t.Skip()
	}

	require := require.New(t)

	ds := []struct {
		DSN        string
		Driver     qube.DBDriver
		AutoCommit bool
	}{
		{
			DSN:        testDSN_MySQL,
			Driver:     qube.DBDriverMySQL,
			AutoCommit: true,
		},
		{
			DSN:        testDSN_MySQL,
			Driver:     qube.DBDriverMySQL,
			AutoCommit: false,
		},
		{
			DSN:        testDSN_PostgreSQL,
			Driver:     qube.DBDriverPostgreSQL,
			AutoCommit: false,
		},
	}

	for _, d := range ds {
		config := &qube.DBConfig{
			DSN:    d.DSN,
			Driver: d.Driver,
			Noop:   false,
		}

		db, err := config.OpenDBWithPing(d.AutoCommit)
		require.NoError(err)
		defer db.Close()

		_, err = db.ExecContext(context.Background(), "select 1")
		require.NoError(err)
		_, err = db.ExecContext(context.Background(), "select 2")
		require.NoError(err)
	}
}

func Test_DBConfig_Noop(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	var buf bytes.Buffer

	config := &qube.DBConfig{
		DSN:       testDSN_MySQL,
		Driver:    qube.DBDriverMySQL,
		Noop:      true,
		NullDBOut: &buf,
	}

	db, err := config.OpenDBWithPing(true)
	require.NoError(err)
	defer db.Close()

	_, err = db.ExecContext(context.Background(), "select 1")
	require.NoError(err)
	_, err = db.ExecContext(context.Background(), "select 2")
	require.NoError(err)

	assert.Equal("select 1\nselect 2\n", buf.String())
}
