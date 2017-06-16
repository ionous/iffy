package reflector

import (
	"fmt"
	"github.com/ionous/errutil"
	r "reflect"
	"strings"
)

type RefEnum struct {
	RefProp
	Enumeration
}

// Choice stores just the original name.
// May be expanded to include id or index for performance reasons.
type Choice string

func (c Choice) IsValid() bool {
	return len(c) > 0
}

func (c Choice) Id() string {
	return MakeId(string(c))
}

func (c Choice) Name() string {
	return string(c)
}

type Enumeration []Choice

func (enum Enumeration) ChoiceToIndex(choice string) (ret int) {
	id := MakeId(choice)
	return enum.choiceToIndex(id)
}

func (enum Enumeration) choiceToIndex(id string) (ret int) {
	found := false
	for i, c := range enum {
		if id == c.Id() {
			ret, found = i, true
			break
		}
	}
	if !found {
		ret = -1
	}
	return
}
func (enum Enumeration) IndexToChoice(idx int) (ret Choice) {
	if idx >= 0 && idx < len(enum) {
		ret = enum[idx]
	}
	return
}

func MakeEnum(rtype r.Type) (ret Enumeration, err error) {
	v := r.New(rtype).Elem()
	finished := false
	for i := int64(0); i < 64; i++ {
		v.SetInt(i) // note: you have to ask for the interface each time because the value is new each time
		if stringer, ok := v.Interface().(fmt.Stringer); !ok {
			err = errutil.New("enum has no strings. generate via stringer?", rtype)
			break
		} else {
			s := stringer.String()
			if strings.ContainsRune(s, '(') {
				finished = true
				break
			} else {
				choice := Choice(s)
				ret = append(ret, choice)
			}
		}
	}
	if !finished {
		err = errutil.New("enum end not found after", len(ret), "choices", rtype)
	} else if len(ret) == 0 {
		err = errutil.New("enum is empty", rtype)
	}
	return
}
