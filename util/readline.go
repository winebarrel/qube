package util

import (
	"bufio"
	"bytes"
)

func ReadLine(r *bufio.Reader) ([]byte, error) {
	var buf bytes.Buffer

	for {
		line, isPrefix, err := r.ReadLine()
		n := len(line)

		if n > 0 {
			buf.Write(line)
		}

		if !isPrefix || err != nil {
			return buf.Bytes(), err
		}
	}
}
