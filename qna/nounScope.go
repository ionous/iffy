package qna

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/rt/scope"
	"github.com/ionous/iffy/tables"
)

// NounScope validates requests for noun names,
// returning the name or an error if the noun doesnt exist.
type NounScope struct {
	scope.EmptyScope
	db *tables.Cache
}

func (ns *NounScope) GetVariable(name string) (ret interface{}, err error) {
	var noun string
	if e := ns.db.QueryRow(
		`select noun 
		from mdl_name
		where name=?
		order by rank
		limit 1 `, name).Scan(&noun); e != nil {
		err = e
	} else if len(noun) == 0 {
		err = errutil.Fmt("unknown noun %q", name)
	} else {
		ret = noun
	}
	return
}
