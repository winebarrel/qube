# qube

[![build](https://github.com/winebarrel/qube/actions/workflows/test.yml/badge.svg)](https://github.com/winebarrel/qube/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/winebarrel/qube)](https://goreportcard.com/report/github.com/winebarrel/qube)

qube is a DB load testing tool.

![](https://github.com/user-attachments/assets/ad0078d7-ec2d-4976-b0c3-836e05e557db)

## Installation

```sh
brew install winebarrel/qube/qube
```

## Usage

```
Usage: qube --data-files=DATA-FILES,... --dsn=DSN [flags]

Flags:
  -h, --help                         Show help.
      --[no-]force                   Do not abort test on error. (default: disabled)
  -f, --data-files=DATA-FILES,...    JSON Lines file list of queries to execute.
      --key="q"                      Key name of the query field in the test data. e.g. {"q":"SELECT ..."}
      --[no-]loop                    Return to the beginning after reading the test data. (default: enabled)
      --[no-]random                  Randomize the starting position of the test data. (default: disabled)
      --commit-rate=UINT             Number of queries to execute "COMMIT".
  -d, --dsn=DSN                      DSN to connect to. (${...} is replaced by environment variables)
                                       - MySQL: https://pkg.go.dev/github.com/go-sql-driver/mysql#readme-dsn-data-source-name
                                       - PostgreSQL: https://pkg.go.dev/github.com/jackc/pgx/v5/stdlib#pkg-overview
      --[no-]noop                    No-op mode. No actual query execution. (default: disabled)
      --[no-]iam-auth                Use RDS IAM authentication.
  -n, --nagents=1                    Number of agents.
  -r, --rate=FLOAT-64                Rate limit (qps). "0" means unlimited.
  -t, --time=DURATION                Maximum execution time of the test. "0" means unlimited.
      --[no-]progress                Show progress report.
  -C, --[no-]color                   Color report JSON.
      --version
```

```
$ echo '{"q":"select 1"}' > data.jsonl
$ echo '{"q":"select 2"}' >> data.jsonl
$ echo '{"q":"select 3"}' >> data.jsonl

$ qube -d 'root@tcp(127.0.0.1:13306)/' -f data.jsonl -n 5 -t 10s
00:05 | 5 agents / exec 95788 queries, 0 errors (23637 qps)
...
{
  "ID": "b1e23c00-1601-46eb-ad2b-fdf01154243d",
  "StartedAt": "2023-11-12T12:08:29.296154+09:00",
  "FinishedAt": "2023-11-12T12:08:39.297268+09:00",
  "ElapsedTime": "10.001173875s",
  "Options": {
    "Force": false,
    "DataFile": [
      "data.jsonl"
    ],
    "Key": "q",
    "Loop": true,
    "Random": false,
    "CommitRate": 0,
    "DSN": "root@tcp(127.0.0.1:13306)/",
    "Driver": "mysql",
    "Noop": false,
    "IAMAuth": false,
    "Nagents": 5,
    "Rate": 0,
    "Time": "10s"
  },
  "GOMAXPROCS": 10,
  "QueryCount": 238001,
  "ErrorQueryCount": 0,
  "AvgQPS": 23797,
  "MaxQPS": 24977,
  "MinQPS": 21623,
  "MedianQPS": 24051.5,
  "Duration": {
    "Time": {
      "Cumulative": "49.569869935s",
      "HMean": "200.366µs",
      "Avg": "208.275µs",
      "P50": "199.75µs",
      "P75": "222.042µs",
      "P95": "288.875µs",
      "P99": "363.375µs",
      "P999": "594.208µs",
      "Long5p": "349.679µs",
      "Short5p": "142.483µs",
      "Max": "2.796209ms",
      "Min": "98.709µs",
      "Range": "2.6975ms",
      "StdDev": "54.681µs"
    },
    "Rate": {
      "Second": 4801.323875008872
    },
    "Samples": 238001,
    "Count": 238001,
    "Histogram": [
      {
        "98µs - 368µs": 235807
      },
      {
        "368µs - 638µs": 2008
      },
      {
        "638µs - 907µs": 117
      },
      {
        "907µs - 1.177ms": 19
      },
      {
        "1.177ms - 1.447ms": 9
      },
      {
        "1.447ms - 1.717ms": 4
      },
      {
        "1.717ms - 1.986ms": 2
      },
      {
        "1.986ms - 2.256ms": 3
      },
      {
        "2.256ms - 2.526ms": 12
      },
      {
        "2.526ms - 2.796ms": 20
      }
    ]
  }
}
```

### Comment out in the data file

Lines starting with `//` are ignored as comments.

```
{"q":"select 1"}
//comment
//{"q":"select 2"}
{"q":"select 3"}
```

### Use Environment Variables in DSN

```
$ export PASSWORD=mypass
$ qube -d 'root:${PASSWORD}@tcp(127.0.0.1:13306)/' -f data.jsonl -n 5 -t 10s
```

## Test

```sh
docker compose up -d
make testacc
```

## Tools to convert logs to JSON Lines

- MySQL
  - https://github.com/winebarrel/genlog
- PostgreSQL
  - https://github.com/winebarrel/poslog
