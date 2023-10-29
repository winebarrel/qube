package qube

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	DSN        string `kong:"short='d',required,help='DSN to connect to. see https://github.com/go-sql-driver/mysql#examples'"`
	Noop       bool   `kong:"negatable,default='false',help='No-op mode. No actual query execution. (default: disabled)'"`
	nconns     int
	autoCommit bool
}

func (config *DBConfig) OpenWithPing() (DBIface, error) {
	if config.Noop {
		return &NullDB{os.Stderr}, nil
	}

	db, err := sql.Open("mysql", config.DSN)

	if err != nil {
		return nil, fmt.Errorf("failed to open DB (%w)", err)
	}

	db.SetConnMaxLifetime(0)
	db.SetConnMaxIdleTime(0)
	db.SetMaxIdleConns(0)

	err = db.Ping()

	if err != nil {
		return nil, fmt.Errorf("failed to ping DB (%w)", err)
	}

	if config.autoCommit {
		_, err = db.Exec("set autocommit = 0")

		if err != nil {
			return nil, fmt.Errorf("failed to disable autocommit (%w)", err)
		}
	}

	db.SetMaxIdleConns(config.nconns)
	db.SetMaxOpenConns(config.nconns)

	return db,
		nil
}
