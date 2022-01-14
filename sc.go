package sc

import (
	"fmt"
	"io"
	"os"
)

func C(c io.Closer) {
	if err := c.Close(); err != nil {
		panic(err)
	}
}

func E(err error) {
	if err != nil {
		panic(err)
	}
}

func P(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func PE(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func T(format string, a ...interface{}) {
	panic(fmt.Errorf(format, a...))
}
