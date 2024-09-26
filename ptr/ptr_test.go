package ptr

import (
	"testing"
)

func TestValOrZero(t *testing.T) {
	type T int

	var val T = 1

	out := ValOrZero(&val)
	if out != val {
		t.Errorf("expected %d, got %d", val, out)
	}

	var ip *int

	zero := ValOrZero(ip)
	if zero != 0 {
		t.Errorf("expected %d, got %d", 0, zero)
	}
}

func TestValOrDef(t *testing.T) {
	type T int

	var val, def T = 1, 0

	out := ValOrDef(&val, def)
	if out != val {
		t.Errorf("expected %d, got %d", val, out)
	}

	out = ValOrDef(nil, def)
	if out != def {
		t.Errorf("expected %d, got %d", def, out)
	}
}

func TestEqual(t *testing.T) {
	type T int

	if !Equal[T](nil, nil) {
		t.Errorf("expected true (nil == nil)")
	}
	if !Equal(Ptr(T(123)), Ptr(T(123))) {
		t.Errorf("expected true (val == val)")
	}
	if Equal(nil, Ptr(T(123))) {
		t.Errorf("expected false (nil != val)")
	}
	if Equal(Ptr(T(123)), nil) {
		t.Errorf("expected false (val != nil)")
	}
	if Equal(Ptr(T(123)), Ptr(T(456))) {
		t.Errorf("expected false (val != val)")
	}
}
