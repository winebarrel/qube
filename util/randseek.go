package util

import (
	"io"
	"math/rand"
	"os"
)

func RandSeek(file *os.File) error {
	fi, err := file.Stat()

	if err != nil {
		return err
	}

	size := fi.Size()
	offset := rand.Int63n(size)
	_, err = file.Seek(offset, io.SeekStart)

	return err
}
