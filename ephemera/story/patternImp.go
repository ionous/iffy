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
	} else if e := imp_pattern_rules(k, patternName, op.MapOf("$PATTERN_RULES")); e != nil {
		err = e
	} else if l := op.MapOf("$PATTERN_LOCALS"); len(l) > 0 {
		err = imp_pattern_locals(k, patternName, l)
	}
	return
}

func imp_pattern_locals(k *Importer, patternName ephemera.Named, r reader.Map) (err error) {
	if op, e := reader.Unpack(r, "pattern_locals"); e != nil {
		err = e
	} else {
		err = reader.Repeats(op.SliceOf("$LOCAL_DECL"),
			func(el reader.Map) (err error) {
				return imp_pattern_local(k, patternName, el)
			})
	}
	return
}

// read a single local_decl from a pattern action block
func imp_pattern_local(k *Importer, patternName ephemera.Named, r reader.Map) (err error) {
	// xxxxxxxxxxxxxxxxxxxxxxx
	// if op, e := reader.Unpack(r, "local_decl"); e != nil {
	// 	err = e
	// } else if localName, localType, e := imp_variable_decl(k, tables.NAMED_LOCAL, op.MapOf("$VARIABLE_DECL")); e != nil {
	// 	err = e
	// } else if hook, e := imp_program_hook(k, op.MapOf("$PROGRAM_RESULT")); e != nil {
	// 	err = e
	// } else if prog, e := k.NewGob(hook.SlotType(), hook.CmdPtr()); e != nil {
	// 	err = e // turn the result generator into a storable program.
	// } else {
	// 	// 	k.NewPatternRule(patternName, patternProg)
	// }
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
