package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"aqwari.net/net/styx/styxproto"
)

// Turn Go types into files

type fakefile struct {
	v      interface{}
	offset int64
	set    func(s string)
}

func (f *fakefile) ReadAt(p []byte, off int64) (int, error) {
	var s string
	if v, ok := f.v.(fmt.Stringer); ok {
		s = v.String()
	} else {
		s = fmt.Sprint(f.v)
	}
	if off > int64(len(s)) {
		return 0, io.EOF
	}
	n := copy(p, s)
	return n, nil
}

func (f *fakefile) WriteAt(p []byte, off int64) (int, error) {
	buf, ok := f.v.(*bytes.Buffer)
	if !ok {
		return 0, errors.New("not supported")
	}
	if off != f.offset {
		return 0, errors.New("no seeking")
	}
	n, err := buf.Write(p)
	f.offset += int64(n)
	return n, err
}

func (f *fakefile) Close() error {
	if f.set != nil {
		f.set(fmt.Sprint(f.v))
	}
	return nil
}

func (f *fakefile) size() int64 {
	if d, ok := f.v.(map[string]interface{}); ok {
		return int64(styxproto.MaxStatLen * len(d))
	}
	return int64(len(fmt.Sprint(f.v)))
}

type stat struct {
	name string
	file *fakefile
}

func (s *stat) Name() string     { return s.name }
func (s *stat) Sys() interface{} { return s.file }

func (s *stat) ModTime() time.Time {
	return time.Now().Truncate(time.Hour)
}

func (s *stat) IsDir() bool {
	return s.Mode().IsDir()
}

func (s *stat) Mode() os.FileMode {
	if _, ok := s.file.v.(map[string]interface{}); ok {
		return os.ModeDir | 0777
	}
	return 0444
}

func (s *stat) Size() int64 {
	return s.file.size()
}

type dir struct {
	c    chan stat
	done chan struct{}
}

func mkdir(val interface{}) *dir {
	c := make(chan stat, 10)
	done := make(chan struct{})
	go func() {
		if m, ok := val.(map[string]interface{}); ok {
		LoopMap:
			for name, v := range m {
				select {
				case c <- stat{name: name, file: &fakefile{v: v}}:
				case <-done:
					break LoopMap
				}
			}
		} else if a, ok := val.([]interface{}); ok {
		LoopArray:
			for i, v := range a {
				name := strconv.Itoa(i)
				select {
				case c <- stat{name: name, file: &fakefile{v: v}}:
				case <-done:
					break LoopArray
				}
			}
		}
		close(c)
	}()
	return &dir{
		c:    c,
		done: done,
	}
}

func (d *dir) Readdir(n int) ([]os.FileInfo, error) {
	var err error
	fi := make([]os.FileInfo, 0, 10)
	for i := 0; i < n; i++ {
		s, ok := <-d.c
		if !ok {
			err = io.EOF
			break
		}
		fi = append(fi, &s)
	}
	return fi, err
}

func (d *dir) Close() error {
	close(d.done)
	return nil
}
