package qube

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/valyala/fastjson"
	"github.com/winebarrel/qube/util"
)

var (
	// End of data
	EOD = errors.New("EOD")
)

type DataOptions struct {
	DataFiles  []string `kong:"short='f',required,help='JSON Lines file list of queries to execute.'"`
	Key        string   `kong:"default='q',help='Key name of the query field in the test data. e.g. {\"q\":\"SELECT ...\"}'"`
	Loop       bool     `kong:"negatable,default='true',help='Return to the beginning after reading the test data. (default: enabled)'"`
	Random     bool     `kong:"negatable,default='false',help='Randomize the starting position of the test data. (default: disabled)'"`
	CommitRate uint     `kong:"help='Number of queries to execute \"COMMIT\".'"`
}

type Data struct {
	*DataOptions
	file   *os.File
	reader *bufio.Reader
	count  uint
	inTxn  bool
}

func NewData(options *Options, agentNum uint64) (*Data, error) {
	dataFile := options.DataFiles[agentNum%uint64(len(options.DataFiles))]
	file, err := os.OpenFile(dataFile, os.O_RDONLY, 0)

	if err != nil {
		return nil, fmt.Errorf("failed to open test data - %s (%w)", dataFile, err)
	}

	fileInfo, err := file.Stat()

	if err != nil {
		return nil, fmt.Errorf("failed to get test data file info - %s (%w)", dataFile, err)
	}

	if fileInfo.Size() == 0 {
		return nil, fmt.Errorf("test data is empty - %s", dataFile)
	}

	if options.Random {
		err = util.RandSeek(file)

		if err != nil {
			file.Close()
			return nil, fmt.Errorf("failed to seek test data (%w)", err)
		}
	}

	reader := bufio.NewReader(file)

	if options.Random {
		// If it is random, skip one line
		_, err = util.ReadLine(reader)

		if err == io.EOF {
			_, err = file.Seek(0, io.SeekStart)

			if err != nil {
				file.Close()
				return nil, fmt.Errorf("failed to rewind test data (%w)", err)
			}
		} else if err != nil {
			file.Close()
			return nil, fmt.Errorf("failed to read test data (%w)", err)
		}
	}

	data := &Data{
		DataOptions: &options.DataOptions,
		file:        file,
		reader:      reader,
	}

	return data, nil
}

func (data *Data) Next() (string, error) {
	data.count++

	if data.CommitRate > 0 && !data.inTxn {
		data.inTxn = true
		return "begin", nil
	}

	if data.CommitRate > 0 && data.count%(data.CommitRate+2) == 0 {
		data.inTxn = false
		return "commit", nil
	}

	for {
		line, err := util.ReadLine(data.reader)

		if err == io.EOF {
			if !data.Loop {
				return "", EOD
			}

			_, err = data.file.Seek(0, io.SeekStart)

			if err != nil {
				return "", fmt.Errorf("failed to rewind test data (%w)", err)
			}

			data.reader.Reset(data.file)
			continue
		}

		if err != nil {
			return "", fmt.Errorf("failed to read test data (%w)", err)
		}

		if len(line) == 0 {
			continue
		}

		if bytes.HasPrefix(line, []byte("//")) {
			continue
		}

		query := fastjson.GetString(line, data.Key)

		if query == "" {
			return "", fmt.Errorf(`failed to get query field "%s" from '%s'`, data.Key, line)
		}

		return query, nil
	}
}

func (data *Data) Close() error {
	return data.file.Close()
}
