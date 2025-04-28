package ptr

// Ptr returns a pointer to the given value.
func Ptr[T any](v T) *T {
	return &v
}

// ValOrZero returns the pointer's value if non-nil, or the zero value of type T if nil.
func ValOrZero[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}
	return *ptr
}

// ValOrDef dereferences ptr and returns the value it points to if non-nil,
// or returns the provided default value if ptr is nil.
func ValOrDef[T any](ptr *T, def T) T {
	if ptr != nil {
		return *ptr
	}
	return def
}

// Equal returns true if both arguments are nil or both arguments
// dereference to the same value.
func Equal[T comparable](a, b *T) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if a == nil {
		return true
	}
	return *a == *b
}

// ToAny converts a value of any type to a pointer to interface{}.
// It first converts the value to interface{} and then returns a pointer to it.
func ToAny[T any](v T) *any {
	a := any(v)
	return &a
}
