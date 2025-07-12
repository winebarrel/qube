package util

import (
	"io"
	"math/rand"
)

func RandSeek(seeker io.Seeker, size int64) error {
	offset := rand.Int63n(size)
	_, err := seeker.Seek(offset, io.SeekStart)
	return err
}
