package builder

import (
	"runtime"
	"strconv"
)

type Location struct {
	File string
	Line int
}

func (l Location) String() (ret string) {
	if len(l.File) > 0 {
		ret = l.File + ":" + strconv.Itoa(l.Line)
	} else {
		ret = "???"
	}
	return
}

// Capture the current goroutine Location. Skip the specified stack frames, a 0 identifies this function's caller.
func Capture(skip int) (ret Location) {
	if _, file, line, ok := runtime.Caller(skip + 1); ok {
		ret = Location{file, line} // otherwise zero-value
	}
	return
}
