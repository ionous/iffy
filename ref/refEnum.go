package ref

import (
	"fmt"
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/id"
	r "reflect"
	"strings"
)

type RefEnum struct {
	RefProp
	Enumeration
}

// Choice stores the name of a single enumerated value.
// ( NOTE: in the future, performance reasons, this may be expanded to include the id or index of the value. )
type Choice string

func (c Choice) IsValid() bool {
	return len(c) > 0
}

func (c Choice) Id() string {
	return id.MakeId(string(c))
}

func (c Choice) Name() string {
	return string(c)
}

// Enumeration collects a number of choices.
type Enumeration []Choice

func (enum Enumeration) setValue(field r.Value, idx int, val bool) (err error) {
	// if the field is a bool, and we found it via getPropertyByChoice,
	// then name must equal the name of the field, and we directly directly setting its status
	if field.Kind() == r.Bool {
		err = CoerceValue(field, val)
	} else {
		// if the field in an int, and the user is trying to set a particular choice
		// we want to set the field to the value of that passed choice.
		if val {
			err = CoerceValue(field, idx)
		} else {
			// if the user is saying unset some choice
			// we have to try to generate an opposite value.
			if invIdx, e := enum.InverseIndex(idx); e != nil {
				err = e
			} else {
				err = CoerceValue(field, invIdx)
			}
		}
	}
	return
}

func (enum Enumeration) InverseIndex(idx int) (ret int, err error) {
	if cnt := len(enum); cnt > 2 {
		err = errutil.New("no opposite value. too many choices", cnt)
	} else {
		// idx= 0; 2-(0+1)=1
		// idx= 1; 2-(1+1)=0
		// ret can be out of range for 1 length enums
		ret = 2 - (idx + 1)
	}
	return
}

// -1 if not found
func (enum Enumeration) ChoiceToIndex(choice string) (ret int) {
	id := id.MakeId(choice)
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

func EnumFromField(field *r.StructField) (ret Enumeration, err error) {
	switch rtype := field.Type; rtype.Kind() {
	default:
		err = errutil.New("unexpected enum", rtype)
	case r.Bool:
		ret = Enumeration{Choice(field.Name)}
	case r.Int:
		ret, err = makeEnum(rtype)
	}
	return
}

func MakeEnum(enum interface{}) (ret Enumeration, err error) {
	if etype := r.TypeOf(enum); etype.Kind() != r.Ptr {
		err = errutil.New("expected pointer (to int)")
	} else if rtype := etype.Elem(); rtype.Kind() != r.Int {
		err = errutil.New("expected an int pointer")
	} else {
		ret, err = makeEnum(rtype)
	}
	return
}

func makeEnum(rtype r.Type) (ret Enumeration, err error) {
	// contruct an enum value of the passed type
	// to generate a list of enumerated choices.
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
