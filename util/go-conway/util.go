package util

import (
	"bufio"
	"fmt"
	"os"

	"github.com/buger/goterm"
)

// WaitForInput pauses until <Enter> is pressed.
func WaitForInput() {
	buf := bufio.NewReader(os.Stdin)
	fmt.Print("\n")
	buf.ReadBytes('\n')
	goterm.MoveCursorUp(1)
	goterm.Flush()
}

// Guard checks an error and calls Fail with the given fmtArgs if it's not nil.
func Guard(err error, fmtArgs ...interface{}) {
	if err != nil {
		if len(fmtArgs) > 0 {
			Fail(fmtArgs[0], fmtArgs[1:]...)
		} else {
			Fail(err)
		}
	}
}

// Fail prints a message with the given Printf args and exits with status 1.
func Fail(msg interface{}, fmtArgs ...interface{}) {
	fullMsg := fmt.Sprintf("conway: error: %s\n", msg)
	fmt.Printf(fullMsg, fmtArgs...)
	os.Exit(1)
}
