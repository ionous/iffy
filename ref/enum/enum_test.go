package enum

import (
	"fmt"
	"github.com/ionous/iffy/ref/unique"
	"github.com/stretchr/testify/assert"
	r "reflect"
	"strings"
	"testing"
)

// TestEnumChoices verifies a stringerified enum generates good choices.
func TestEnumChoices(t *testing.T) {
	assert := assert.New(t)
	reg := make(Registry)
	if choices, e := reg.Register((*YesNoMaybe)(nil), _YesNoMaybe_name, _YesNoMaybe_index[:]); assert.NoError(e, "enum should generate") {
		//
		var reduce []string
		for _, c := range choices {
			reduce = append(reduce, c)
		}
		choiceToIndex := func(s string) int {
			ret := -1
			for i, c := range choices {
				if s == c {
					ret = i
					break
				}
			}
			return ret
		}
		assert.EqualValues("NoYesMaybe", strings.Join(reduce, ""))
		assert.EqualValues(0, choiceToIndex("No"))
		assert.EqualValues(-1, choiceToIndex("Never"))
		assert.EqualValues("Maybe", choices[2])
	}
}

// TestEnumEmpty verifies an empty enum is an error.
func TestEnumEmpty(t *testing.T) {
	if tp, e := unique.TypePtr(r.Int, (*EmptyState)(nil)); e != nil {
		t.Fatal("couldnt get type pointer for TooLongState", e)
	} else if _, e := Stringify(tp); e == nil {
		t.Fatal("enum shouldnt generate")
	}
}

// TestEnumTooLong verifies an enum with too many choices is an error.
func TestEnumTooLong(t *testing.T) {
	if tp, e := unique.TypePtr(r.Int, (*TooLongState)(nil)); e != nil {
		t.Fatal("couldnt get type pointer for TooLongState", e)
	} else if _, e := Stringify(tp); e == nil {
		t.Fatal("enum shouldnt generate")
	}
}

// YesNoMaybe provides an enum with three choices for testing.
//go:generate stringer -type=YesNoMaybe
type YesNoMaybe int

const (
	No YesNoMaybe = iota
	Yes
	Maybe
)

const _YesNoMaybe_name = "NoYesMaybe"

var _YesNoMaybe_index = [...]uint8{0, 2, 5, 10}

func (i YesNoMaybe) String() string {
	if i < 0 || i >= YesNoMaybe(len(_YesNoMaybe_index)-1) {
		return fmt.Sprintf("YesNoMaybe(%d)", i)
	}
	return _YesNoMaybe_name[_YesNoMaybe_index[i]:_YesNoMaybe_index[i+1]]
}

// EmptyState provides an enum with choices, but without stringer.
type EmptyState int

const (
	NotEmpty EmptyState = iota
)

// TooLongState simulates an enum with an infinite number of values.
type TooLongState int

func (i TooLongState) String() string {
	return "repeats"
}
