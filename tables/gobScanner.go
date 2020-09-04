package tables

import (
	"bytes"
	"encoding/gob"
	r "reflect"

	"github.com/ionous/errutil"
)

type GobScanner struct {
	Target r.Value
}

func NewGobScanner(ptr interface{}) *GobScanner {
	return &GobScanner{r.ValueOf(ptr).Elem()}
}

func (gs *GobScanner) Scan(val interface{}) (err error) {
	if b, ok := val.([]byte); !ok {
		err = errutil.Fmt("gob scanner received unexpected type %T", val)
	} else {
		dec := gob.NewDecoder(bytes.NewBuffer(b))
		err = dec.DecodeValue(gs.Target)
	}
	return
}
