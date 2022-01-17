//go:build go1.18
// +build go1.18

package sc

func M[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
