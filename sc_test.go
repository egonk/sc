package sc_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	. "github.com/egonk/sc"
)

func expectRecover(t *testing.T, expect interface{}) {
	if r := recover(); !reflect.DeepEqual(r, expect) {
		t.Fatal(r)
	}
}

type closer struct {
	err error
}

func (c closer) Close() error { return c.err }

func captureOutput(f **os.File, fn func()) (output []byte) {
	tmp, err := os.CreateTemp("", "")
	E(err)
	defer func() {
		output, err = ioutil.ReadFile(tmp.Name())
		E(err)
		E(os.Remove(tmp.Name()))
	}()
	defer C(tmp)
	orig := *f
	defer func() { *f = orig }()
	*f = tmp
	fn()
	return
}

func TestC(t *testing.T) {
	defer expectRecover(t, errors.New("err"))
	C(closer{errors.New("err")})
}

func TestC_nil(t *testing.T) {
	C(closer{})
}

func TestE(t *testing.T) {
	defer expectRecover(t, errors.New("err"))
	E(errors.New("err"))
}

func TestE_nil(t *testing.T) {
	E(nil)
}

func TestP(t *testing.T) {
	output := captureOutput(&os.Stdout, func() {
		P("%v: %v", "write to", "stdout")
	})
	if !bytes.Equal(output, []byte(`write to: stdout`)) {
		t.Fatal(output)
	}
}

func TestPE(t *testing.T) {
	output := captureOutput(&os.Stderr, func() {
		PE("%v: %v", "write to", "stderr")
	})
	if !bytes.Equal(output, []byte(`write to: stderr`)) {
		t.Fatal(output)
	}
}

func TestT(t *testing.T) {
	defer expectRecover(t, errors.New("err: 123"))
	T("%v: %v", "err", "123")
}
