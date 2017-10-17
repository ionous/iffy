package chart

import (
	"github.com/ionous/errutil"
	"math"
	"testing"
)

type endpointError int

func (e endpointError) Error() string {
	return errutil.New("ended at", int(e)).Error()
}

// parse the string, if the state machine ends before the string is empty
// return the one-index point of failure; 0 therefore is "ok".
func parse(try State, str string) (ret int) {
	for i, r := range str {
		if next := try.NewRune(r); next != nil {
			try = next
		} else {
			ret = i + 1
			break
		}
	}
	if ret == 0 {
		try.NewRune(eof)
	}
	return
}

func TestNum(t *testing.T) {
	// returns point of failure
	run := func(str string) (val float64, err error) {
		var num numParser
		if end := parse(&num, str); end > 0 {
			val = math.NaN()
			err = endpointError(end)
		} else if v, e := num.GetValue(); e != nil {
			err = e
		} else {
			val = v
		}
		return
	}
	tests := []struct {
		input    string
		endpoint int // 0 means okay, -1 incomplete, >0 the one-index of the failure point.
		value    float64
	}{
		{"0.", -1, math.NaN()},
		{".0", 1, math.NaN()},
		{"0.0", 0, 0},
		{"72.40", 0, 72.4},
		{"072.40", 0, 72.4},
		{"2.71828", 0, 2.71828},
		{"1.e+0", 0, 1},
		{"6.67428e-11", 0, 6.67428e-11},
		{"0e6", 0, 0},
		{"1E6", 0, 1e6},
		{".25", 1, math.NaN()},
		{"0.25", 0, 0.25},
		{"0.12345E+5", 0, 0.12345E+5},
		{"0.12345E+", -1, math.NaN()},
		{"0.12345E", -1, math.NaN()},
		{"42", 0, 42},
		{"0600", 0, 600},
		{"0xFACADE", 0, 0xfacade},
		{"uncle", 1, math.NaN()},
		{"0uncle", 2, math.NaN()},
		{"-5", 0, -5},
	}
	// out of range:
	// {"170141183460469231731687303715884105727", 0, 1.7014118346046923e+38},
	for i, test := range tests {
		if v, e := run(test.input); e == nil {
			// no error returned, then our values should match
			if v != test.value {
				t.Fatalf("test %d mismatched value '%s'; expected: %g, got: %g", i,
					test.input, test.value, v)
				break
			}
		} else {
			// error returned, check the expected error
			if test.endpoint > 0 && e != endpointError(test.endpoint) {
				t.Fatalf("test %d mismatched endpoint '%s' at %s", i, test.input, e)
			} else if c, ok := e.(endpointError); ok {
				t.Fatalf("test %d expected value error '%s' %d at %s", i, test.input, c, e)
			}
		}
		break
	}
}
