package story

import (
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/tables"
)

// "an {activity:patterned_activity} or a {value:variable_type}");
func imp_pattern_type(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	err = reader.Option(r, "pattern_type", reader.ReadMaps{
		"$ACTIVITY": func(m reader.Map) (err error) {
			ret, err = imp_patterned_activity(k, m)
			return
		},
		"$VALUE": func(m reader.Map) (err error) {
			ret, err = imp_variable_type(k, m)
			return
		},
	})
	return
}

// make.str("patterned_activity", "{activity}");
// returns "prog" as the name of a type  ( eases the difference b/t user named kinds, and internally named types )
func imp_patterned_activity(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	if e := reader.Const(r, "patterned_activity", "$ACTIVITY"); e != nil {
		err = e
	} else {
		ret = k.NewName("execute", tables.NAMED_TYPE, reader.At(r))
	}
	return
}

// make.run("pattern_actions", ... "To {name:pattern_name}: {+pattern_rule}")
func imp_pattern_actions(k *Importer, r reader.Map) (err error) {
	if op, e := reader.Unpack(r, "pattern_actions"); e != nil {
		err = e
	} else if patternName, e := imp_pattern_name(k, op.MapOf("$NAME")); e != nil {
		err = e
	} else {
		err = imp_pattern_rules(k, patternName, op.MapOf("$PATTERN_RULES"))
	}
	return
}

func imp_pattern_rules(k *Importer, patternName ephemera.Named, r reader.Map) (err error) {
	if op, e := reader.Unpack(r, "pattern_rules"); e != nil {
		err = e
	} else {
		err = reader.Repeats(op.SliceOf("$PATTERN_RULE"),
			func(el reader.Map) (err error) {
				return imp_pattern_rule(k, patternName, el)
			})
	}
	return
}

func imp_pattern_rule(k *Importer, patternName ephemera.Named, r reader.Map) (err error) {
	var guard rt.BoolEval
	if op, e := reader.Unpack(r, "pattern_rule"); e != nil {
		err = e
	} else if e := k.DecodeSlot(op.MapOf("$GUARD"), "bool_eval", &guard); e != nil {
		err = e
	} else if slot, e := imp_program_hook(k, op.MapOf("$HOOK")); e != nil {
		err = e
	} else {
		name, rule := slot.NewRule(guard)
		//
		if patternProg, e := k.NewGob(name, rule); e != nil {
			err = e
		} else {
			k.NewPatternRule(patternName, patternProg)
		}
	}
	return
}
