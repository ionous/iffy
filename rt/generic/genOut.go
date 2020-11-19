package generic

import "github.com/ionous/iffy/affine"

func RecordsToValue(ds []*Record) []interface{} {
	var els []interface{}
	for _, d := range ds {
		els = append(els, RecordToValue(d))
	}
	return els
}

// future: json encoding instead
func RecordToValue(d *Record) map[string]interface{} {
	m := make(map[string]interface{})
	for i, f := range d.kind.fields {
		if rv, e := d.GetFieldByIndex(i); e != nil {
			panic(e)
		} else {
			var el interface{}
			switch a := rv.Affinity(); a {
			case affine.Record:
				if d, e := rv.GetRecord(); e != nil {
					panic(e)
				} else {
					el = RecordToValue(d)
				}
			case affine.RecordList:
				if ds, e := rv.GetRecordList(); e != nil {
					panic(e)
				} else {
					el = RecordsToValue(ds)
				}
			default:
				el = rv.(refValue).v.Interface()
			}
			m[f.Name] = el
		}
	}
	return m
}
