package lib

import "fmt"

func Map[T, U any](in []T, transform func(T) U) (out []U) {
	out = make([]U, len(in))
	for i, each := range in {
		out[i] = transform(each)
	}
	return
}

func StringsOf[T fmt.Stringer](in []T) []string {
	return Map(
		in,
		func(t T) string { return t.String() },
	)
}

func IndirectsOf[T any](in []T) []*T {
	return Map(
		in,
		func(t T) *T { return &t },
	)
}

// Drops the first argument passed to this function.
// Useful for ignoring the first return value of a function call.
//
//	// Example:
//	_, err := foo()
//	return err
//
//	// Equivalent to:
//	return Just(foo())
func Just[T any](_ T, err error) error {
	return err
}
