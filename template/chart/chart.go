package chart

import (
	"fmt"
)

// func Parse(str string) (ret []Block, err error) {
// 	p := newBlocks(newTopBlockParser)
// 	if e := parse(p, str); e != nil {
// 		err = e
// 	} else {
// 		ret = p.GetBlocks()
// 	}
// 	return
// }

type endpointError struct {
	str  string
	end  int
	last State
}

func (e endpointError) Error() string {
	return fmt.Sprintf("parsing of %s ended in %T at %q(%d)",
		e.str, e.last, e.str[e.end-1], e.end)
}

// parse the string, if the state machine ends before the string is empty
// return the one-index point of failure; 0 therefore is "ok".
func parse(try State, str string) (err error) {
	if end, last := innerParse(try, str); end > 0 {
		err = endpointError{str, end, last}
	}
	return
}

func innerParse(try State, str string) (ret int, last State) {
	for i, r := range str {
		if next := try.NewRune(r); next != nil {
			try = next
		} else {
			ret = i + 1
			last = try
			break
		}
	}
	if ret == 0 {
		try.NewRune(eof)
	}
	return
}
