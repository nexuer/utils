package sets

import (
	"errors"
	"fmt"
	"testing"
)

func TestEmptyErrors(t *testing.T) {
	var slice []error
	var err error

	err = NewErrors(slice)
	if err != nil {
		t.Errorf("expected nil, got %#v", err)
	}
}

func TestErrorsWithNil(t *testing.T) {
	var slice []error
	slice = []error{nil}
	var err error

	err = NewErrors(slice)
	if err != nil {
		t.Errorf("expected nil, got %#v", err)
	}

	// Append a non-nil error
	slice = append(slice, fmt.Errorf("err"))
	err = NewErrors(slice)
	if err == nil {
		t.Errorf("expected non-nil")
	}

	if s := err.Error(); s != "err" {
		t.Errorf("expected 'err', got %q", s)
	}

}
func TestJoinError_Error(t *testing.T) {
	tests := []struct {
		input error
		want  string
	}{
		// Plural
		{
			input: NewErrors([]error{fmt.Errorf("abc"), fmt.Errorf("123")}),
			want:  "[abc, 123]",
		},
		// DedupePlural
		{
			input: NewErrors([]error{fmt.Errorf("abc"), fmt.Errorf("abc"), fmt.Errorf("123")}),
			want:  "[abc, 123]",
		},
		// FlattenAndDedupe
		{
			input: NewErrors([]error{fmt.Errorf("abc"), fmt.Errorf("abc"), NewErrors([]error{fmt.Errorf("abc")})}),
			want:  "abc",
		},
		{
			input: NewErrors([]error{fmt.Errorf("abc"), fmt.Errorf("abc"), NewErrors([]error{fmt.Errorf("abc"), fmt.Errorf("123")})}),
			want:  "[abc, 123]",
		},
		{
			input: NewErrors([]error{fmt.Errorf("abc"), fmt.Errorf("abc"), errors.Join(fmt.Errorf("abc"), fmt.Errorf("123"))}),
			want:  "[abc, 123]",
		},
		// Flatten
		{
			input: NewErrors([]error{fmt.Errorf("abc"), fmt.Errorf("abc"), NewErrors([]error{fmt.Errorf("abc"), fmt.Errorf("def"),
				NewErrors([]error{fmt.Errorf("def"), fmt.Errorf("ghi")})}),
			}),
			want: "[abc, def, ghi]",
		},
	}

	for _, tt := range tests {
		if tt.input.Error() != tt.want {
			t.Errorf("got %s, want %s", tt.input.Error(), tt.want)
		}
	}
}

func TestJoinError_Is(t *testing.T) {
	err1 := errors.New("some specific error")
	err2 := fmt.Errorf("wrapped error: %w", err1)
	tests := []struct {
		input  error
		target error
		Is     bool
	}{
		{
			input:  errors.Join(err1, fmt.Errorf("123")),
			target: err1,
			Is:     true,
		},
		{
			input:  errors.Join(err2, fmt.Errorf("123")),
			target: err1,
			Is:     true,
		},
		{
			input:  NewErrors([]error{err1, fmt.Errorf("123")}),
			target: err1,
			Is:     true,
		},
		{
			input:  NewErrors([]error{err2, fmt.Errorf("123")}),
			target: err1,
			Is:     true,
		},
	}

	for _, tt := range tests {
		if errors.Is(tt.input, tt.target) != tt.Is {
			t.Errorf("got %v, want %v", tt.input.Error(), tt.target.Error())
		}
	}
}
