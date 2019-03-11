// go_snitcher Copyright © 2019 Chris Wojno. All rights reserved.

package snitcher

import (
	reopen "github.com/wojnosystems/go_reopen"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestFile_WriteRead(t *testing.T) {
	calledRead := false
	calledWrite := false

	f, err := ioutil.TempFile("", "go-test-snitch-*")
	if err != nil {
		t.Fatal(err)
	}

	rf, err := reopen.OpenFile(f.Name(), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		t.Fatal(err)
	}
	snitch := NewFile(rf, func() {
		calledRead = true
	}, func() {
		calledWrite = true
	}, nil, nil)

	if calledRead {
		t.Error("read should not be called yet")
	}
	if calledWrite {
		t.Error("write should not be called yet")
	}

	n, err := snitch.Write([]byte("boom"))
	if err != nil {
		t.Error(err)
	}
	if n != 4 {
		t.Error("expected to write 4 bytes")
	}

	if calledRead {
		t.Error("read should not be called yet")
	}
	if !calledWrite {
		t.Error("write should be called")
	}

	_, _ = rf.Seek(0, 0)
	data := make([]byte, 10)
	n, err = snitch.Read(data)
	if err != nil && err != io.EOF {
		t.Error(err)
	}
	if n != 4 {
		t.Error("expected to read 4 bytes, but read ", n)
	}

	if !calledRead {
		t.Error("read should be called")
	}
}

func TestFile_ReOpen(t *testing.T) {
	calledReOpen := false

	f, err := ioutil.TempFile("", "go-test-snitch-*")
	if err != nil {
		t.Fatal(err)
	}

	rf, err := reopen.OpenFile(f.Name(), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		t.Fatal(err)
	}
	snitch := NewFile(rf, nil, nil, nil, func() {
		calledReOpen = true
	})

	_ = snitch.ReOpen()

	if !calledReOpen {
		t.Error("expected ReOpen to be called")
	}
}

func TestFile_Close(t *testing.T) {
	calledClose := false

	f, err := ioutil.TempFile("", "go-test-snitch-*")
	if err != nil {
		t.Fatal(err)
	}

	rf, err := reopen.OpenFile(f.Name(), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		t.Fatal(err)
	}
	snitch := NewFile(rf, nil, nil, func() {
		calledClose = true
	}, nil)

	_ = snitch.Close()

	if !calledClose {
		t.Error("expected Close to be called")
	}
}
