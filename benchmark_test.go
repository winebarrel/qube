package qube_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/winebarrel/qube"
)

func BenchmarkMySQL(b *testing.B) {
	f := filepath.Join(b.TempDir(), "data.jsonl")
	data := []byte(`{"q":"select 1"}` + "\n")
	os.WriteFile(f, data, 0400)

	task := &qube.Task{
		Options: &qube.Options{
			AgentOptions: qube.AgentOptions{
				Force: false,
			},
			DataOptions: qube.DataOptions{
				DataFiles: []string{f},
				Key:       "q",
				Loop:      true,
			},
			DBConfig: qube.DBConfig{
				Noop: false,
			},
			Nagents:  uint64(runtime.GOMAXPROCS(0)),
			Time:     20 * time.Second,
			Progress: false,
		},
		ID: testUUID,
	}

	task.DSN = qube.DSN(testDSN_MySQL)
	task.Driver = qube.DBDriverMySQL
	report, err := task.Run()
	require.NoError(b, err)
	b.Logf("QueryCount: %d, AvgQPS: %.0f, MedianQPS: %.0f\n", report.QueryCount, report.AvgQPS, report.MedianQPS)
}

func BenchmarkPostgreSQL(b *testing.B) {
	f := filepath.Join(b.TempDir(), "data.jsonl")
	data := []byte(`{"q":"select 1"}` + "\n")
	os.WriteFile(f, data, 0400)

	task := &qube.Task{
		Options: &qube.Options{
			AgentOptions: qube.AgentOptions{
				Force: false,
			},
			DataOptions: qube.DataOptions{
				DataFiles: []string{f},
				Key:       "q",
				Loop:      true,
			},
			DBConfig: qube.DBConfig{
				Noop: false,
			},
			Nagents:  uint64(runtime.GOMAXPROCS(0)),
			Time:     20 * time.Second,
			Progress: false,
		},
		ID: testUUID,
	}

	task.DSN = qube.DSN(testDSN_PostgreSQL)
	task.Driver = qube.DBDriverPostgreSQL
	report, err := task.Run()
	require.NoError(b, err)
	b.Logf("QueryCount:%d, AvgQPS: %.0f, MedianQPS: %.0f\n", report.QueryCount, report.AvgQPS, report.MedianQPS)
}
