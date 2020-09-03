package tables

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	r "reflect"
)

type Prog struct {
	Type     string      // underscore_case name of type
	Fragment interface{} // nil pointer identifying the go implementation of the type name.
}

func (p *Prog) NewFragment() r.Value {
	rtype := r.TypeOf(p.Fragment).Elem()
	return r.New(rtype)
}

func (p *Prog) NewFragments() r.Value {
	rtype := r.SliceOf(r.TypeOf(p.Fragment))
	return r.New(rtype).Elem() // note: fragment containers are arrays, not pointers to arrays
}

func (p *Prog) Aggregate(rows *sql.Rows) (ret interface{}, err error) {
	var prog []byte
	rs := p.NewFragments()
	if e := ScanAll(rows, func() (err error) {
		rl := p.NewFragment()
		dec := gob.NewDecoder(bytes.NewBuffer(prog))
		if e := dec.DecodeValue(rl); e != nil {
			err = e
		} else {
			rs = r.Append(rs, rl)
		}
		return
	}, &prog); e != nil {
		err = e
	} else {
		ret = rs.Interface()
	}
	return
}
