package qube

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"strings"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/winebarrel/qube/rds"
)

type DBDriver string

const (
	DBDriverMySQL      DBDriver = "mysql"
	DBDriverPostgreSQL DBDriver = "pgx"
)

type DBConfig struct {
	DSN       string    `kong:"short='d',required,help='DSN to connect to. \n - MySQL: https://pkg.go.dev/github.com/go-sql-driver/mysql#readme-dsn-data-source-name \n - PostgreSQL: https://pkg.go.dev/github.com/jackc/pgx/v5/stdlib#pkg-overview'"`
	Driver    DBDriver  `kong:"-"`
	Noop      bool      `kong:"negatable,default='false',help='No-op mode. No actual query execution. (default: disabled)'"`
	IAMAuth   bool      `kong:"negatable,default='false',help='Use IAM authentication.'"`
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
	default:
		err = fmt.Errorf("unimplemented driver - %s", config.Driver)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to open DB (%w)", err)
	}

	db := sql.OpenDB(connector)
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
	mycfg, err := mysql.ParseDSN(cfg.DSN)

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
	pgcfg, err := pgx.ParseConfig(cfg.DSN)

	if err != nil {
		return nil, err
	}

	if cfg.IAMAuth {
		host, err := rds.ResolveCNAME(pgcfg.Config.Host)

		if err != nil {
			return nil, err
		}

		endpoint := fmt.Sprintf("%s:%d", host, pgcfg.Config.Port)
		user := pgcfg.Config.User

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
