package qube

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/valyala/fastjson"
	"github.com/winebarrel/qube/util"
)

var (
	EOD = errors.New("EOD")
)

type DataOptions struct {
	DataFile   string `kong:"short='f',required,help='NDJSON file path of queries to execute.'"`
	Key        string `kong:"default='q',help='Key name of the query field in the test data. e.g. {\"q\":\"SELECT ...\"}'"`
	Loop       bool   `kong:"negatable,default='true',help='Return to the beginning after reading the test data. (default: enabled)'"`
	Random     bool   `kong:"negatable,default='false',help='Randomize the starting position of the test data. (default: disabled)'"`
	CommitRate int    `kong:"help='Number of queries to execute \"COMMIT\".'"`
}

type Data struct {
	*DataOptions
	file   *os.File
	reader *bufio.Reader
	count  int
}

func NewData(options *Options) (*Data, error) {
	file, err := os.OpenFile(options.DataFile, os.O_RDONLY, 0)

	if err != nil {
		return nil, fmt.Errorf("failed to open test data - %s (%w)", options.DataFile, err)
	}

	if options.Random {
		err = util.RandSeek(file)

		if err != nil {
			return nil, fmt.Errorf("failed to seek test data (%w)", err)
		}
	}

	reader := bufio.NewReader(file)

	if options.Random {
		_, err = util.ReadLine(reader)

		if err == io.EOF {
			_, err = file.Seek(0, io.SeekStart)

			if err != nil {
				return nil, fmt.Errorf("failed to rewind test data (%w)", err)
			}
		} else if err != nil {
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

	if data.CommitRate > 0 && data.count%(data.CommitRate+1) == 0 {
		return "commit", nil
	}

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
		line, err = util.ReadLine(data.reader)
	}

	if err != nil {
		return "", fmt.Errorf("failed to read test data (%w)", err)
	}

	query := fastjson.GetString(line, data.Key)

	if query == "" {
		return "", fmt.Errorf("failed to get query field from '%s'", line)
	}

	return query, nil
}

func (data *Data) Close() error {
	return data.file.Close()
}
