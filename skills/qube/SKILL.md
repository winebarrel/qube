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

JSON report including:
- Test metadata (ID, timestamps, duration, options)
- Query statistics (total count, error count, avg/min/max/median QPS)
- Latency metrics (percentiles, histogram, HMean, StdDev)
- Runtime info (GOMAXPROCS, version)
