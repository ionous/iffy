package generic

import "github.com/ionous/iffy/affine"

// primarily for testing:
// convert the contents of a bunch of records into a printable format.
func RecordsToValue(ds []*Record) []interface{} {
	var els []interface{}
	for _, d := range ds {
		els = append(els, RecordToValue(d))
	}
	return els
}

// primarily for testing:
// convert the contents of a record into a printable format.
// future: json encoding instead
func RecordToValue(d *Record) map[string]interface{} {
	m := make(map[string]interface{})
	for i, f := range d.kind.fields {
		if rv, e := d.GetFieldByIndex(i); e != nil {
			panic(e)
		} else {
			var el interface{}
			switch a := rv.Affinity(); a {
			case affine.TextList:
				el = rv.Strings()
			case affine.NumList:
				el = rv.Floats()
			case affine.Record:
				el = RecordToValue(rv.Record())
			case affine.RecordList:
				el = RecordsToValue(rv.Records())
			default:
				el = rv.(refValue).i
			}
			m[f.Name] = el
		}
	}
	return m
}
