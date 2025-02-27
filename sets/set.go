package sets

import (
	"cmp"
	"fmt"
	"sort"
	"strings"
)

// Empty is public since it is used by some internal API objects for conversions between external
// string arrays and internal sets, and conversion logic requires public types today.
type Empty struct{}

// Set is a set of the same type elements, implemented via map[comparable]struct{} for minimal memory consumption.
type Set[T comparable] map[T]Empty

// New creates a Set from a list of values.
// NOTE: type param must be explicitly instantiated if given items are empty.
func New[T comparable](items ...T) Set[T] {
	ss := make(Set[T], len(items))
	if len(items) > 0 {
		ss.Add(items...)
	}
	return ss
}

// NewWithSize creates a Set from a list of values with size.
func NewWithSize[T comparable](size int, items ...T) Set[T] {
	ss := make(Set[T], size)
	if len(items) > 0 {
		ss.Add(items...)
	}
	return ss
}

// KeySet creates a Set from a keys of a map[comparable](? extends interface{}).
// If the value passed in is not actually a map, this will panic.
func KeySet[T comparable, V any](theMap map[T]V) Set[T] {
	ret := make(Set[T], len(theMap))
	for keyValue := range theMap {
		ret.Add(keyValue)
	}
	return ret
}

// Add adds items to the set.
func (s Set[T]) Add(items ...T) Set[T] {
	for _, item := range items {
		s[item] = Empty{}
	}
	return s
}

// Remove removes all items from the set.
func (s Set[T]) Remove(items ...T) Set[T] {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

// Clear empties the set.
// It is preferable to replace the set with a newly constructed set,
// but not all callers can do that (when there are other references to the map).
func (s Set[T]) Clear() Set[T] {
	clear(s)
	return s
}

// Has returns true if and only if item is contained in the set.
func (s Set[T]) Has(item T) bool {
	_, contained := s[item]
	return contained
}

// HasAny returns true if any items are contained in the set.
func (s Set[T]) HasAny(items ...T) bool {
	for _, item := range items {
		if s.Has(item) {
			return true
		}
	}
	return false
}

// HasAll returns true if and only if all items are contained in the set.
func (s Set[T]) HasAll(items ...T) bool {
	for _, item := range items {
		if !s.Has(item) {
			return false
		}
	}
	return true
}

// IsZero returns true if set does not contain any elements.
func (s Set[T]) IsZero() bool {
	return s.Len() == 0
}

// Clone returns a new set which is a copy of the current set.
func (s Set[T]) Clone() Set[T] {
	result := make(Set[T], len(s))
	for key := range s {
		result.Add(key)
	}
	return result
}

// UnsortedList returns the slice with contents in random order.
func (s Set[T]) UnsortedList() []T {
	res := make([]T, 0, len(s))
	for key := range s {
		res = append(res, key)
	}
	return res
}

// Range all items by fn.
func (s Set[T]) Range(fn func(item T) error) error {
	for key := range s {
		if err := fn(key); err != nil {
			return err
		}
	}
	return nil
}

// Difference returns a set of objects that are not in s2.
// For example:
// s1 = {a1, a2, a3}
// s2 = {a1, a2, a4, a5}
// s1.Difference(s2) = {a3}
// s2.Difference(s1) = {a4, a5}
func (s Set[T]) Difference(s2 Set[T]) Set[T] {
	result := New[T]()
	for key := range s {
		if !s2.Has(key) {
			result.Add(key)
		}
	}
	return result
}

// Len returns the size of the set.
func (s Set[T]) Len() int {
	return len(s)
}

// SymmetricDifference returns a set of elements which are in either of the sets, but not in their intersection.
// For example:
// s1 = {a1, a2, a3}
// s2 = {a1, a2, a4, a5}
// s1.SymmetricDifference(s2) = {a3, a4, a5}
// s2.SymmetricDifference(s1) = {a3, a4, a5}
func (s Set[T]) SymmetricDifference(s2 Set[T]) Set[T] {
	return s.Difference(s2).Union(s2.Difference(s))
}

// Union returns a new set which includes items in either s1 or s2.
// For example:
// s1 = {a1, a2}
// s2 = {a3, a4}
// s1.Union(s2) = {a1, a2, a3, a4}
// s2.Union(s1) = {a1, a2, a3, a4}
func (s Set[T]) Union(s2 Set[T]) Set[T] {
	result := make(Set[T], len(s)+len(s2))
	for key := range s {
		result.Add(key)
	}
	for key := range s2 {
		result.Add(key)
	}
	return result
}

// Intersection returns a new set which includes the item in BOTH s1 and s2
// For example:
// s1 = {a1, a2}
// s2 = {a2, a3}
// s1.Intersection(s2) = {a2}
func (s Set[T]) Intersection(s2 Set[T]) Set[T] {
	var walk, other Set[T]
	result := New[T]()
	if s.Len() < s2.Len() {
		walk = s
		other = s2
	} else {
		walk = s2
		other = s
	}
	for key := range walk {
		if other.Has(key) {
			result.Add(key)
		}
	}
	return result
}

// IsSuperset returns true if and only if s1 is a superset of s2.
func (s Set[T]) IsSuperset(s2 Set[T]) bool {
	for item := range s2 {
		if !s.Has(item) {
			return false
		}
	}
	return true
}

// Equal returns true if and only if s1 is equal (as a set) to s2.
// Two sets are equal if their membership is identical.
// (In practice, this means same elements, order doesn't matter)
func (s Set[T]) Equal(s2 Set[T]) bool {
	return len(s) == len(s2) && s.IsSuperset(s2)
}

// PopAny returns a single element from the set.
func (s Set[T]) PopAny() (T, bool) {
	for key := range s {
		s.Remove(key)
		return key, true
	}
	var zeroValue T
	return zeroValue, false
}

// String returns a string representation of container
func (s Set[T]) String() string {
	items := make([]string, 0, len(s))
	for key := range s {
		items = append(items, fmt.Sprintf("%v", key))
	}
	return fmt.Sprintf("Set{%s}", strings.Join(items, ", "))
}

type sortableSliceOfGeneric[T cmp.Ordered] []T

func (g sortableSliceOfGeneric[T]) Len() int           { return len(g) }
func (g sortableSliceOfGeneric[T]) Less(i, j int) bool { return less[T](g[i], g[j]) }
func (g sortableSliceOfGeneric[T]) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }

// List returns the contents as a sorted T slice.
//
// This is a separate function and not a method because not all types supported
// by Generic are ordered and only those can be sorted.
func List[T cmp.Ordered](s Set[T]) []T {
	res := make(sortableSliceOfGeneric[T], 0, len(s))
	for key := range s {
		res = append(res, key)
	}
	sort.Sort(res)
	return res
}

func less[T cmp.Ordered](lhs, rhs T) bool {
	return lhs < rhs
}
