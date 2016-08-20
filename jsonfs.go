// Command jsonfs allows for the consumption and manipulation of a JSON
// object as a file system hierarchy.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"aqwari.net/net/styx"
)

var (
	addr = flag.String("a", ":5640", "Port to listen on")
)

type server struct {
	file map[string]interface{}
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	var srv server
	if f, err := os.Open(flag.Arg(0)); err != nil {
		log.Fatal(err)
	} else {
		d := json.NewDecoder(f)
		if err := d.Decode(&srv.file); err != nil {
			log.Fatal(err)
		}
	}
	var styxServer styx.Server
	styxServer.ErrorLog = log.New(os.Stderr, "", 0)
	styxServer.TraceLog = log.New(os.Stderr, "", 0)
	styxServer.Addr = *addr
	styxServer.Handler = &srv

	log.Fatal(styxServer.ListenAndServe())
}

func walkTo(v interface{}, loc string) (map[string]interface{}, interface{}, bool) {
	cwd := v
	parts := strings.FieldsFunc(loc, func(r rune) bool { return r == '/' })
	var parent map[string]interface{}

	for _, p := range parts {
		m, ok := cwd.(map[string]interface{})
		if !ok {
			return nil, nil, false
		}
		parent = m
		if child, ok := m[p]; !ok {
			return nil, nil, false
		} else {
			cwd = child
		}
	}
	return parent, cwd, true
}

func (srv *server) Serve9P(s *styx.Session) {
	for t := range s.Requests {
		parent, file, ok := walkTo(srv.file, t.Path())
		if !ok {
			t.Rerror("no such file or directory")
			continue
		}
		switch t := t.(type) {
		case styx.Twalk:
			switch file.(type) {
			case map[string]interface{}:
				t.Rwalk(true, os.ModeDir)
			default:
				t.Rwalk(true, 0)
			}
		case styx.Topen:
			switch v := file.(type) {
			case map[string]interface{}:
				t.Ropen(mkdir(v), os.ModeDir)
			default:
				t.Ropen(strings.NewReader(fmt.Sprint(v)), 0)
			}
		case styx.Tstat:
			fi := &stat{name: path.Base(t.Path()), file: &fakefile{v: file}}
			t.Rstat(fi)
		case styx.Tcreate:
			switch v := file.(type) {
			case map[string]interface{}:
				if t.Perm.IsDir() {
					dir := make(map[string]interface{})
					v[t.Name] = dir
					t.Rcreate(mkdir(dir))
				} else {
					v[t.Name] = new(bytes.Buffer)
					t.Rcreate(&fakefile{
						v:   v[t.Name],
						set: func(s string) { v[t.Name] = s },
					})
				}
			default:
				t.Rerror("%s is not a directory", t.Path())
			}
		case styx.Tremove:
			switch v := file.(type) {
			case map[string]interface{}:
				if len(v) > 0 {
					t.Rerror("directory is not empty")
					break
				}
				if parent != nil {
					delete(parent, path.Base(t.Path()))
					t.Rremove()
				} else {
					t.Rerror("permission denied")
				}
			default:
				if parent != nil {
					delete(parent, path.Base(t.Path()))
					t.Rremove()
				} else {
					t.Rerror("permission denied")
				}
			}
		default:
			t.Rerror("not supported")
		}
	}
}
