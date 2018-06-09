package fmtatt

import (
	"io"
)

type fmtErrs struct {
	err error
	msg io.Reader
}
type fmted struct {
	msg io.Reader
}
type Fmter interface {
	Gofmt(fmtConfig FmtConfig, fmted chan fmted, errs chan fmtErrs) (done chan bool)
}

type FileList struct {
	WhiteList []string
	BlackList []string
}

type FmtConfig struct {
	Commands []string
	Target string
	PerFile bool
	Files FileList
	Dirs FileList
}