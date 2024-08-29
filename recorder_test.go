package qube_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/winebarrel/qube"
)

func Test_Recorder(t *testing.T) {
	assert := assert.New(t)

	options := &qube.Options{
		AgentOptions: qube.AgentOptions{
			Force: true,
		},
		DataOptions: qube.DataOptions{
			DataFile:   "data.jsonl",
			Key:        "q",
			Loop:       true,
			Random:     false,
			CommitRate: 0,
		},
		DBConfig: qube.DBConfig{
			DSN:    testDSN_MySQL,
			Driver: qube.DBDriverMySQL,
			Noop:   false,
		},
		Nagents:  1,
		Rate:     0,
		Time:     1 * time.Second,
		Progress: false,
	}

	deps := []qube.DataPointWithErr{
		// 1 qps
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 14, 0, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 15, 0, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 16, 0, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		// 2 qps
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 17, 0, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 17, 1, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		// 6 qps
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 18, 0, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 18, 1, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 18, 2, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 18, 3, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 18, 4, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 18, 5, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		// 7 qps
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 19, 0, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 19, 1, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 19, 2, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 19, 3, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 19, 4, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 19, 5, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 19, 6, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
	}

	rec := qube.NewRecorder(testUUID, options)
	rec.Start()
	rec.Add(deps)
	rec.Close()

	report := rec.Report()
	// Set dummy value
	report.StartedAt = time.Time{}
	report.FinishedAt = time.Time{}
	report.ElapsedTime = 0
	report.AvgQPS = 0
	report.GOMAXPROCS = 10

	assert.Equal(`{
  "ID": "473d2574-4d1c-46cf-a275-5f3541eb47b7",
  "StartedAt": "0001-01-01T00:00:00Z",
  "FinishedAt": "0001-01-01T00:00:00Z",
  "ElapsedTime": "0s",
  "Options": {
    "Force": true,
    "DataFile": "data.jsonl",
    "Key": "q",
    "Loop": true,
    "Random": false,
    "CommitRate": 0,
    "DSN": "root@tcp(127.0.0.1:13306)/",
    "Driver": "mysql",
    "Noop": false,
    "Nagents": 1,
    "Rate": 0,
    "Time": "0s"
  },
  "GOMAXPROCS": 10,
  "QueryCount": 18,
  "ErrorQueryCount": 0,
  "AvgQPS": 0,
  "MaxQPS": 7,
  "MinQPS": 1,
  "MedianQPS": 4,
  "Duration": {
    "Count": 18,
    "Histogram": [
      {
        "1ms - 1ms": 18
      }
    ],
    "Rate": {
      "Second": 1000
    },
    "Samples": 18,
    "Time": {
      "Avg": "1ms",
      "Cumulative": "18ms",
      "HMean": "999.999µs",
      "Long5p": "1ms",
      "Max": "1ms",
      "Min": "1ms",
      "P50": "1ms",
      "P75": "1ms",
      "P95": "1ms",
      "P99": "1ms",
      "P999": "1ms",
      "Range": "0s",
      "Short5p": "1ms",
      "StdDev": "0s"
    }
  }
}`, report.JSON())
}

func Test_Recorder_WithError(t *testing.T) {
	assert := assert.New(t)

	options := &qube.Options{
		AgentOptions: qube.AgentOptions{
			Force: true,
		},
		DataOptions: qube.DataOptions{
			DataFile:   "data.jsonl",
			Key:        "q",
			Loop:       true,
			Random:     false,
			CommitRate: 0,
		},
		DBConfig: qube.DBConfig{
			DSN:    testDSN_MySQL,
			Driver: qube.DBDriverMySQL,
			Noop:   false,
		},
		Nagents:  1,
		Rate:     0,
		Time:     1 * time.Second,
		Progress: false,
	}

	dpes := []qube.DataPointWithErr{
		// 1 qps
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 14, 0, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 15, 0, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 16, 0, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		// 2 qps
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 17, 0, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 17, 1, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		// 6 qps
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 18, 0, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 18, 1, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 18, 2, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 18, 3, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 18, 4, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 18, 5, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		// 7 qps
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 19, 0, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 19, 1, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 19, 2, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 19, 3, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 19, 4, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 19, 5, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 19, 6, time.UTC).Unix(), Duration: 1 * time.Millisecond}, false},
		// error
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 20, 0, time.UTC).Unix(), Duration: 1 * time.Millisecond}, true},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 21, 0, time.UTC).Unix(), Duration: 1 * time.Millisecond}, true},
		{qube.DataPoint{Time: time.Date(2023, 10, 11, 12, 13, 22, 0, time.UTC).Unix(), Duration: 1 * time.Millisecond}, true},
	}

	rec := qube.NewRecorder(testUUID, options)
	rec.Start()
	rec.Add(dpes)
	rec.Close()

	report := rec.Report()
	// Set dummy value
	report.StartedAt = time.Time{}
	report.FinishedAt = time.Time{}
	report.ElapsedTime = 0
	report.AvgQPS = 0
	report.GOMAXPROCS = 10

	assert.Equal(`{
  "ID": "473d2574-4d1c-46cf-a275-5f3541eb47b7",
  "StartedAt": "0001-01-01T00:00:00Z",
  "FinishedAt": "0001-01-01T00:00:00Z",
  "ElapsedTime": "0s",
  "Options": {
    "Force": true,
    "DataFile": "data.jsonl",
    "Key": "q",
    "Loop": true,
    "Random": false,
    "CommitRate": 0,
    "DSN": "root@tcp(127.0.0.1:13306)/",
    "Driver": "mysql",
    "Noop": false,
    "Nagents": 1,
    "Rate": 0,
    "Time": "0s"
  },
  "GOMAXPROCS": 10,
  "QueryCount": 21,
  "ErrorQueryCount": 3,
  "AvgQPS": 0,
  "MaxQPS": 7,
  "MinQPS": 1,
  "MedianQPS": 4,
  "Duration": {
    "Count": 18,
    "Histogram": [
      {
        "1ms - 1ms": 18
      }
    ],
    "Rate": {
      "Second": 1000
    },
    "Samples": 18,
    "Time": {
      "Avg": "1ms",
      "Cumulative": "18ms",
      "HMean": "999.999µs",
      "Long5p": "1ms",
      "Max": "1ms",
      "Min": "1ms",
      "P50": "1ms",
      "P75": "1ms",
      "P95": "1ms",
      "P99": "1ms",
      "P999": "1ms",
      "Range": "0s",
      "Short5p": "1ms",
      "StdDev": "0s"
    }
  }
}`, report.JSON())
}
