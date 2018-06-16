package remote

import (
	"github.com/google/go-github/github"
	"fmt"
	"errors"
)

type RemoteError struct {
	Err error
	Status string
	StatusCode int
	Msg string
}

func (r RemoteError) String() string {
	return fmt.Sprintf("stat=%s code=%d msg=%s err=%v", r.Status, r.StatusCode, r.Msg, r.Err)
}

func wrapGHRespErr(res *github.Response, err error) (ok bool, ee error){
	if err != nil {
		if res != nil {
			if res.StatusCode < 200 ||res.StatusCode >= 300 {}
			e := RemoteError{
				Msg:        res.String(),
				Err:        err,
			}
			ee = errors.New(e.String())
		} else {
			ee = errors.New("nil response: err= "+err.Error())
		}
		return
	}
	if res == nil {
		ee = errors.New("nil response from github")
	}
	ok = true
	return
}