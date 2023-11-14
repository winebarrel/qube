package qube_test

import (
	"os"
	"testing"
)

var (
	testAcc = false
)

const (
	testUUID          = "473d2574-4d1c-46cf-a275-5f3541eb47b7"
	testMySQLURL      = "root@tcp(127.0.0.1:13306)/"
	testPostgreSQLURL = "postgres://postgres@localhost:15432"
)

func TestMain(m *testing.M) {
	if v := os.Getenv("TEST_ACC"); v == "1" {
		testAcc = true
	}

	m.Run()
}
