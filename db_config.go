package qube

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/winebarrel/esub"
	"github.com/winebarrel/qube/rds"
)

type DBDriver string

func (driver DBDriver) String() string {
	return string(driver)
}

const (
	DBDriverMySQL      DBDriver = "mysql"
	DBDriverPostgreSQL DBDriver = "pgx"
)

type DSN string

func (dsn DSN) String() string {
	return string(dsn)
}

func (dsn DSN) Fill() string {
	dsnStr := string(dsn)
	out, err := esub.Fill(dsnStr)

	if err != nil {
		return dsnStr
	} else {
		return out
	}
}

type DBConfig struct {
	DSN       DSN       `kong:"short='d',required,help='DSN to connect to. (${...} is replaced by environment variables)\n - MySQL: https://pkg.go.dev/github.com/go-sql-driver/mysql#readme-dsn-data-source-name \n - PostgreSQL: https://pkg.go.dev/github.com/jackc/pgx/v5/stdlib#pkg-overview'"`
	Driver    DBDriver  `kong:"-"`
	Noop      bool      `kong:"negatable,default='false',help='No-op mode. No actual query execution. (default: disabled)'"`
	IAMAuth   bool      `kong:"negatable,default='false',help='Use RDS IAM authentication.'"`
	NullDBOut io.Writer `json:"-" kong:"-"`
}

func (config *DBConfig) OpenDBWithPing(autoCommit bool) (DBIface, error) {
	if config.Noop {
		return &NullDB{config.NullDBOut}, nil
	}

	var connector driver.Connector
	var err error

	switch config.Driver {
	case DBDriverMySQL:
		connector, err = config.getMySQLConnector()
	case DBDriverPostgreSQL:
		connector, err = config.getPostgreSQLConnector()
	}

	var db *sql.DB

	if err == nil {
		if connector != nil {
			db = sql.OpenDB(connector)
		} else {
			db, err = sql.Open(config.Driver.String(), config.DSN.Fill())
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to open DB (%w)", err)
	}

	db.SetConnMaxLifetime(0)
	db.SetConnMaxIdleTime(0)
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)

	err = db.Ping()

	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping DB (%w)", err)
	}

	if config.Driver == DBDriverMySQL {
		if autoCommit {
			_, err = db.Exec("set autocommit = 1")
		} else {
			_, err = db.Exec("set autocommit = 0")
		}

		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to disable autocommit (%w)", err)
		}
	}

	return db, nil
}

func (cfg *DBConfig) getMySQLConnector() (driver.Connector, error) {
	mycfg, err := mysql.ParseDSN(cfg.DSN.Fill())

	if err != nil {
		return nil, err
	}

	if cfg.IAMAuth {
		hostPort := strings.SplitN(mycfg.Addr, ":", 2)
		host, err := rds.ResolveCNAME(hostPort[0])

		if err != nil {
			return nil, err
		}

		port := hostPort[1]
		endpoint := host + ":" + port
		user := mycfg.User

		bc := func(ctx context.Context, mc *mysql.Config) error {
			token, err := rds.BuildIAMAuthToken(ctx, endpoint, user)

			if err != nil {
				return err
			}

			mc.Passwd = token
			return nil
		}

		err = mycfg.Apply(mysql.BeforeConnect(bc))

		if err != nil {
			return nil, err
		}

		mycfg.AllowCleartextPasswords = true

		if mycfg.TLSConfig == "" {
			mycfg.TLSConfig = "preferred"
		}
	}

	return mysql.NewConnector(mycfg)
}

func (cfg *DBConfig) getPostgreSQLConnector() (driver.Connector, error) {
	opts := []stdlib.OptionOpenDB{}
	pgcfg, err := pgx.ParseConfig(cfg.DSN.Fill())

	if err != nil {
		return nil, err
	}

	if cfg.IAMAuth {
		host, err := rds.ResolveCNAME(pgcfg.Host)

		if err != nil {
			return nil, err
		}

		endpoint := fmt.Sprintf("%s:%d", host, pgcfg.Port)
		user := pgcfg.User

		opts = append(opts, stdlib.OptionBeforeConnect(func(ctx context.Context, cc *pgx.ConnConfig) error {
			token, err := rds.BuildIAMAuthToken(ctx, endpoint, user)

			if err != nil {
				return err
			}

			cc.Password = token
			return nil
		}))
	}

	connector := stdlib.GetConnector(*pgcfg, opts...)
	return connector, nil
}
