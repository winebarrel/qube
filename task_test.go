package qube_test

import (
	"bytes"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/qube"
)

func TestAcc_Task(t *testing.T) {
	if !testAcc {
		t.Skip()
	}

	t.Setenv("MYSQL_PARAM", "tls=skip-verify")
	t.Setenv("POSTGRES_PARAM", "sslmode=disable")

	assert := assert.New(t)
	require := require.New(t)

	tt := []struct {
		Key     string
		Nagents uint64
		Rate    float64
		Loop    bool
		Random  bool
	}{
		// Default
		{
			Key:     "q",
			Nagents: 1,
			Rate:    0,
			Loop:    true,
			Random:  false,
		},
		// No-loop
		{
			Key:     "q",
			Nagents: 1,
			Rate:    0,
			Loop:    false,
			Random:  false,
		},
		// Randome
		{
			Key:     "q",
			Nagents: 1,
			Rate:    0,
			Loop:    true,
			Random:  true,
		},
		// Mult-agents
		{
			Key:     "q",
			Nagents: 3,
			Rate:    0,
			Loop:    true,
			Random:  false,
		},
		// Non-default key
		{
			Key:     "query",
			Nagents: 1,
			Rate:    0,
			Loop:    true,
			Random:  false,
		},
		// Limit rate
		{
			Key:     "q",
			Nagents: 1,
			Rate:    1,
			Loop:    true,
			Random:  false,
		},
	}

	for _, t := range tt {
		f, _ := os.CreateTemp("", "")
		defer os.Remove(f.Name())
		f.WriteString(`{"` + t.Key + `":"select 1"}` + "\n")
		f.Sync()

		task := &qube.Task{
			Options: &qube.Options{
				AgentOptions: qube.AgentOptions{
					Force: false,
				},
				DataOptions: qube.DataOptions{
					DataFiles:  []string{f.Name()},
					Key:        t.Key,
					Loop:       t.Loop,
					Random:     t.Random,
					CommitRate: 0,
				},
				DBConfig: qube.DBConfig{
					Noop: false,
				},
				Nagents:  t.Nagents,
				Rate:     t.Rate,
				Time:     1 * time.Second,
				Progress: false,
			},
			ID: testUUID,
		}

		ds := []struct {
			DSN    string
			Driver qube.DBDriver
		}{
			{
				DSN:    testDSN_MySQL,
				Driver: qube.DBDriverMySQL,
			},
			{
				DSN:    testDSN_PostgreSQL,
				Driver: qube.DBDriverPostgreSQL,
			},
			{
				DSN:    testDSN_MySQL + "?${MYSQL_PARAM}",
				Driver: qube.DBDriverMySQL,
			},
			{
				DSN:    testDSN_PostgreSQL + "?${POSTGRES_PARAM}",
				Driver: qube.DBDriverPostgreSQL,
			},
		}

		for _, d := range ds {
			task.DSN = qube.DSN(d.DSN)
			task.Driver = d.Driver
			report, err := task.Run()

			require.NoError(err)
			assert.Equal(testUUID, report.ID)
			assert.NotEqual(0, report.QueryCount)
			assert.Equal(0, report.ErrorQueryCount)
			assert.NotEqual(0, report.AvgQPS)
		}
	}
}

func TestAcc_Task_CommitRate(t *testing.T) {
	if !testAcc {
		t.Skip()
	}

	assert := assert.New(t)
	require := require.New(t)

	f, _ := os.CreateTemp("", "")
	defer os.Remove(f.Name())
	f.WriteString(`{"q":"select 1"}` + "\n")
	f.Sync()

	var buf bytes.Buffer

	task := &qube.Task{
		Options: &qube.Options{
			AgentOptions: qube.AgentOptions{
				Force: false,
			},
			DataOptions: qube.DataOptions{
				DataFiles:  []string{f.Name()},
				Key:        "q",
				Loop:       true,
				Random:     false,
				CommitRate: 1,
			},
			DBConfig: qube.DBConfig{
				DSN:       testDSN_MySQL,
				Driver:    qube.DBDriverMySQL,
				Noop:      true,
				NullDBOut: &buf,
			},
			Nagents:  1,
			Rate:     0,
			Time:     1 * time.Second,
			Progress: false,
		},
		ID: testUUID,
	}

	report, err := task.Run()

	require.NoError(err)
	assert.Equal(testUUID, report.ID)
	assert.Regexp("begin", buf.String())
	assert.Regexp("select 1", buf.String())
	assert.Regexp("commit", buf.String())
}

