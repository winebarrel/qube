package qube

import (
	"database/sql"
	"fmt"
	"io"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBDriver string

const (
	DBDriverMySQL      DBDriver = "mysql"
	DBDriverPostgreSQL DBDriver = "pgx"
)

type DBConfig struct {
	DSN       string    `kong:"short='d',required,help='DSN to connect to. \n - MySQL: https://github.com/go-sql-driver/mysql#examples \n - PostgreSQL: https://github.com/jackc/pgx/blob/df5d00e/stdlib/sql.go'"`
	Driver    DBDriver  `kong:"-"`
	Noop      bool      `kong:"negatable,default='false',help='No-op mode. No actual query execution. (default: disabled)'"`
	NullDBOut io.Writer `kong:"-"`
}

func (config *DBConfig) OpenDBWithPing(autoCommit bool) (DBIface, error) {
	if config.Noop {
		return &NullDB{config.NullDBOut}, nil
	}

	db, err := sql.Open(string(config.Driver), config.DSN)

	if err != nil {
		return nil, fmt.Errorf("failed to open DB (%w)", err)
	}

	db.SetConnMaxLifetime(0)
	db.SetConnMaxIdleTime(0)
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)

	err = db.Ping()

	if err != nil {
		return nil, fmt.Errorf("failed to ping DB (%w)", err)
	}

	if config.Driver == DBDriverMySQL {
		if autoCommit {
			_, err = db.Exec("set autocommit = 1")
		} else {
			_, err = db.Exec("set autocommit = 0")
		}

		if err != nil {
			return nil, fmt.Errorf("failed to disable autocommit (%w)", err)
		}
	}

	return db,
		nil
}
