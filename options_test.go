package qube_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/qube"
)

func Test_Options_AfterApply(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	ds := []struct {
		DSN    string
		Driver qube.DBDriver
	}{
		{
			DSN:    testDSN_MySQL,
			Driver: qube.DBDriverMySQL,
		},
		{
			DSN:    testDSN_PostgreSQL,
			Driver: qube.DBDriverPostgreSQL,
		},
	}

	for _, d := range ds {
		options := qube.Options{
			Time: 3 * time.Second,
			DBConfig: qube.DBConfig{
				DSN: d.DSN,
			},
		}

		err := options.AfterApply()
		require.NoError(err)

		assert.Equal(qube.JSONDuration(3*time.Second), options.X_Time)
		assert.Equal(d.Driver, options.Driver)
	}
}

func Test_Options_AfterApply_InvalidDSN(t *testing.T) {
	assert := assert.New(t)

	options := qube.Options{
		DBConfig: qube.DBConfig{
			DSN: "**invalid**",
		},
	}

	err := options.AfterApply()
	assert.ErrorContains(err, "cannot parse DSN - **invalid**")
}