func TestAcc_Task_Force(t *testing.T) {
	if !testAcc {
		t.Skip()
	}

	assert := assert.New(t)
	require := require.New(t)

	f, _ := os.CreateTemp("", "")
	defer os.Remove(f.Name())
	f.WriteString(`{"q":"invalid"}` + "\n")
	f.Sync()

	task := &qube.Task{
		Options: &qube.Options{
			AgentOptions: qube.AgentOptions{
				Force: true,
			},
			DataOptions: qube.DataOptions{
				DataFiles:  []string{f.Name()},
				Key:        "q",
				Loop:       true,
				Random:     false,
				CommitRate: 0,
			},
			DBConfig: qube.DBConfig{
				Noop: false,
			},
			Nagents:  1,
			Rate:     0,
			Time:     1 * time.Second,
			Progress: false,
		},
		ID: testUUID,
	}

	ds := []struct {
		DSN    string
		Driver qube.DBDriver
	}{
		{
			DSN:    testDSN_MySQL,
			Driver: qube.DBDriverMySQL,
		},
		{
			DSN:    testDSN_PostgreSQL,
			Driver: qube.DBDriverPostgreSQL,
		},
	}

	for _, d := range ds {
		task.DSN = qube.DSN(d.DSN)
		task.Driver = d.Driver
		report, err := task.Run()

		require.NoError(err)
		assert.Equal(testUUID, report.ID)
		assert.NotEqual(0, report.QueryCount)
		assert.NotEqual(0, report.ErrorQueryCount)
		assert.Equal(report.QueryCount, report.ErrorQueryCount)
		assert.Equal(float64(0), report.AvgQPS)
	}
}

func TestAcc_Task_MultiData(t *testing.T) {
	if !testAcc {
		t.Skip()
	}

	assert := assert.New(t)
	require := require.New(t)

	f1, _ := os.CreateTemp("", "")
	defer os.Remove(f1.Name())
	f1.WriteString(`{"q":"select 1"}` + "\n")
	f1.Sync()

	f2, _ := os.CreateTemp("", "")
	defer os.Remove(f2.Name())
	f2.WriteString(`{"q":"select 2"}` + "\n")
	f2.Sync()

	task := &qube.Task{
		Options: &qube.Options{
			AgentOptions: qube.AgentOptions{
				Force: false,
			},
			DataOptions: qube.DataOptions{
				DataFiles:  []string{f1.Name(), f2.Name()},
				Key:        "q",
				Loop:       true,
				Random:     false,
				CommitRate: 0,
			},
			DBConfig: qube.DBConfig{
				Noop: false,
			},
			Nagents:  10,
			Rate:     0,
			Time:     3 * time.Second,
			Progress: false,
		},
		ID: testUUID,
	}

	ds := []struct {
		DSN    string
		Driver qube.DBDriver
	}{
		{
			DSN:    testDSN_MySQL,
			Driver: qube.DBDriverMySQL,
		},
		{
			DSN:    testDSN_PostgreSQL,
			Driver: qube.DBDriverPostgreSQL,
		},
	}

	for _, d := range ds {
		task.DSN = qube.DSN(d.DSN)
		task.Driver = d.Driver
		report, err := task.Run()

		require.NoError(err)
		assert.Equal(testUUID, report.ID)
		assert.NotEqual(0, report.QueryCount)
		assert.Equal(0, report.ErrorQueryCount)
		assert.NotEqual(0, report.AvgQPS)
	}
}

func TestAcc_Task_MySQLSession(t *testing.T) {
	if !testAcc {
		t.Skip()
	}

	assert := assert.New(t)
	require := require.New(t)

	db, err := sql.Open("mysql", testDSN_MySQL)
	require.NoError(err)
	_, err = db.Exec("create table qube_test (data text not null)")
	require.NoError(err)

	t.Cleanup(func() {
		_, err := db.Exec("drop table qube_test")
		require.NoError(err)
		db.Close()
	})

	f, _ := os.CreateTemp("", "")
	defer os.Remove(f.Name())
	f.WriteString(`
{"q":"set @sess_data = 'ABC'"}
{"q":"insert into qube_test values (@sess_data)"}
`)
	f.Sync()

	task := &qube.Task{
		Options: &qube.Options{
			AgentOptions: qube.AgentOptions{
				Force: false,
			},
			DataOptions: qube.DataOptions{
				DataFiles:  []string{f.Name()},
				Key:        "q",
				Loop:       true,
				Random:     false,
				CommitRate: 0,
			},
			DBConfig: qube.DBConfig{
				DSN:    testDSN_MySQL,
				Driver: qube.DBDriverMySQL,
				Noop:   false,
			},
			Nagents:  10,
			Rate:     0,
			Time:     3 * time.Second,
			Progress: false,
		},
		ID: testUUID,
	}

	report, err := task.Run()

	require.NoError(err)
	assert.Equal(testUUID, report.ID)
	assert.NotEqual(0, report.QueryCount)
	assert.Equal(0, report.ErrorQueryCount)
	assert.NotEqual(0, report.AvgQPS)

	var data string
	err = db.QueryRow("select data from qube_test limit 1").Scan(&data)
	require.NoError(err)
	assert.Equal("ABC", data)
}

