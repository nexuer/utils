package sets

import (
	"strings"
)

type joinError struct {
	errs []error
}

func (e *joinError) visit(f func(err error) bool) bool {
	for _, err := range e.errs {
		switch errAny := err.(type) {
		case *joinError:
			if match := errAny.visit(f); match {
				return match
			}
		case interface{ Unwrap() []error }:
			je := joinError{errs: errAny.Unwrap()}
			if match := je.visit(f); match {
				return match
			}
		default:
			if match := f(err); match {
				return match
			}
		}
	}
	return false
}

func (e *joinError) Error() string {
	if len(e.errs) == 0 {
		return ""
	}
	if len(e.errs) == 1 {
		return e.errs[0].Error()
	}
	var builder strings.Builder
	builder.WriteByte('[')
	set := NewWithSize[string](len(e.errs))
	e.visit(func(err error) bool {
		msg := err.Error()
		if set.Has(msg) {
			return false
		}
		set.Add(msg)
		if set.Len() > 1 {
			builder.WriteString(", ")
		}
		builder.WriteString(msg)
		return false
	})
	if set.Len() == 1 {
		msg, _ := set.PopAny()
		return msg
	}
	builder.WriteByte(']')
	return builder.String()
}

func (e *joinError) Unwrap() []error {
	return e.errs
}

func NewErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	n := 0
	for _, err := range errs {
		if err != nil {
			n++
		}
	}
	if n == 0 {
		return nil
	}
	e := &joinError{
		errs: make([]error, 0, n),
	}
	for _, err := range errs {
		if err != nil {
			e.errs = append(e.errs, err)
		}
	}
	return e
}
