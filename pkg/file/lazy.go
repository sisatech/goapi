package file

import (
	"errors"
	"io"
	"os"
)

// LazyReadCloser is an implementation of io.ReadCloser
// that defers its own initialization until the first
// attempted read.
func LazyReadCloser(openFunc func() (io.Reader, error),
	closeFunc func() error) io.ReadCloser {
	return &lazyReadCloser{
		openFunc:  openFunc,
		closeFunc: closeFunc,
	}
}

type lazyReadCloser struct {
	opened    bool
	closed    bool
	r         io.Reader
	openFunc  func() (io.Reader, error)
	closeFunc func() error
}

func (rc *lazyReadCloser) Read(p []byte) (n int, err error) {
	if rc.closed {
		err = errors.New("lazy readcloser is closed")
		return
	}

	if rc.r == nil {
		rc.r, err = rc.openFunc()
		if err != nil {
			return
		}
		rc.opened = true
	}

	return rc.r.Read(p)
}

func (rc *lazyReadCloser) Close() error {
	if rc.closed {
		return errors.New("lazy readcloser already closed")
	}
	rc.closed = true
	return rc.closeFunc()
}

// LazyOpen is an alternative implementation of Open that
// defers actually opening the file until the first
// attempted read.
func LazyOpen(path string) (File, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	var f *os.File

	openFunc := func() (io.Reader, error) {
		f, err = os.Open(path)
		if err != nil {
			return nil, err
		}
		return f, nil
	}

	closeFunc := func() error {
		if f != nil {
			return f.Close()
		}
		return nil
	}

	return CustomFile(CustomFileArgs{
		Name:       fi.Name(),
		Size:       int(fi.Size()),
		ModTime:    fi.ModTime(),
		IsDir:      fi.IsDir(),
		ReadCloser: LazyReadCloser(openFunc, closeFunc),
	}), nil
}
