package express

import (
	r "reflect"

	"github.com/ionous/errutil"
)

type rstack struct {
	els []r.Value
}

//
func (k *rstack) flush() (ret interface{}, err error) {
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
func (k *rstack) push(cmd ...r.Value) {
	k.els = append(k.els, cmd...)
}

func (k *rstack) pop(cnt int) (ret []r.Value, err error) {
	if end := len(k.els) - cnt; end < 0 {
		err = errutil.New("stack underflow")
	} else {
		ret, k.els = k.els[end:], k.els[:end]
	}
	return
}

func (k *rstack) peek(cnt int) (ret []r.Value, err error) {
	if end := len(k.els) - cnt; end < 0 {
		err = errutil.New("stack underflow")
	} else {
		ret = k.els[end:]
	}
	return
}
