package assembly

import "github.com/ionous/errutil"

type TargetField struct {
	Target, Field string
}

type TargetValue struct {
	Target, Field string
	Value         interface{}
}

// create some fake model hierarchy using Target and Field
func AddTestHierarchy(m *Modeler, keyValues []TargetField) (err error) {
	for _, tv := range keyValues {
		if e := m.WriteAncestor(tv.Target, tv.Field); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// mdl_noun:  kind, ret id noun
// mdl_name: noun, name, rank
func AddTestNouns(m *Modeler, els []TargetField) (err error) {
	for _, tv := range els {
		noun, kind := tv.Target, tv.Field
		if e := m.WriteNounWithNames(noun, kind); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// create some fake model hierarchy using Value as the field type
func AddTestFields(m *Modeler, targetValues []TargetValue) (err error) {
	for _, tv := range targetValues {
		if e := m.WriteField(tv.Target, tv.Field, tv.Value.(string)); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// write aspect, trait pairs using Target and Field
func AddTestTraits(m *Modeler, keyValues []TargetField) (err error) {
	for _, tv := range keyValues {
		// rank is not set yet, see DetermineAspects
		if e := m.WriteTrait(tv.Target, tv.Field, 0); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// write some noun, field, value triplets
func AddTestStarts(m *Modeler, targetValues []TargetValue) (err error) {
	for _, tv := range targetValues {
		if e := m.WriteValue(tv.Target, tv.Field, tv.Value); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// AddTestDefaults writes some kind, field, value triplets
func AddTestDefaults(m *Modeler, targetValues []TargetValue) (err error) {
	for _, tv := range targetValues {
		if e := m.WriteDefault(tv.Target, tv.Field, tv.Value); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

// add modeled data: relation, kind, cardinality, otherKind
func AddTestRelations(m *Modeler, relations [][4]string) (err error) {
	for _, el := range relations {
		relation, kind, cardinality, otherKind := el[0], el[1], el[2], el[3]
		if e := m.WriteRelation(relation, kind, cardinality, otherKind); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}

func AddTestVerbs(m *Modeler, relationStem [][2]string) (err error) {
	for _, el := range relationStem {
		rel, verb := el[0], el[1]
		if e := m.WriteVerb(rel, verb); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}
