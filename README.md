# qube

[![build](https://github.com/winebarrel/qube/actions/workflows/build.yml/badge.svg)](https://github.com/winebarrel/qube/actions/workflows/build.yml)

qube is a DB load testing tool.

## Usage

```
Usage: qube --data-file=STRING --dsn=STRING

Flags:
  -h, --help                 Show context-sensitive help.
      --[no-]abort-on-err    Abort test on error. (default: disabled)
  -f, --data-file=STRING     NDJSON file path of queries to execute.
      --key="q"              Key name of the query field in the test data. e.g. {"q":"SELECT ..."}
      --[no-]loop            Return to the beginning after reading the test data. (default: enabled)
      --[no-]random          Randomize the starting position of the test data. (default: disabled)
      --commit-rate=INT      Number of queries to execute "COMMIT".
  -d, --dsn=STRING           DSN to connect to.
                               - MySQL: https://github.com/go-sql-driver/mysql#examples
                               - PostgreSQL: https://github.com/jackc/pgx/blob/df5d00e/stdlib/sql.go
      --[no-]noop            No-op mode. No actual query execution. (default: disabled)
  -n, --nagents=1            Number of agents.
  -r, --rate=INT             Rate limit (qps). "0" means unlimited.
  -t, --time=DURATION        Maximum execution time of the test. "0" means unlimited.
      --[no-]progress        Show progress report. (default: enabled)
      --version
```

```
$  echo '{"q":"select 1"}' > data.jsonl
$  echo '{"q":"select 2"}' >> data.jsonl
$  echo '{"q":"select 3"}' >> data.jsonl

$ qube -d 'root@tcp(127.0.0.1:13306)/' -f data.jsonl -n 5 -t 10s
00:07 | 5 agents / exec 147489 queries, 0 errors (24276 qps)
...
{
  "ID": "21111f0e-c40a-465b-999c-734674f57721",
  "StartedAt": "2023-11-12T12:03:33.320236+09:00",
  "FinishedAt": "2023-11-12T12:03:43.321728+09:00",
  "ElapsedTime": "10.001689666s",
  "Options": {
    "AbortOnErr": false,
    "DataFile": "data.jsonl",
    "Key": "q",
    "Loop": true,
    "Random": false,
    "CommitRate": 0,
    "DSN": "root@tcp(127.0.0.1:13306)/",
    "Driver": "mysql",
    "Noop": false,
    "Nagents": 5,
    "Rate": 0,
    "Time": "10s"
  },
  "GOMAXPROCS": 10,
  "QueryCount": 245221,
  "ErrorQueryCount": 0,
  "AvgQPS": 24517,
  "MaxQPS": 25487,
  "MinQPS": 23609,
  "MedianQPS": 24663,
  "Duration": {
    "Time": {
      "Cumulative": "49.578293032s",
      "HMean": "195.952µs",
      "Avg": "202.178µs",
      "P50": "197.125µs",
      "P75": "215.583µs",
      "P95": "264.708µs",
      "P99": "332.125µs",
      "P999": "592.166µs",
      "Long5p": "319.822µs",
      "Short5p": "140.433µs",
      "Max": "2.671708ms",
      "Min": "87.417µs",
      "Range": "2.584291ms",
      "StdDev": "46.378µs"
    },
    "Rate": {
      "Second": 4946.136403722566
    },
    "Samples": 245221,
    "Count": 245221,
    "Histogram": [
      {
        "87µs - 345µs": 243387
      },
      {
        "345µs - 604µs": 1610
      },
      {
        "604µs - 862µs": 146
      },
      {
        "862µs - 1.121ms": 29
      },
      {
        "1.121ms - 1.379ms": 17
      },
      {
        "1.379ms - 1.637ms": 4
      },
      {
        "1.637ms - 1.896ms": 10
      },
      {
        "1.896ms - 2.154ms": 3
      },
      {
        "2.154ms - 2.413ms": 1
      },
      {
        "2.413ms - 2.671ms": 14
      }
    ]
  }
}
```
