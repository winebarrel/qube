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
  -d, --dsn=STRING           DSN to connect to. see https://github.com/go-sql-driver/mysql#examples
      --[no-]noop            No-op mode. No actual query execution. (default: disabled)
  -n, --nagents=1            Number of agents.
  -r, --rate=INT             Rate limit (qps). "0" means unlimited.
  -t, --time=DURATION        Maximum execution time of the test. "0" means unlimited.
      --version
```

```
$  echo '{"q":"select 1"}' > data.jsonl
$  echo '{"q":"select 2"}' >> data.jsonl
$  echo '{"q":"select 3"}' >> data.jsonl

$ qube -d 'root@tcp(127.0.0.1:13306)/' -f data.jsonl -n 5 -t 10s
00:05 | 5 agents / exec 88756 queries, 0 errors (22101 qps)
...
{
  "ID": "bea6f7c0-fd09-46ed-a801-440a380581c6",
  "DSN": "root@tcp(127.0.0.1:13306)/",
  "StartedAt": "2023-11-11T14:56:17.048404+09:00",
  "FinishedAt": "2023-11-11T14:56:27.049823+09:00",
  "ElapsedTime": "10.001450833s",
  "Options": {
    "AbortOnErr": false,
    "DataFile": "data.jsonl",
    "Key": "q",
    "Loop": true,
    "Random": false,
    "CommitRate": 0,
    "DSN": "root@tcp(127.0.0.1:13306)/",
    "Noop": false,
    "Nagents": 5,
    "Rate": 0,
    "Time": "10s"
  },
  "GOMAXPROCS": 10,
  "QueryCount": 221386,
  "ErrorQueryCount": 5,
  "AvgQPS": 22135,
  "MaxQPS": 23252,
  "MinQPS": 21054,
  "MedianQPS": 22307,
  "Duration": {
    "Time": {
      "Cumulative": "49.111531104s",
      "HMean": "0s",
      "Avg": "221.836µs",
      "P50": "217µs",
      "P75": "244.458µs",
      "P95": "295.542µs",
      "P99": "347.75µs",
      "P999": "495.042µs",
      "Long5p": "335.835µs",
      "Short5p": "149.046µs",
      "Max": "2.21125ms",
      "Min": "0s",
      "Range": "2.21125ms",
      "StdDev": "47.517µs"
    },
    "Rate": {
      "Second": 4507.821178109609
    },
    "Samples": 221386,
    "Count": 221386,
    "Histogram": [
      {
        "0s - 221µs": 120362
      },
      {
        "221µs - 442µs": 100638
      },
      {
        "442µs - 663µs": 315
      },
      {
        "663µs - 884µs": 34
      },
      {
        "884µs - 1.105ms": 7
      },
      {
        "1.105ms - 1.326ms": 4
      },
      {
        "1.326ms - 1.547ms": 1
      },
      {
        "1.547ms - 1.769ms": 3
      },
      {
        "1.769ms - 1.99ms": 4
      },
      {
        "1.99ms - 2.211ms": 18
      }
    ]
  }
}
```
