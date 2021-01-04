package internal

import (
	"github.com/PaesslerAG/jsonpath"
	"github.com/ionous/errutil"
)

//
func replicate(doc interface{}, fromParent, fromField, toParent, toField string) (err error) {
	if src, e := jsonpath.Get(fromParent, doc); e != nil {
		err = e
	} else if fromEls, ok := src.([]interface{}); !ok {
		err = errutil.Fmt("unknown src %T", src)
	} else if dst, e := jsonpath.Get(toParent, doc); e != nil {
		err = e
	} else if toEls, ok := dst.([]interface{}); !ok {
		err = errutil.Fmt("unknown dst %T", dst)
	} else if fromCnt, toCnt := len(fromEls), len(toEls); fromCnt != toCnt {
		err = errutil.Fmt("mismatched copy, from %d to %d", fromCnt, toCnt)
	} else {
		for i := 0; i < fromCnt; i++ {
			toEl, fromEl := toEls[i], fromEls[i]
			if from, ok := fromEl.(map[string]interface{}); !ok {
				err = errutil.Fmt("expected a slice of objects; got %T", fromEl)
				break
			} else if to, ok := toEl.(map[string]interface{}); !ok {
				err = errutil.Fmt("expected a slice of objects; got %T", toEl)
				break
			} else {
				to[toField] = from[fromField]
			}
		}
	}
	return
}

func replace(doc interface{}, parent, field string, value interface{}) (err error) {
	if tgt, e := jsonpath.Get(parent, doc); e != nil {
		err = e
	} else if els, ok := tgt.([]interface{}); !ok {
		err = errutil.Fmt("unknown target %T", tgt)
	} else {
		for _, el := range els {
			if obj, ok := el.(map[string]interface{}); !ok {
				err = errutil.Fmt("expected a slice of objects; got %T", el)
				break
			} else if value == nil {
				delete(obj, field)
			} else {
				obj[field] = value
			}
		}
	}
	return
}

//
// func rename(doc interface{}, parent, field, newField string) (err error) {
// 	if tgt, e := jsonpath.Get(parent, doc); e != nil {
// 		err = e
// 	} else if els, ok := tgt.([]interface{}); !ok {
// 		err = errutil.Fmt("unknown target %T", tgt)
// 	} else {
// 		for _, el := range els {
// 			if obj, ok := el.(map[string]interface{}); !ok {
// 				err = errutil.Fmt("expected a slice of objects; got %T", el)
// 				break
// 			} else {
// 				if len(newField) > 0 {
// 					obj[newField] = obj[field]
// 				}
// 				delete(obj, field)
// 			}
// 		}
// 	}
// 	return
// }

//
