package sc_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"testing"

	. "github.com/egonk/sc"
)

func expectRecover(t *testing.T, expect interface{}) {
	if r := recover(); !reflect.DeepEqual(r, expect) {
		t.Fatal(r)
	}
}

func expectRecoverErr(t *testing.T, expect string) {
	r := recover()
	if err, _ := r.(error); err == nil || err.Error() != expect {
		t.Fatal(r)
	}
}

type closer struct {
	err error
}

func (c closer) Close() error { return c.err }

func setInput(data []byte, fn func()) {
	tmp, err := ioutil.TempFile("", "")
	E(err)
	_, err = tmp.Write(data)
	E(err)
	_, err = tmp.Seek(0, 0)
	E(err)
	defer func() {
		E(os.Remove(tmp.Name()))
	}()
	defer C(tmp)
	orig := os.Stdin
	defer func() { os.Stdin = orig }()
	os.Stdin = tmp
	fn()
}

func captureOutput(f **os.File, fn func()) (output []byte) {
	tmp, err := ioutil.TempFile("", "")
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

func TestW(t *testing.T) {
	defer expectRecover(t, errors.New("err: 123: 456"))
	defer W("%v: %v", "err", "123")
	T("456")
}

func testX(t *testing.T, stderr bool) {
	var output []byte
	setInput([]byte(`xyz`), func() {
		if stderr {
			output = captureOutput(&os.Stderr, func() {
				switch runtime.GOOS {
				case "windows":
					X("powershell", `foreach ($l in $input) { [Console]::Error.WriteLine($l) }`)
				default:
					X("sh", "-c", `cat >&2`)
				}
			})
		} else {
			output = captureOutput(&os.Stdout, func() {
				switch runtime.GOOS {
				case "windows":
					X("powershell", `$input | write-output`)
				default:
					X("cat")
				}
			})
		}
	})
	suffix := ""
	if runtime.GOOS == "windows" {
		suffix = "\r\n"
	}
	if !bytes.Equal(output, []byte(`xyz`+suffix)) {
		t.Fatalf("%v\n%s", output, output)
	}
}

func TestX(t *testing.T) {
	for _, stderr := range []bool{false, true} {
		t.Run(fmt.Sprintf("stderr=%v", stderr), func(t *testing.T) {
			testX(t, stderr)
		})
	}
}

func TestX_err(t *testing.T) {
	defer expectRecoverErr(t, "exit status 1")
	switch runtime.GOOS {
	case "windows":
		X("powershell", `exit 1`)
	default:
		X("sh", "-c", `exit 1`)
	}
}
