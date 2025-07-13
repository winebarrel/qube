package qube

import (
	"io"
	"os"

	seekable "github.com/SaveTheRbtz/zstd-seekable-format-go/pkg"
	"github.com/klauspost/compress/zstd"
)

var (
	_ io.ReadSeekCloser = (*ZstdFile)(nil)
)

type ZstdFile struct {
	file    *os.File
	decoder *zstd.Decoder
	reader  seekable.Reader
	size    int64
}

func NewZstdFile(file *os.File) (*ZstdFile, error) {
	decoder, err := zstd.NewReader(nil)

	if err != nil {
		return nil, err
	}

	reader, err := seekable.NewReader(file, decoder)

	if err != nil {
		decoder.Close()
		return nil, err
	}

	size, err := reader.Seek(0, io.SeekEnd)

	if err != nil {
		reader.Close()
		decoder.Close()
		return nil, err
	}

	_, err = reader.Seek(0, io.SeekStart)

	if err != nil {
		reader.Close()
		decoder.Close()
		return nil, err
	}

	zstdFile := &ZstdFile{
		file:    file,
		decoder: decoder,
		reader:  reader,
		size:    size,
	}

	return zstdFile, nil
}

func (zf *ZstdFile) Read(p []byte) (n int, err error) {
	return zf.reader.Read(p)
}

func (zf *ZstdFile) Seek(offset int64, whence int) (int64, error) {
	return zf.reader.Seek(offset, whence)
}

func (zf *ZstdFile) Size() int64 {
	return zf.size
}

func (zf *ZstdFile) Close() error {
	zf.reader.Close()
	zf.decoder.Close()
	zf.file.Close()
	return nil
}
