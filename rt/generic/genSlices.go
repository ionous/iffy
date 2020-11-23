package generic

import (
	"github.com/ionous/iffy/affine"
)

// duplicate the passed slice of floats
// ( b/c golang's built in copy doesnt allocate )
func copyFloats(src []float64) []float64 {
	out := make([]float64, len(src))
	copy(out, src)
	return out
}

//  duplicate the passed slice of strings.
func copyStrings(src []string) []string {
	out := make([]string, len(src))
	copy(out, src)
	return out
}

// duplicate the passed slice of record pointers
// note: doesnt copy the contents of the records, just the pointers.
func copyRecords(src []*Record) []*Record {
	out := make([]*Record, len(src))
	copy(out, src)
	return out
}

// rare, just for splice for now.
func safeAffinity(v Value) (ret affine.Affinity) {
	if v != nil {
		ret = v.Affinity()
	}
	return
}

// change a number or number_list into a slice of floats
// panics if the passed value isnt one of those two types.
func normalizeFloats(v Value) (ret []float64) {
	switch a := safeAffinity(v); a {
	case "":
	case affine.Number:
		one := v.Float()
		ret = []float64{one}
	case affine.NumList:
		ret = v.Floats()
	default:
		panic("cant create floats from " + a.String())
	}
	return
}

// change a string or string_list into a slice of strings
// panics if the passed value isnt one of those two types.
func normalizeStrings(v Value) (ret []string) {
	switch a := safeAffinity(v); a {
	case "":
	case affine.Text:
		one := v.String()
		ret = []string{one}
	case affine.TextList:
		ret = v.Strings()
	default:
		panic("cant create strings from " + a.String())
	}
	return
}

// change a record or record_list into a slice of record pointers
// panics if the passed value isnt one of those two types.
func normalizeRecords(v Value) (ret []*Record) {
	switch a := safeAffinity(v); a {
	case "": // nil
	case affine.Record:
		one := v.Record()
		ret = []*Record{one}
	case affine.RecordList:
		ret = v.Records()
	default:
		panic("cant create records from " + a.String())
	}
	return
}
