package story

import (
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/imp"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/tables"
)

// make.str("pattern_name");
func imp_pattern_name(k *imp.Porter, r reader.Map) (ret ephemera.Named, err error) {
	if str, e := reader.String(r, "pattern_name"); e != nil {
		err = e
	} else {
		ret = k.NewName("pattern_name", str, reader.At(r))
	}
	return
}

// "an {activity:patterned_activity} or a {value:variable_type}");
func imp_pattern_type(k *imp.Porter, r reader.Map) (ret ephemera.Named, err error) {
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
func imp_patterned_activity(k *imp.Porter, r reader.Map) (ret ephemera.Named, err error) {
	if e := reader.Const(r, "patterned_activity", "$ACTIVITY"); e != nil {
		err = e
	} else {
		ret = k.NewName(tables.NAMED_TYPE, "execute", reader.At(r))
	}
	return
}

// make.run("pattern_handler", {name:pattern_name}{filters?pattern_filters}: {hook:pattern_hook}
func imp_pattern_handler(k *imp.Porter, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "pattern_handler"); e != nil {
		err = e
	} else if n, e := imp_pattern_name(k, m.MapOf("$NAME")); e != nil {
		err = e
	} else if rule, e := imp_pattern_hook(k, m.MapOf("$HOOK")); e != nil {
		err = e
	} else {
		if v := m.MapOf("$FILTERS"); len(v) != 0 {
			err = imp_pattern_filters(k, rule, v)
		}
		if err == nil {
			if prog, e := k.NewProg(rule.name, rule.distill()); e != nil {
				err = e
			} else {
				k.NewPatternHandler(n, prog)
			}
		}
	}
	return
}

//run("pattern_filters", " when {filter+bool_eval}"
func imp_pattern_filters(k *imp.Porter, rule *ruleBuilder, r reader.Map) (err error) {
	if m, e := reader.Unpack(r, "pattern_filters"); e != nil {
		err = e
	} else {
		err = reader.Repeats(m.SliceOf("$FILTER"),
			func(m reader.Map) (err error) {
				if i, e := k.DecodeSlot(m, "bool_eval"); e != nil {
					err = e
				} else {
					rule.addFilter(i.(rt.BoolEval))
				}
				return
			})
	}
	return
}

// opt("pattern_hook", "{activity:pattern_activity} or {result:pattern_return}");
func imp_pattern_hook(k *imp.Porter, r reader.Map) (ret *ruleBuilder, err error) {
	err = reader.Option(r, "pattern_hook", reader.ReadMaps{
		"$ACTIVITY": func(m reader.Map) (err error) {
			ret, err = imp_pattern_activity(k, m)
			return
		},
		"$RESULT": func(m reader.Map) (err error) {
			ret, err = imp_pattern_return(k, m)
			return
		},
	})
	return
}

// run("pattern_return", "return {result:pattern_result}");
// note: this slat exists for composer formatting reasons only...
func imp_pattern_return(k *imp.Porter, r reader.Map) (ret *ruleBuilder, err error) {
	if m, e := reader.Unpack(r, "pattern_return"); e != nil {
		err = e
	} else {
		ret, err = imp_pattern_result(k, m.MapOf("$RESULT"))
	}
	return
}

// opt("pattern_result", "a {simple value%primitive:primitive_func} or an {object:object_func}");
func imp_pattern_result(k *imp.Porter, r reader.Map) (ret *ruleBuilder, err error) {
	err = reader.Option(r, "pattern_result", reader.ReadMaps{
		"$PRIMITIVE": func(m reader.Map) (err error) {
			ret, err = imp_primitive_func(k, m)
			return
		},
		"$OBJECT": func(m reader.Map) (err error) {
			ret, err = imp_object_func(k, m)
			return
		},
	})
	return
}

// opt("primitive_func", "{a number%number_eval}, {some text%text_eval}, {a true/false value%bool_eval}")
func imp_primitive_func(k *imp.Porter, r reader.Map) (ret *ruleBuilder, err error) {
	err = reader.Option(r, "primitive_func", reader.ReadMaps{
		"$NUMBER_EVAL": func(m reader.Map) (err error) {
			if i, e := k.DecodeSlot(m, "number_eval"); e != nil {
				err = e
			} else {
				ret = newNumberRule(i.(rt.NumberEval))
			}
			return
		},
		"$TEXT_EVAL": func(m reader.Map) (err error) {
			if i, e := k.DecodeSlot(m, "text_eval"); e != nil {
				err = e
			} else {
				ret = newTextRule(i.(rt.TextEval))
			}
			return
		},
		"$BOOL_EVAL": func(m reader.Map) (err error) {
			if i, e := k.DecodeSlot(m, "bool_eval"); e != nil {
				err = e
			} else {
				ret = newBoolRule(i.(rt.BoolEval))
			}
			return
		},
	})
	return
}

// run("pattern_activity", "run: {activity%go+execute|ghost}");
func imp_pattern_activity(k *imp.Porter, r reader.Map) (ret *ruleBuilder, err error) {
	var exes []rt.Execute
	if m, e := reader.Unpack(r, "pattern_activity"); e != nil {
		err = e
	} else if e := reader.Repeats(m.SliceOf("$GO"),
		func(m reader.Map) (err error) {
			if i, e := k.DecodeSlot(m, "execute"); e != nil {
				err = e
			} else {
				exes = append(exes, i.(rt.Execute))
			}
			return
		}); e != nil {
		err = e
	} else {
		ret = newExecuteRule(exes)
	}
	return
}

// run("object_func", "an object named {name%text_eval}");
func imp_object_func(k *imp.Porter, r reader.Map) (ret *ruleBuilder, err error) {
	if m, e := reader.Unpack(r, "object_func"); e != nil {
		err = e
	} else if i, e := k.DecodeSlot(m, "text_eval"); e != nil {
		err = e
	} else {
		// FIX: we should wrap the text with a runtime "test the object matches" command
		// and, -- for simple text values -- add ephemera for "name, type"
		ret = newTextRule(i.(rt.TextEval))
	}
	return
}