func TestAcc_Task_PostgreSQLSession(t *testing.T) {
	if !testAcc {
		t.Skip()
	}

	assert := assert.New(t)
	require := require.New(t)

	db, err := sql.Open("pgx", testDSN_PostgreSQL)
	require.NoError(err)
	_, err = db.Exec("create table qube_test (data text not null)")
	require.NoError(err)

	t.Cleanup(func() {
		_, err := db.Exec("drop table qube_test")
		require.NoError(err)
		db.Close()
	})

	f, _ := os.CreateTemp("", "")
	defer os.Remove(f.Name())
	f.WriteString(`
{"q":"select set_config('my.sess_data', 'ABC', false)"}
{"q":"insert into qube_test values (current_setting('my.sess_data'))"}
`)
	f.Sync()

	task := &qube.Task{
		Options: &qube.Options{
			AgentOptions: qube.AgentOptions{
				Force: false,
			},
			DataOptions: qube.DataOptions{
				DataFiles:  []string{f.Name()},
				Key:        "q",
				Loop:       true,
				Random:     false,
				CommitRate: 0,
			},
			DBConfig: qube.DBConfig{
				DSN:    testDSN_PostgreSQL,
				Driver: qube.DBDriverPostgreSQL,
				Noop:   false,
			},
			Nagents:  10,
			Rate:     0,
			Time:     3 * time.Second,
			Progress: false,
		},
		ID: testUUID,
	}

	report, err := task.Run()

	require.NoError(err)
	assert.Equal(testUUID, report.ID)
	assert.NotEqual(0, report.QueryCount)
	assert.Equal(0, report.ErrorQueryCount)
	assert.NotEqual(0, report.AvgQPS)

	var data string
	err = db.QueryRow("select data from qube_test limit 1").Scan(&data)
	require.NoError(err)
	assert.Equal("ABC", data)
}

func TestAcc_Task_MySQLTx(t *testing.T) {
	if !testAcc {
		t.Skip()
	}

	assert := assert.New(t)
	require := require.New(t)

	db, err := sql.Open("mysql", testDSN_MySQL)
	require.NoError(err)
	_, err = db.Exec("create table qube_test (data text not null)")
	require.NoError(err)

	t.Cleanup(func() {
		_, err := db.Exec("drop table qube_test")
		require.NoError(err)
		db.Close()
	})

	f, _ := os.CreateTemp("", "")
	defer os.Remove(f.Name())
	f.WriteString(`
{"q":"begin"}
{"q":"insert into qube_test values ('ABC')"}
{"q":"rollback"}
`)
	f.Sync()

	task := &qube.Task{
		Options: &qube.Options{
			AgentOptions: qube.AgentOptions{
				Force: false,
			},
			DataOptions: qube.DataOptions{
				DataFiles:  []string{f.Name()},
				Key:        "q",
				Loop:       true,
				Random:     false,
				CommitRate: 0,
			},
			DBConfig: qube.DBConfig{
				DSN:    testDSN_MySQL,
				Driver: qube.DBDriverMySQL,
				Noop:   false,
			},
			Nagents:  10,
			Rate:     0,
			Time:     3 * time.Second,
			Progress: false,
		},
		ID: testUUID,
	}

	report, err := task.Run()

	require.NoError(err)
	assert.Equal(testUUID, report.ID)
	assert.NotEqual(0, report.QueryCount)
	assert.Equal(0, report.ErrorQueryCount)
	assert.NotEqual(0, report.AvgQPS)

	var cnt int
	err = db.QueryRow("select count(*) from qube_test").Scan(&cnt)
	require.NoError(err)
	assert.Equal(0, cnt)
}

func TestAcc_Task_PostgreSQLTx(t *testing.T) {
	if !testAcc {
		t.Skip()
	}

	assert := assert.New(t)
	require := require.New(t)

	db, err := sql.Open("pgx", testDSN_PostgreSQL)
	require.NoError(err)
	_, err = db.Exec("create table qube_test (data text not null)")
	require.NoError(err)

	t.Cleanup(func() {
		_, err := db.Exec("drop table qube_test")
		require.NoError(err)
		db.Close()
	})

	f, _ := os.CreateTemp("", "")
	defer os.Remove(f.Name())
	f.WriteString(`
{"q":"begin"}
{"q":"insert into qube_test values ('ABC')"}
{"q":"rollback"}
`)
	f.Sync()

	task := &qube.Task{
		Options: &qube.Options{
			AgentOptions: qube.AgentOptions{
				Force: false,
			},
			DataOptions: qube.DataOptions{
				DataFiles:  []string{f.Name()},
				Key:        "q",
				Loop:       true,
				Random:     false,
				CommitRate: 0,
			},
			DBConfig: qube.DBConfig{
				DSN:    testDSN_PostgreSQL,
				Driver: qube.DBDriverPostgreSQL,
				Noop:   false,
			},
			Nagents:  10,
			Rate:     0,
			Time:     3 * time.Second,
			Progress: false,
		},
		ID: testUUID,
	}

	report, err := task.Run()

	require.NoError(err)
	assert.Equal(testUUID, report.ID)
	assert.NotEqual(0, report.QueryCount)
	assert.Equal(0, report.ErrorQueryCount)
	assert.NotEqual(0, report.AvgQPS)

	var cnt int
	err = db.QueryRow("select count(*) from qube_test").Scan(&cnt)
	require.NoError(err)
	assert.Equal(0, cnt)
}
