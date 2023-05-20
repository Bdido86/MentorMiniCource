package src

import (
	"bytes"
	"io"
	"strings"
)

type Reader interface {
	Read(p []byte) (int, error)
	ReadAll(bufSize int) (string, error)
	BytesRead() int64
}

type CountingToLowerReaderImpl struct {
	Reader         io.Reader
	TotalBytesRead int64
}

func (cr *CountingToLowerReaderImpl) Read(p []byte) (int, error) {
	i, err := cr.Reader.Read(p)
	if err != nil && err != io.EOF {
		return i, err
	}

	cr.TotalBytesRead += int64(i)
	copy(p, bytes.ToLower(p))

	return i, err
}

func (cr *CountingToLowerReaderImpl) ReadAll(bufSize int) (string, error) {
	var err error
	var ir int
	ss := make([]byte, 0, bufSize)
	for {
		s := make([]byte, bufSize)
		if ir, err = cr.Read(s); err != nil && err != io.EOF {
			return "", err
		}
		s = s[:ir]

		ss = append(ss, s...)
		if err != nil {
			break
		}
	}

	return strings.TrimSpace(string(ss)), nil
}

func (cr *CountingToLowerReaderImpl) BytesRead() int64 {
	return cr.TotalBytesRead
}

func NewCountingReader(r io.Reader) *CountingToLowerReaderImpl {
	return &CountingToLowerReaderImpl{
		Reader: r,
	}
}
