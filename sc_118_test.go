//go:build go1.18
// +build go1.18

package sc_test

import (
	"errors"
	"testing"

	. "github.com/egonk/sc"
)

func mFunc(s string, err error) (string, error) {
	return s, err
}

func TestM(t *testing.T) {
	defer expectRecover(t, errors.New("err"))
	M(mFunc("str", errors.New("err")))
}

func TestM_nil(t *testing.T) {
	if m := M(mFunc("str", nil)); m != "str" {
		t.Fatal(m)
	}
}
