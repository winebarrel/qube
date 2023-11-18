# qube

[![build](https://github.com/winebarrel/qube/actions/workflows/test.yml/badge.svg)](https://github.com/winebarrel/qube/actions/workflows/test.yml)

qube is a DB load testing tool.

![](https://github.com/winebarrel/qube/assets/117768/d7c42cb0-c3eb-4522-b74c-0b05a6fcc2ed)

## Installation

```sh
brew install winebarrel/qube/qube
```

## Usage

```
Usage: qube --data-file=STRING --dsn=STRING

Flags:
  -h, --help                Show context-sensitive help.
      --[no-]force          Do not abort test on error. (default: disabled)
  -f, --data-file=STRING    NDJSON file path of queries to execute.
      --key="q"             Key name of the query field in the test data. e.g. {"q":"SELECT ..."}
      --[no-]loop           Return to the beginning after reading the test data. (default: enabled)
      --[no-]random         Randomize the starting position of the test data. (default: disabled)
      --commit-rate=INT     Number of queries to execute "COMMIT".
  -d, --dsn=STRING          DSN to connect to.
                              - MySQL: https://github.com/go-sql-driver/mysql#examples
                              - PostgreSQL: https://github.com/jackc/pgx/blob/df5d00e/stdlib/sql.go
      --[no-]noop           No-op mode. No actual query execution. (default: disabled)
  -n, --nagents=1           Number of agents.
  -r, --rate=INT            Rate limit (qps). "0" means unlimited.
  -t, --time=DURATION       Maximum execution time of the test. "0" means unlimited.
      --[no-]progress       Show progress report. (default: enabled)
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

## Testing

```sh
docker compose up -d
make testacc
```

## Tools to convert logs to NDJSON

* MySQL
    * https://github.com/winebarrel/genlog
* PostgreSQL
    * https://github.com/winebarrel/poslog
