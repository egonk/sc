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
//  			P("%s\n", s.Bytes())
//  		}
//  	}
//  	E(s.Err())
//  }
package sc

import (
	"fmt"
	"io"
	"os"
)

// C calls c.Close() and panics on error.
//  f := M(os.Create("example"))
//  defer C(f)
//  // write to f
func C(c io.Closer) {
	if err := c.Close(); err != nil {
		panic(err)
	}
}

// E panics if err is not nil.
//  E(os.Chdir("example"))
func E(err error) {
	if err != nil {
		panic(err)
	}
}

// P calls fmt.Printf(format, a...).
//  P("=== writing to: %v\n", fn)
func P(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

// PE calls fmt.Fprintf(os.Stderr, format, a...).
//  PE("=== error writing to: %v: %v\n", fn, err)
func PE(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

// T panics with fmt.Errorf(format, a...).
//  T("invalid argument: %v", arg)
func T(format string, a ...interface{}) {
	panic(fmt.Errorf(format, a...))
}

// W panics with fmt.Errorf(format+": %v", append(a, recover())...) if recover()
// is not nil.
//  defer W("file: %v", fn) // panic: file: <fn>: invalid argument: <arg>
//  T("invalid argument: %v", arg)
func W(format string, a ...interface{}) {
	if r := recover(); r != nil {
		panic(fmt.Errorf(format+": %v", append(a, r)...))
	}
}
