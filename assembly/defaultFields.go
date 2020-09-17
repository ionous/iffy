package assembly

import (
	"fmt"
	"reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/tables"
)

// reads mdl_field, mdl_kind
func assembleDefaultFields(asm *Assembler) (err error) {
	var store valueStore
	var curr, last valueInfo
	if e := tables.QueryAll(asm.cache.DB(),
		`select asm.kind, mf.field, mf.type, asm.value
 			from asm_default as asm
 		join mdl_field mf
 			on (asm.prop = mf.field)
 			and (mf.type != 'aspect')
 		/* is the field's declared kind in the path of the user specified kind */		
 		where instr((
 			select mk.kind || "," || mk.path || ","
			from mdl_kind mk 
			where mk.kind = asm.kind
		),  mf.kind || ",")
		order by asm.kind, mf.field`,
		func() (err error) {
			if nv, e := convertField(curr.fieldType, curr.value); e != nil {
				err = e
			} else if last.target != curr.target || last.field != curr.field {
				store.add(last)
				last, last.value = curr, nv
			} else if !reflect.DeepEqual(last.value, nv) {
				err = errutil.Fmt("conflicting defaults: %s != %s:%T", last.String(), curr.String(), nv)
			}
			return
		},
		&curr.target, &curr.field, &curr.fieldType,
		&curr.value,
	); e != nil {
		err = e
	} else {
		store.add(last)
		err = store.writeDefaultFields(asm)
	}
	return
}

type valueInfo struct {
	target, field, fieldType string
	value                    interface{}
}

func (n *valueInfo) String() string {
	return n.target + "." + n.field + ":" + n.fieldType + fmt.Sprintf("(%v:%T)", n.value, n.value)
}

type valueStore struct {
	list []valueInfo
}

func (store *valueStore) add(n valueInfo) {
	if len(n.target) > 0 {
		store.list = append(store.list, n)
	}
}

func (store *valueStore) writeDefaultFields(asm *Assembler) (err error) {
	for _, n := range store.list {
		if e := asm.WriteDefault(n.target, n.field, n.value); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

func (store *valueStore) writeInitialFields(asm *Assembler) (err error) {
	for _, n := range store.list {
		if e := asm.WriteStart(n.target, n.field, n.value); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// out types are currently: int, float32, or string.
func convertField(fieldType string, value interface{}) (ret interface{}, err error) {
	switch v := reflect.ValueOf(value); fieldType {
	case tables.PRIM_DIGI:
		switch k := v.Kind(); k {
		case reflect.Float64:
			ret = float32(v.Float())
		case reflect.Int64:
			ret = int(v.Int())
		default:
			err = errutil.Fmt("can't convert [%v](%s) to %s", value, k, fieldType)
		}
	case tables.PRIM_TEXT:
		switch k := v.Kind(); k {
		case reflect.String:
			ret = v.String()
		default:
			err = errutil.Fmt("can't convert [%v](%s) to %s", value, k, fieldType)
		}
	default:
		err = errutil.New("convertField: unhandled field type", fieldType)
	}
	return
}
