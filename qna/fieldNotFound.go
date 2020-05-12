package qna

import "github.com/ionous/errutil"

type fieldNotFound struct {
	owner, field string
}

func (f fieldNotFound) Error() string {
	return errutil.New("field not found", f.owner, f.field).Error()
}
