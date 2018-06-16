package fmter

import (
	"bytes"
	"io"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
	"github.com/pkg/errors"
	"fmt"
	"github.com/kballard/go-shellquote"
)

type FmtErr struct {
	err error
	msg io.Reader
}
func (e FmtErr) String() string {
	var b bytes.Buffer
	io.ReadFull(e.msg, b.Bytes())
	return fmt.Sprintf("[fmterr] err=%v msg=%s", e.err, b.String())
}
type FmtOut struct {
	msg io.Reader
}
func (o FmtOut) String() string {
	var b bytes.Buffer
	io.ReadFull(o.msg, b.Bytes())
	return fmt.Sprintf("[fmtout] msg=%s", b.String())
}

type FileList struct {
	WhiteList []string
	BlackList []string
}

type FmtConfig struct {
	Commands []string
	Target   string `json:"-",toml:"-",yaml:"-"` // this will get set in the code as the cloned repo dir
	PerFile  bool
	Files    *FileList
	Dirs     *FileList
}

func (f FmtConfig) Print() string {
	return fmt.Sprintf(`[fmtatt] -> %s
$ %s
`, f.Target, strings.Join(f.Commands, "\n- "), )
}

var errSkippingFile = errors.New("skipping file")

func Fmt(fmtConfig FmtConfig, target string) (fouts []FmtOut, ferrs []FmtErr, err error) {
	// this is NOT an asynchronous loop; commands should be run in the order specified
	fmtConfig.Target = target

	for _, c := range fmtConfig.Commands {
		words, e := shellquote.Split(c)
		if e != nil {
			err = e
			return
		}

		command := exec.Command(words[0], words[1:]...)
		command.Dir = fmtConfig.Target
		command.Args = append(command.Args, ".")

		outB, errB := new(bytes.Buffer), new(bytes.Buffer)
		command.Stdout = outB
		command.Stderr = errB

		runCmd := func() {
			if err := command.Run(); err != nil {
				ferrs = append(ferrs, FmtErr{msg: errB, err: err})
			} else {
				fouts = append(fouts, FmtOut{msg: outB})
			}
		}

		if fmtConfig.PerFile {
			fs, readDirErr := ioutil.ReadDir(command.Dir)
			if readDirErr != nil {
				ferrs = append(ferrs, FmtErr{err: readDirErr})
				continue
			}

			d := ""
			for _, f := range fs {
				if f.IsDir() {
					d = f.Name()
				}
				if !filterFileByList(&fmtConfig, d, f.Name()) {
					// TODO maybe add an error or event here
					ferrs = append(ferrs, FmtErr{err: fmt.Errorf("%v: %s", errSkippingFile, f.Name())})
					continue
				}
				command.Args[len(command.Args)-1] = f.Name()
				runCmd()
			}
		} else {
			runCmd()
		}
	}
	return
}

func filterFileByList(config *FmtConfig, dir, ffile string) bool {
	if config.Dirs != nil {
		for _, bd := range config.Dirs.BlackList {
			re := regexp.MustCompile(bd)
			if re.MatchString(dir) {
				return false
			}
		}
		for _, wd := range config.Dirs.WhiteList {
			re := regexp.MustCompile(wd)
			if !re.MatchString(dir) {
				return false
			}
		}
		if dir == ffile {
			return true
		}
	}
	if config.Files != nil {
		for _, bf := range config.Files.BlackList {
			re := regexp.MustCompile(bf)
			if re.MatchString(ffile) {
				return false
			}
		}
		for _, wf := range config.Files.WhiteList {
			re := regexp.MustCompile(wf)
			if !re.MatchString(ffile) {
				return false
			}
		}
	}
	return true
}
