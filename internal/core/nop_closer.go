package core

import (
	"io"
)

type NopCloser struct {
	io.Reader
}

func (NopCloser) Close() error { return nil }
