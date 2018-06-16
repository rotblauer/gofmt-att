package logger

import (
	"log"
	"fmt"
	"runtime"
	"os"
	"path/filepath"
)

var CWD string

func init() {
	cwd, err := os.Getwd()
	if err != nil {
		panic("WD WD DO?")
	}
	CWD = cwd
}

func (s StdLogger) SetLevel(level int) {
	s.level = level
	// log.SetPrefix()
}

func header() (h string) {
	_, file, no, ok := runtime.Caller(2)
	if ok {
		f, e := filepath.Rel(CWD, file)
		if e != nil {
			panic(e)
		}
		h = fmt.Sprintf("%s:%d", f, no)
	}
	return
}

func (s StdLogger) F(args ...interface{}) {
	if s.level <= 0 {
		as := append([]interface{}{header() + " "}, args...)
		log.Println(as...)
		panic("FIXME FATAL")
	}
}
func (s StdLogger) Ff(fmt string, args ...interface{}) {
	if s.level <= 0 {
		log.Printf(header() + "  "+fmt+"\n", args...)
		panic("FIXME FATAL")
	}
}
func (s StdLogger) E(args ...interface{}) {
	if s.level <= 1 {
		as := append([]interface{}{header() + " "}, args...)
		log.Println(as...)
	}
}
func (s StdLogger) Ef(fmt string, args ...interface{}) {
	if s.level <= 1 {
		log.Printf(header() + "  "+fmt+"\n", args...)
	}
}
func (s StdLogger) W(args ...interface{}) {
	if s.level <= 2 {
		as := append([]interface{}{header() + " "}, args...)
		log.Println(as...)
	}
}
func (s StdLogger) Wf(fmt string, args ...interface{}) {
	if s.level <= 2 {
		log.Printf(header() + "  "+fmt+"\n", args...)
	}
}
func (s StdLogger) I(args ...interface{}) {
	if s.level <= 3 {
		as := append([]interface{}{header() + " "}, args...)
		log.Println(as...)
	}
}
func (s StdLogger) If(fmt string, args ...interface{}) {
	if s.level <= 3 {
		log.Printf(header() + "  "+fmt+"\n", args...)
	}
}
func (s StdLogger) D(args ...interface{}) {
	if s.level <= 4 {
		as := append([]interface{}{header() + " "}, args...)
		log.Println(as...)
	}
}
func (s StdLogger) Df(fmt string, args ...interface{}) {
	if s.level <= 4 {
		log.Printf(header() + "  "+fmt+"\n", args...)
	}
}
func (s StdLogger) R(args ...interface{}) {
	if s.level <= 5 {
		as := append([]interface{}{header() + " "}, args...)
		log.Println(as...)
	}
}
func (s StdLogger) Rf(fmt string, args ...interface{}) {
	if s.level <= 5 {
		log.Printf(header() + "  "+fmt+"\n", args...)
	}
}
