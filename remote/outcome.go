package remote

import (
	"fmt"
	"time"
)

type OutcomeStatus int

const (
	DBErr   OutcomeStatus = iota // unknown to db
	Invalid                      // filtered out by metadata
	Clean                        // fmting did nothing (repo worktree was clean)
	Valid                        // cleared by metadata for formatting
	Dirty // fmting did something
	Committed
	Forked        // did fork
	PullRequested // did pr
)

type Outcome struct {
	Status           OutcomeStatus
	// 2018/06/15 21:28:38 json: cannot unmarshal object into Go struct field Outcome.Error of type error
	// 2018/06/15 21:28:38 {"Status":6,"Error":{},"FilteredE":null,"Timestamp":"2018-06-15T21:21:13.274234557+02:00","FormattedOutcome":{"CommitHash":"08aef8a43b86e50173b0dd24fbb7f8103717fed2","GitStatus":"M  running.go\nM  running_test.go\n"},"ForkedOutcome":{"GitUrl":"git://github.com/whilei/gorunning.git","HTMLUrl":"https://github.com/whilei/gorunning","CloneUrl":"https://github.com/whilei/gorunning.git"},"PROutcome":null}

	Error            string // gotcha. Error, not Err // nope. fuck it.
	FilteredE        *ErrFilteredT
	Timestamp        time.Time
	FormattedOutcome *FormattedOutcome
	ForkedOutcome    *ForkedOutcome
	PROutcome        *PullRequestT
}

func (o *Outcome) SetErr(e error) {
	if e != nil {
		o.Error = e.Error()
	}
}

type ForkedOutcome struct {
	GitUrl   string
	HTMLUrl  string
	CloneUrl string
}

type FormattedOutcome struct {
	CommitHash string
	GitStatus  string
}

func (s Outcome) String() string {
	status := ""
	switch s.Status {
	case DBErr:
		status = "DBERR"
	case Invalid:
		status = "INVALID"
	case Valid:
		status = "VALID"
	case Clean:
		status = "CLEAN"
	case Dirty:
		status = "DIRTY"
	case Committed:
		status = "COMMITTED"
	case Forked:
		status = "FORKED"
	case PullRequested:
		status = "PULL REQUESTED"
	default:
		status = "NIL"
	}
	// TODO maybe
	// this is not best practice, children

	err := ""
	if s.Error != "" {
		err = s.Error
	} else {
		err = s.FilteredE.String()
	}

	if err == "" {
		err = "_"
	}
	out := fmt.Sprintf("🏷 %s err=%s %s", status, err, s.Timestamp.Round(time.Second).String())
	if s.PROutcome != nil {
		out += s.PROutcome.String()
		return out
	}
	if s.ForkedOutcome != nil {
		out += s.ForkedOutcome.String()
		return out
	}
	if s.FormattedOutcome != nil {
		out += "\n\n" + s.FormattedOutcome.String()
		return out
	}
	return out
}

func (fo *FormattedOutcome) String() string {
	ch := "<->"
	if len(fo.CommitHash) > 0 {
		ch = fo.CommitHash[:9]
	}
	return fmt.Sprintf(`[%s]
%s`, ch, fo.GitStatus)
}

func (fo *ForkedOutcome) String() string {
	return fmt.Sprintf(`⌥ %s`, fo.CloneUrl)
}
