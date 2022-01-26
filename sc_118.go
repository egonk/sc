//go:build go1.18
// +build go1.18

package sc

// M panics if err is not nil. v is returned if err is nil. M is a shorthand for
// Must.
//  f := M(os.Create("example"))
//  defer C(f)
//  // write to f
func M[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
