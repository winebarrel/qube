---
name: qube
description: Database load testing tool that measures query performance and throughput against MySQL or PostgreSQL
allowed-tools: Read, Grep
---

# qube

Database load testing tool that measures query performance and throughput against MySQL or PostgreSQL.

## Overview

qube executes a large number of queries concurrently via multiple agents and reports detailed latency/QPS statistics as JSON output.

## Key Features

- **Multi-agent concurrency**: `-n` flag to set parallel agents
- **MySQL & PostgreSQL support**: Auto-detected from DSN
- **JSONL data files**: Queries provided as JSON Lines (supports `.zst` compressed files)
- **Rate limiting**: `-r` flag to cap QPS
- **Time-based testing**: `-t` flag to set test duration
- **RDS IAM authentication**: `--iam-auth` flag
- **No-op mode**: `--noop` to validate data without executing queries
- **Progress reporting**: Real-time stats (elapsed time, CPU, QPS)
- **Detailed reporting**: Percentile latencies (P50/P75/P95/P99/P99.9), QPS stats, histogram

## CLI Usage

```bash
qube -d <DSN> -f <DATA_FILE> [options]
```

### Required Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--dsn` | `-d` | Database connection string |
| `--data-files` | `-f` | JSON Lines file(s) containing queries |

### Optional Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--nagents` | `-n` | 1 | Number of parallel agents |
| `--rate` | `-r` | 0 (unlimited) | Rate limit in QPS |
| `--time` | `-t` | 0 (unlimited) | Max test duration |
| `--key` | | `"q"` | JSON field name for query |
| `--[no-]loop` | | enabled | Restart from beginning after reading all data |
| `--[no-]random` | | disabled | Randomize starting position in data file |
| `--commit-rate` | | 0 | Execute COMMIT every N queries |
| `--[no-]force` | | disabled | Continue on query errors |
| `--[no-]noop` | | disabled | No-op mode |
| `--[no-]iam-auth` | | disabled | Use RDS IAM authentication |
| `--[no-]progress` | | auto | Show progress report |
| `-C, --[no-]color` | | auto | Color JSON output |

## Data File Format

JSON Lines with one query per line:

```jsonl
{"q":"SELECT * FROM users WHERE id = 1"}
{"q":"INSERT INTO logs VALUES (1, 'test')"}
```

- Lines starting with `//` are treated as comments
- DSN supports environment variable substitution: `${VARIABLE_NAME}`
- Zstandard Seekable Format (`.zst`) supported for compression

## Examples

```bash
# 5 agents, 10 second test
qube -d 'root@tcp(127.0.0.1:13306)/' -f data.jsonl -n 5 -t 10s

# Rate limited to 5000 QPS
qube -d 'root@tcp(127.0.0.1:13306)/' -f data.jsonl -n 5 -r 5000

# PostgreSQL
qube -d 'postgres://user@localhost:5432/mydb' -f data.jsonl -n 10

# Validate data file without executing
qube -d 'root@tcp(127.0.0.1:13306)/' -f data.jsonl --noop
```

## Output

JSON report with the following fields:

| Field | Description |
|-------|-------------|
| `ID` | Unique test run ID (UUID) |
| `StartedAt` / `FinishedAt` | Timestamps of test start and end |
| `ElapsedTime` | Total test duration |
| `Options` | Test configuration (Force, DataFiles, Key, Loop, Random, CommitRate, DSN, Driver, Noop, IAMAuth, Nagents, Rate, Time) |
| `GOMAXPROCS` | Go runtime parallelism setting |
| `QueryCount` | Total number of executed queries |
| `ErrorQueryCount` | Number of failed queries |
| `AvgQPS` / `MaxQPS` / `MinQPS` / `MedianQPS` | Queries per second statistics |
| `Duration.Time.Cumulative` | Total cumulative query execution time |
| `Duration.Time.HMean` | Harmonic mean of query duration |
| `Duration.Time.Avg` | Average query duration |
| `Duration.Time.P50` / `P75` / `P95` / `P99` / `P999` | Latency percentiles |
| `Duration.Time.Long5p` / `Short5p` | Average of longest/shortest 5% of queries |
| `Duration.Time.Max` / `Min` / `Range` | Min/max/range of query duration |
| `Duration.Time.StdDev` | Standard deviation of query duration |
| `Duration.Rate.Second` | Throughput rate per second |
| `Duration.Samples` / `Count` | Number of samples used for statistics |
| `Duration.Histogram` | Distribution of query durations across 10 buckets |
| `Version` | qube version |

### Example Output

```json
{
  "ID": "b1e23c00-1601-46eb-ad2b-fdf01154243d",
  "StartedAt": "2023-11-12T12:08:29.296154+09:00",
  "FinishedAt": "2023-11-12T12:08:39.297268+09:00",
  "ElapsedTime": "10.001173875s",
  "Options": {
    "Force": false,
    "DataFiles": ["data.jsonl"],
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
      {"98µs - 368µs": 235807},
      {"368µs - 638µs": 2008},
      {"638µs - 907µs": 117},
      {"907µs - 1.177ms": 19},
      {"1.177ms - 1.447ms": 9},
      {"1.447ms - 1.717ms": 4},
      {"1.717ms - 1.986ms": 2},
      {"1.986ms - 2.256ms": 3},
      {"2.256ms - 2.526ms": 12},
      {"2.526ms - 2.796ms": 20}
    ]
  },
  "Version": "1.7.2"
}
```
