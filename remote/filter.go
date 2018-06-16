package remote

import (
	"reflect"
	"regexp"
	"time"
	"fmt"
)

type ErrFilteredT struct {
	Reason   string
	Value    interface{}
	Resource string
}

func (e *ErrFilteredT) String() string {
	if e == nil {
		return "" // hm
	}
	return fmt.Sprintf("%s %s %s %v", "filtered", e.Resource, e.Reason, e.Value)
}

func newFilterError(subject string, spec interface{}, val interface{}) *ErrFilteredT {
	return &ErrFilteredT{
		Reason:   getSpecName(spec),
		Value:    val,
		Resource: subject,
	}
}

func filterNSpec(subject string, spec *MatchNSpec, n int) *ErrFilteredT {
	if spec == nil {
		return nil
	}
	nn := int64(n)
	if nn < spec.Min || nn > spec.Max {
		return newFilterError(subject, spec, n)
	}
	return nil
}

func filterTimeSpec(subject string, spec *MatchTimeSpec, t time.Time) *ErrFilteredT {
	if spec == nil {
		return nil
	}
	if t.Before(spec.Earliest) || t.After(spec.Latest) {
		return newFilterError(subject, spec, t)
	}
	return nil
}

func filterTextSpec(subject string, spec *MatchTextSpec, s string) *ErrFilteredT {
	if spec == nil {
		return nil
	}
	// TODO FIXME check s == "", etc... maybe some tests...!
	// blacklist first
	for _, b := range spec.BlackList {
		re := regexp.MustCompile(b)
		if re.MatchString(s) {
			return newFilterError(subject, spec, s)
		}
	}
	for _, w := range spec.WhiteList {
		re := regexp.MustCompile(w)
		if !re.MatchString(s) {
			return newFilterError(subject, spec, s)
		}
	}
	return nil
}

func filterBoolSpec(subject string, spec bool, b bool) bool {
	return spec == b
}

func getSpecName(spec interface{}) (name string) {
	// Just a sneaky way to abstract error val assignment
	val := reflect.Indirect(reflect.ValueOf(spec))
	name = val.Type().Field(0).Name
	return
}
