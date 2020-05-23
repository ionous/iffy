package qna

import "github.com/ionous/errutil"

type fieldNotFound struct {
	owner, field string
}

func (f fieldNotFound) Error() string {
	return errutil.Sprint("field not found", f.owner, f.field)
}
