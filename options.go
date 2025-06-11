package qube

import (
	"fmt"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5"
	"github.com/mattn/go-isatty"
)

type Options struct {
	AgentOptions
	DataOptions
	DBConfig
	Nagents  uint64        `kong:"short='n',default='1',help='Number of agents.'"`
	Rate     float64       `kong:"short='r',help='Rate limit (qps). \"0\" means unlimited.'"`
	Time     time.Duration `json:"-" kong:"short='t',help='Maximum execution time of the test. \"0\" means unlimited.'"`
	X_Time   JSONDuration  `json:"Time" kong:"-"` // for report
	Progress bool          `json:"-" kong:"negatable,help='Show progress report.'"`
	Color    bool          `json:"-" kong:"negatable,short='C',help='Color report JSON.'"`
}

// Kong hook
// see https://github.com/alecthomas/kong#hooks-beforereset-beforeresolve-beforeapply-afterapply-and-the-bind-option
func (options *Options) BeforeApply() error {
	options.Color = isatty.IsTerminal(os.Stdout.Fd())
	return nil
}

func (options *Options) AfterApply() error {
	options.X_Time = JSONDuration(options.Time)
	options.NullDBOut = os.Stderr
	options.Progress = isatty.IsTerminal(os.Stderr.Fd())

	if _, err := mysql.ParseDSN(options.DSN.Fill()); err == nil {
		options.Driver = DBDriverMySQL
	} else if _, err := pgx.ParseConfig(options.DSN.Fill()); err == nil {
		options.Driver = DBDriverPostgreSQL
	} else {
		return fmt.Errorf("cannot parse DSN - %s", options.DSN)
	}

	return nil
}
