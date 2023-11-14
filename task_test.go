package qube_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/qube"
)

func Test_Task_Acc(t *testing.T) {
	if !testAcc {
		t.Skip()
	}

	assert := assert.New(t)
	require := require.New(t)

	tt := []struct {
		URL    string
		Driver qube.DBDriver
	}{
		{URL: testMySQLURL, Driver: qube.DBDriverMySQL},
		{URL: testPostgreSQLURL, Driver: qube.DBDriverPostgreSQL},
	}

	for _, t := range tt {
		f, _ := os.CreateTemp("", "")
		defer os.Remove(f.Name())
		f.WriteString(`{"q":"select 1"}` + "\n") //nolint:errcheck
		f.Sync()                                 //nolint:errcheck

		task := &qube.Task{
			Options: &qube.Options{
				AgentOptions: qube.AgentOptions{
					Force: false,
				},
				DataOptions: qube.DataOptions{
					DataFile:   f.Name(),
					Key:        "q",
					Loop:       true,
					Random:     false,
					CommitRate: 0,
				},
				DBConfig: qube.DBConfig{
					DSN:    t.URL,
					Driver: t.Driver,
					Noop:   false,
				},
				Nagents:  1,
				Rate:     0,
				Time:     1 * time.Second,
				Progress: false,
			},
			ID: testUUID,
		}

		report, err := task.Run()

		require.NoError(err)
		assert.Equal(testUUID, report.ID)
		assert.NotEqual(0, report.AvgQPS)
	}
}
