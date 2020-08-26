package qna

import "github.com/ionous/errutil"

type fieldNotFound struct {
	owner, field string
}

func (f fieldNotFound) Error() string {
	return errutil.Sprintf("field not found '%s.%s'", f.owner, f.field)
}
