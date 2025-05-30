package rand

import (
	"math/rand"
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	valid := rng.letterBytes
	for _, l := range []int{0, 1, 2, 10, 123} {
		s := String(l)
		if len(s) != l {
			t.Errorf("expected string of size %d, got %q", l, s)
		}
		for _, c := range s {
			if !strings.ContainsRune(valid, c) {
				t.Errorf("expected valid characters, got %v", c)
			}
		}
	}
}

// Confirm that panic occurs on invalid input.
func TestRangePanic(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("Panic didn't occur!")
		}
	}()
	// Should result in an error...
	Intn(0)
}

func TestIntn(t *testing.T) {
	// 0 is invalid.
	for _, max := range []int{1, 2, 10, 123} {
		inrange := Intn(max)
		if inrange < 0 || inrange > max {
			t.Errorf("%v out of range (0,%v)", inrange, max)
		}
	}
}

func TestPerm(t *testing.T) {
	Seed(5)
	rand.Seed(5)
	for i := 1; i < 20; i++ {
		actual := Perm(i)
		expected := rand.Perm(i)
		for j := 0; j < i; j++ {
			if actual[j] != expected[j] {
				t.Errorf("Perm call result is unexpected")
			}
		}
	}
}

const (
	maxRangeTestCount = 500
	testStringLength  = 32
)

func TestIntnRange(t *testing.T) {
	// 0 is invalid.
	for min, max := range map[int]int{1: 2, 10: 123, 100: 500} {
		for i := 0; i < maxRangeTestCount; i++ {
			inrange := IntnRange(min, max)
			if inrange < min || inrange >= max {
				t.Errorf("%v out of range (%v,%v)", inrange, min, max)
			}
		}
	}
}

func TestInt63nRange(t *testing.T) {
	// 0 is invalid.
	for min, max := range map[int64]int64{1: 2, 10: 123, 100: 500} {
		for i := 0; i < maxRangeTestCount; i++ {
			inrange := Int63nRange(min, max)
			if inrange < min || inrange >= max {
				t.Errorf("%v out of range (%v,%v)", inrange, min, max)
			}
		}
	}
}

func BenchmarkRandomStringGeneration(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	var s string
	for i := 0; i < b.N; i++ {
		s = String(testStringLength)
	}
	b.StopTimer()
	if len(s) == 0 {
		b.Fatal(s)
	}
}
