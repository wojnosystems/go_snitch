// go_snitcher Copyright © 2019 Chris Wojno. All rights reserved.

package snitcher

import (
	reopen "github.com/wojnosystems/go_reopen"
	"io"
)

// Filer Snitch File reports back to a channel when read, write, and closes occur.
// This is a testing library and is not intended for inclusion in production code, but whatever.
type Filer interface {
	io.ReadWriter
	io.Closer
	reopen.ReOpener
}

// file stores the file and the callbacks for each event it snitches on
type file struct {
	obj         Filer
	afterRead   func()
	afterWrite  func()
	afterClose  func()
	afterReOpen func()
}

// NewFile creates a new file snitcher
func NewFile(f Filer, afterRead, afterWrite, afterClose, afterReOpen func()) Filer {
	return &file{
		obj:         f,
		afterRead:   afterRead,
		afterWrite:  afterWrite,
		afterClose:  afterClose,
		afterReOpen: afterReOpen,
	}
}

// Read performs the file read and calls the afterRead callback
func (s *file) Read(d []byte) (int, error) {
	n, err := s.obj.Read(d)
	if s.afterRead != nil {
		s.afterRead()
	}
	return n, err
}

// Write performs the file write and calls the afterWrite callback
func (s *file) Write(d []byte) (int, error) {
	n, err := s.obj.Write(d)
	if s.afterWrite != nil {
		s.afterWrite()
	}
	return n, err
}

// Close performs the file close and calls the afterClose callback
func (s *file) Close() error {
	err := s.obj.Close()
	if s.afterClose != nil {
		s.afterClose()
	}
	return err
}

// ReOpen performs the file reopen and calls the afterReOpen callback
func (s *file) ReOpen() error {
	err := s.obj.ReOpen()
	if s.afterReOpen != nil {
		s.afterReOpen()
	}
	return err
}
