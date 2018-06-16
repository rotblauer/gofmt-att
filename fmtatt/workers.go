package fmtatt

import (
	"github.com/rotblauer/gofmt-att/remote"
	"github.com/rotblauer/gofmt-att/fmter"
	"os"
	"strings"
	"regexp"
	"io/ioutil"
	"path/filepath"
)

func (f *FmtAtt) teardown(r *remote.RepoT) error {
	f.Logger.I("teardown", r.String())
	if _, err := os.Stat(r.Target); err == nil {
		f.Logger.D("removing", r.Target)
		if e := os.RemoveAll(r.Target); e != nil {
			return e
		}
		ownerFolder := filepath.Dir(r.Target)
		ds, err := ioutil.ReadDir(ownerFolder)
		if err != nil {
			return err
		}
		if len(ds) == 0 {
			f.Logger.D("removing", ownerFolder)
			return os.RemoveAll(ownerFolder)

		}
	}
	f.repoPool.Remove(r)
	// if only repo from that owner, splice him from owner pool
	if rs := f.repoPool.GetWhere(func(rr *remote.RepoT, o *remote.Outcome) bool {
		return rr.Owner.Name == r.Owner.Name
	}); len(rs) == 0 {
		f.ownerPool.Splice(r.Owner)
	}
	return nil
}

func (f *FmtAtt) fmter(r *remote.RepoT) (err error) {
	// NOTE that these formatters are NOT running asynchronously; they may depend on cardinality
	var allOuts []fmter.FmtOut
	var allErrs []fmter.FmtErr

	for _, ft := range f.Fmters {
		// deref, so can safely assign unique target
		ftft := *ft
		ftft.Target = r.Target
		f.Logger.I("running", ftft.Print())
		outs, errs, e := fmter.Fmt(ftft, r.Target)
		if e != nil {
			err = e
			return
		}
		allOuts = append(allOuts, outs...)
		allErrs = append(allErrs, errs...)
	}
	// FIXME; messy, out of order
	for _, o := range allOuts {
		f.Logger.I(o.String())
	}
	for _, e := range allErrs {
		f.Logger.I(e.String())
	}
	return
}

var fileAddMatchReg = regexp.MustCompile(`M\s*(?P<FILE>.*)\b`)

func (f *FmtAtt) add(r *remote.RepoT, outcome *remote.Outcome, status string) (added int, err error) {
	// collect a set of add-able paths, if params are spec'd
	if addFileSpecs := f.Config.GitConfig.AddPaths; addFileSpecs != nil {
		// collect addable files
		addFiles := []string{}
		compileRegexes := func(list []string)(res []*regexp.Regexp) {
			for _, s := range list {
				res = append(res, regexp.MustCompile(s))
			}
			return
		}
		whites, blacks := compileRegexes(addFileSpecs.WhiteList), compileRegexes(addFileSpecs.BlackList)
	lines:
		for _, line := range strings.Split(status, "\n") { // https://stackoverflow.com/questions/14493867/what-is-the-most-portable-cross-platform-way-to-represent-a-newline-in-go-golang#14494187
			// if 0 whites, will continue past loop
			for _, re := range whites {
				// always add whitelist lines (prioritized)
				// b/c the 'continue'
				if re.MatchString(line) {
					addFiles = append(addFiles, line)
					continue lines
				}
			}
			for _, re := range blacks {
				if re.MatchString(line) {
					// don't add line
					continue lines
				}
			}
			// either 0 specs or no matches in BOTH lists
			// default is add
			addFiles = append(addFiles, line)
		}

		// git add <files>
		for _, af := range addFiles {
			match := fileAddMatchReg.FindStringSubmatch(af)
			if len(match) != 2 {
				continue
			}

			fName := match[1]
			if fName == "" {
				continue
			}

			err = f.Giter.Add(r.Target, fName)
			if err != nil {
				f.Logger.Ef("error adding file: '%s'", fName)
			} else {
				added++
			}
		}

	} else {
		// add eerrything
		added++ // because the repo was _dirty_, so there must have been lines
		err = f.Giter.Add(r.Target, ".")
		if err != nil {
			outcome.SetErr(err)
			return
		}
	}
	return added, nil
}

func (f *FmtAtt) stripeStatus(status string) (repos []*remote.RepoT) {
	for _, stripe := range f.Config.GitConfig.StripeList {
		for _, statusLine := range strings.Split(status, "\n") {
			params := getParams(stripe, statusLine)
			owner, ok := params["OWNER"]
			if !ok {
				continue
			}
			repo, ok := params["REPO"]
			if !ok {
				continue
			}
			repos = append(repos, remote.RepoizeRef(owner+"/"+repo))
		}
	}
	return
}

/**
 * Parses url with the given regular expression and returns the
 * group values defined in the expression.
 *
https://stackoverflow.com/questions/30483652/how-to-get-capturing-group-functionality-in-golang-regular-expressions#30483899

You can use this function like:

params := getParams(`(?P<Year>\d{4})-(?P<Month>\d{2})-(?P<Day>\d{2})`, `2015-05-27`)
fmt.Println(params)
and the output will be:

map[Year:2015 Month:05 Day:27]
 */

// func findAddFileMatch(re *regexp.Regexp, s string) string {
// 	out := ""
// 	match := fileAddMatchReg.FindStringSubmatch(s)
//
//
// 	return out
// }

func getParams(regEx, url string) (paramsMap map[string]string) {

	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(url)

	paramsMap = make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return
}

// 2018/06/13 21:05:03 fmtatt/workers.go:73   [fmtout] msg=
// 	2018/06/13 21:05:03 fmtatt/workers.go:78  fmting finished (3/3) Organization:go-task/task clone=https://github.com/go-task/task.git target=/Users/ia/gofmt-att/clones/go-task/task)
// 2018/06/13 21:05:03 fmtatt/workers.go:110   dirty: /Users/ia/gofmt-att/clones/go-task/task
// M vendor/github.com/Masterminds/sprig/numeric.go
// 	M vendor/github.com/Masterminds/sprig/regex.go
// 	M vendor/github.com/huandu/xstrings/convert.go
// 	M vendor/gopkg.in/yaml.v2/readerc.go
// 	M vendor/gopkg.in/yaml.v2/resolve.go
// 	M vendor/gopkg.in/yaml.v2/sorter.go
//
// 	2018/06/13 21:05:05 fmtatt/workers.go:120  committed: 313154d00
// M  vendor/github.com/Masterminds/sprig/numeric.go
// 	M  vendor/github.com/Masterminds/sprig/regex.go
// 	M  vendor/github.com/huandu/xstrings/convert.go
// 	M  vendor/gopkg.in/yaml.v2/readerc.go
// 	M  vendor/gopkg.in/yaml.v2/resolve.go
// 	M  vendor/gopkg.in/yaml.v2/sorter.go