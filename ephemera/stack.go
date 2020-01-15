package ephemera

import "github.com/ionous/errutil"

// Stack multiple queues to act as one
type Stack struct {
	qs []Queue
}

// NewStack of queue(s)
func NewStack(qs ...Queue) Queue {
	return &Stack{qs: qs}
}

// Prep implements Queue
func (j *Stack) Prep(which string, cols ...Col) {
	for _, q := range j.qs {
		q.Prep(which, cols...)
	}
}

// Write implements Queue
func (j *Stack) Write(which string, args ...interface{}) (ret Queued, err error) {
	for _, q := range j.qs {
		if r, e := q.Write(which, args...); e != nil {
			err = errutil.Append(err, e)
		} else {
			ret = r
		}
	}
	return
}
