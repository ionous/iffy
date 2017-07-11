package stream

import (
	"github.com/ionous/iffy/rt"
)

// ObjectBinder splats joins multiple object streams into once.
// Its unused and untested: it might be used for patterns to prefix/postfix streams
type ObjectBinder struct {
	evals []rt.ObjListEval
}

type ObjectBindIt struct {
	run   rt.Runtime
	evals []rt.ObjListEval
	it    rt.ObjectStream
	next  rt.Object
	err   error
}

func (b *ObjectBinder) Append(eval rt.ObjListEval) {
	if chain, ok := eval.(*ObjectBinder); ok {
		b.evals = append(b.evals, chain.evals...)
	} else {
		b.evals = append(b.evals, eval)
	}
}

func (b *ObjectBinder) GetObjectStream(run rt.Runtime) (rt.ObjectStream, error) {
	it := &ObjectBindIt{run: run, evals: b.evals, it: &ObjectIt{}}
	it.next, it.err = it.advance()
	return it, nil
}

func (b *ObjectBindIt) advance() (next rt.Object, err error) {
	if b.it.HasNext() {
		next, err = b.it.GetNext()
	} else {
		for len(b.evals) > 0 {
			var eval rt.ObjListEval
			eval, b.evals = b.evals[0], b.evals[1:]
			if s, e := eval.GetObjectStream(b.run); e != nil {
				err = e
				break
			} else if s.HasNext() {
				next, err = s.GetNext()
				break
			}
		}
	}
	return
}

func (b *ObjectBindIt) HasNext() bool {
	return b.err != nil && b.next != nil
}

func (b *ObjectBindIt) GetNext() (ret rt.Object, err error) {
	if b.err != nil {
		err = b.err
	} else if b.next == nil {
		err = rt.StreamExceeded
	} else {
		b.next, b.err = b.advance()
	}
	return
}
