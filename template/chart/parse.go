package chart

import (
	"github.com/ionous/errutil"
)

// Parse is the main function of chart.
func Parse(try State, str string) (err error) {
	if end, last := innerParse(try, str); end == 0 {
		last.NewRune(eof) // FIX: this is odd, i know.
	} else {
		err = endpointError{str, end, last}
	}
	return
}

type endpointError struct {
	str  string
	end  int
	last State
}

func (e endpointError) Error() (ret string) {
	return errutil.Fmt("parsing `%s` ended in %T(%s) at %q(%d)",
		e.str, e.last, e.last.StateName(), e.str[e.end-1], e.end).Error()
}

// if the state machine ends before the string is empty
// return the one-index point of failure; 0 therefore is "ok".
func innerParse(try State, str string) (ret int, last State) {
	for i, r := range str {
		if next := try.NewRune(r); next != nil {
			try = next
		} else {
			ret = i + 1
			break
		}
	}
	last = try
	return
}
