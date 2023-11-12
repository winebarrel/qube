package qube

import (
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5"
)

type Options struct {
	AgentOptions
	DataOptions
	DBConfig
	Nagents  int           `kong:"short='n',default='1',help='Number of agents.'"`
	Rate     int           `kong:"short='r',help='Rate limit (qps). \"0\" means unlimited.'"`
	Time     time.Duration `json:"-" kong:"short='t',help='Maximum execution time of the test. \"0\" means unlimited.'"`
	X_Time   JSONDuration  `json:"Time" kong:"-"` // for report
	Progress bool          `json:"-" kong:"negatable,default='true',help='Show progress report. (default: enabled)'"`
}

// Kong hook
// see https://github.com/alecthomas/kong#hooks-beforereset-beforeresolve-beforeapply-afterapply-and-the-bind-option
func (options *Options) AfterApply() error {
	options.autoCommit = options.CommitRate == 0
	options.X_Time = JSONDuration(options.Time)

	if _, err := mysql.ParseDSN(options.DSN); err == nil {
		options.Driver = DBDriverMySQL
	} else if _, err := pgx.ParseConfig(options.DSN); err == nil {
		options.Driver = DBDriverPostgreSQL
	} else {
		return fmt.Errorf("cannot parse DSN - %s", options.DSN)
	}

	return nil
}
