// sc is a convenience package to support writing short automation scripts that
// simply panic on errors. It is similar to shell scripting with "set -e"
// option.
//  package main
//
//  import (
//  	"bufio"
//  	"os"
//  	"regexp"
//
//  	. "github.com/egonk/sc"
//  )
//
//  func main() {
//  	// set -e
//  	// grep "abc" example
//  	re := regexp.MustCompile(`abc`)
//  	f := M(os.Open("example"))
//  	defer C(f)
//  	s := bufio.NewScanner(f)
//  	for s.Scan() {
//  		if re.Match(s.Bytes()) {
//  			M(fmt.Printf("%s\n", s.Bytes()))
//  		}
//  	}
//  	E(s.Err())
//  }
package sc

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// C calls c.Close() and panics on error. C is a shorthand for Close.
//  f := M(os.Create("example"))
//  defer C(f)
//  // write to f
func C(c io.Closer) {
	if err := c.Close(); err != nil {
		panic(err)
	}
}

// E panics if err is not nil. E is a shorthand for Error.
//  E(os.Chdir("example"))
func E(err error) {
	if err != nil {
		panic(err)
	}
}

// P calls fmt.Printf(format, a...). Write errors are ignored. P is a shorthand
// for Printf.
//  P("=== writing to: %v\n", fn)
func P(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

// PE calls fmt.Fprintf(os.Stderr, format, a...). Write errors are ignored. PE
// is a shorthand for Printf to stdErr.
//  PE("=== error writing to: %v: %v\n", fn, err)
func PE(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

// T panics with fmt.Errorf(format, a...). T is a shorthand for Throw inspired
// by other programming languages.
//  T("invalid argument: %v", arg)
func T(format string, a ...interface{}) {
	panic(fmt.Errorf(format, a...))
}

// W panics with fmt.Errorf(format+": %v", append(a, recover())...) if recover()
// is not nil. W is a shorthand for Wrap.
//  defer W("file: %v", fn) // panic: file: <fn>: invalid argument: <arg>
//  T("invalid argument: %v", arg)
func W(format string, a ...interface{}) {
	if r := recover(); r != nil {
		panic(fmt.Errorf(format+": %v", append(a, r)...))
	}
}

// X calls XS(name, arg...).Run() and panics on error. X is a shorthand for
// eXecute.
//  X("git", "status")
func X(name string, arg ...string) {
	if err := XS(name, arg...).Run(); err != nil {
		panic(err)
	}
}

// XS calls exec.Command(name, arg...) and sets up os.Stdin, os.Stdout and
// os.Stderr in the returned exec.Cmd. XS is a shorthand for eXecute with Stdio.
//  c := XS("go", "build", "-v", ".")
//  c.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64")
//  E(c.Run())
func XS(name string, arg ...string) *exec.Cmd {
	c := exec.Command(name, arg...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c
}
