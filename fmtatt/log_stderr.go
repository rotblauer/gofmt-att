package fmtatt

import (
	"log"
)

type StdLogger struct {
	level int
}

func (s StdLogger) SetLevel(level int) {
	s.level = level
}

func (s StdLogger) F(args []interface{}) {
	if s.level <= 0 {
		log.Println(args)
	}
}
func (s StdLogger) Ff(fmt string, args []interface{}) {
	if s.level <= 0 {
		log.Printf(fmt+"\n", args)
	}
}
func (s StdLogger) E(args []interface{}) {
	if s.level <= 1 {
		log.Println(args)
	}
}
func (s StdLogger) Ef(fmt string, args []interface{}) {
	if s.level <= 1 {
		log.Printf(fmt+"\n", args)
	}
}
func (s StdLogger) W(args []interface{}) {
	if s.level <= 2 {
		log.Println(args)
	}
}
func (s StdLogger) Wf(fmt string, args []interface{}) {
	if s.level <= 2 {
		log.Printf(fmt+"\n", args)
	}
}
func (s StdLogger) I(args []interface{}) {
	if s.level <= 3 {
		log.Println(args)
	}
}
func (s StdLogger) If(fmt string, args []interface{}) {
	if s.level <= 3  {
		log.Printf(fmt+"\n", args)
	}
}
func (s StdLogger) D(args []interface{}) {
	if s.level <= 4 {
		log.Println(args)
	}
}
func (s StdLogger) Df(fmt string, args []interface{}) {
	if s.level <= 4 {
		log.Printf(fmt+"\n", args)
	}
}
func (s StdLogger) R(args []interface{}) {
	if s.level <= 5 {
		log.Println(args)
	}
}
func (s StdLogger) Rf(fmt string, args []interface{}) {
	if s.level <= 5 {
		log.Printf(fmt+"\n", args)
	}
}

