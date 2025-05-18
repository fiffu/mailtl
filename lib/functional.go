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

// DropError suppresses the error argument passed to this function.
// Useful for ignoring errors(?), use with care!
func DropError[T any](val T, _ error) T {
	return val
}

// DropResult suppresses the first argument passed to this function.
//
//	// Example:
//	_, err := foo()
//	return err
//
//	// Equivalent to:
//	return DropResult(foo())
func DropResult[T any](_ T, err error) error {
	return err
}
