package express

import (
	r "reflect"

	"github.com/ionous/errutil"
)

// contains the results of converting postfix template functions to iffy commands.
// the fields of the commands are filled by the time they get into this stack.
type cmdStack struct {
	els []r.Value
}

//
func (k *cmdStack) flush() (ret interface{}, err error) {
	if cnt := len(k.els); cnt == 0 {
		err = errutil.New("empty output")
	} else if cnt > 1 {
		err = errutil.New("unparsed output")
	} else if cmd := k.els[0]; !cmd.IsValid() {
		err = errutil.New("convert returned invalid value")
	} else {
		ret, k.els = cmd.Interface(), nil
	}
	return
}

// add a new command pointer to the output stack.
func (k *cmdStack) push(cmd ...r.Value) {
	k.els = append(k.els, cmd...)
}

// removes a section of the stack
func (k *cmdStack) pop(cnt int) (ret []r.Value, err error) {
	if end := len(k.els) - cnt; end < 0 {
		err = underflow(end)
	} else {
		ret, k.els = k.els[end:], k.els[:end]
	}
	return
}

// returns a section of the stack without changing the stack.
func (k *cmdStack) peek(cnt int) (ret []r.Value, err error) {
	if end := len(k.els) - cnt; end < 0 {
		err = underflow(end)
	} else {
		ret = k.els[end:]
	}
	return
}

type underflow int

func (i underflow) Error() string {
	return errutil.Sprintf("stack underflow %d", int(i))
}
