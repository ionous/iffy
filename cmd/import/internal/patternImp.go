package internal

import (
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/tables"
)

// make.str("pattern_name");
func imp_pattern_name(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	return k.namedType(r, "pattern_name")
}

// "an {activity:patterned_activity} or a {value:variable_type}");
func imp_pattern_type(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	err = k.expectOpt(r, "pattern_type", map[string]Parse{
		"$ACTIVITY": func(k *Importer, m reader.Map) (err error) {
			ret, err = imp_patterned_activity(k, m)
			return
		},
		"$VALUE": func(k *Importer, m reader.Map) (err error) {
			ret, err = imp_variable_type(k, m)
			return
		},
	})
	return
}

// make.str("patterned_activity", "{activity}");
// returns "prog" as the name of a type  ( eases the difference b/t user named kinds, and internally named types )
func imp_patterned_activity(k *Importer, r reader.Map) (ret ephemera.Named, err error) {
	if e := k.expectConst(r, "patterned_activity", "$ACTIVITY"); e != nil {
		err = e
	} else {
		ret = k.Named(tables.NAMED_TYPE, "execute", at(r))
	}
	return
}

// make.run("pattern_handler", {name:pattern_name}{filters?pattern_filters}: {hook:pattern_hook}
func imp_pattern_handler(k *Importer, r reader.Map) (err error) {
	if m, e := k.expectSlat(r, "pattern_handler"); e != nil {
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
			if prog, e := k.newProg(rule.name, rule.distill()); e != nil {
				err = e
			} else {
				k.NewPatternHandler(n, prog)
			}
		}
	}
	return
}

//run("pattern_filters", " when {filter+bool_eval}"
func imp_pattern_filters(k *Importer, rule *ruleBuilder, r reader.Map) (err error) {
	if m, e := k.expectSlat(r, "pattern_filters"); e != nil {
		err = e
	} else {
		err = k.repeats(m.SliceOf("$FILTER"),
			func(k *Importer, m reader.Map) (err error) {
				if i, e := k.expectProg(m, "bool_eval"); e != nil {
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
func imp_pattern_hook(k *Importer, r reader.Map) (ret *ruleBuilder, err error) {
	err = k.expectOpt(r, "pattern_hook", map[string]Parse{
		"$ACTIVITY": func(k *Importer, m reader.Map) (err error) {
			// FIX: what's the "name" for again...?
			ret, err = imp_pattern_activity(k, m)
			return
		},
		"$RESULT": func(k *Importer, m reader.Map) (err error) {
			ret, err = imp_object_func(k, m)
			return
		},
	})
	return
}

// run("pattern_return", "return {result:pattern_result}");
// note: this slat exists for composer formatting reasons only...
func imp_pattern_return(k *Importer, r reader.Map) (ret *ruleBuilder, err error) {
	if m, e := k.expectSlat(r, "pattern_return"); e != nil {
		err = e
	} else {
		ret, err = imp_pattern_result(k, m)
	}
	return
}

// opt("pattern_result", "a {simple value%primitive:primitive_func} or an {object:object_func}");
func imp_pattern_result(k *Importer, r reader.Map) (ret *ruleBuilder, err error) {
	err = k.expectOpt(r, "pattern_result", map[string]Parse{
		"$PRIMITIVE": func(k *Importer, m reader.Map) (err error) {
			ret, err = imp_primitive_func(k, m)
			return
		},
		"$OBJECT": func(k *Importer, m reader.Map) (err error) {
			ret, err = imp_object_func(k, m)
			return
		},
	})
	return
}

// opt("primitive_func", "{a number%number_eval}, {some text%text_eval}, {a true/false value%bool_eval}")
func imp_primitive_func(k *Importer, r reader.Map) (ret *ruleBuilder, err error) {
	err = k.expectOpt(r, "primitive_func", map[string]Parse{
		"$NUMBER_EVAL": func(k *Importer, m reader.Map) (err error) {
			if i, e := k.expectProg(m, "number_eval"); e != nil {
				err = e
			} else {
				ret = newNumberRule(i.(rt.NumberEval))
			}
			return
		},
		"$TEXT_EVAL": func(k *Importer, m reader.Map) (err error) {
			if i, e := k.expectProg(m, "text_eval"); e != nil {
				err = e
			} else {
				ret = newTextRule(i.(rt.TextEval))
			}
			return
		},
		"$BOOL_EVAL": func(k *Importer, m reader.Map) (err error) {
			if i, e := k.expectProg(m, "bool_eval"); e != nil {
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
func imp_pattern_activity(k *Importer, r reader.Map) (ret *ruleBuilder, err error) {
	var exes []rt.Execute
	if m, e := k.expectSlat(r, "pattern_activity"); e != nil {
		err = e
	} else if e := k.repeats(m.SliceOf("$GO"),
		func(k *Importer, m reader.Map) (err error) {
			if i, e := k.expectProg(m, "execute"); e != nil {
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
func imp_object_func(k *Importer, r reader.Map) (ret *ruleBuilder, err error) {
	if m, e := k.expectSlat(r, "object_func"); e != nil {
		err = e
	} else if i, e := k.expectProg(m, "text_eval"); e != nil {
		err = e
	} else {
		// FIX: we should wrap the text with a runtime "test the object matches" command
		// and, -- for simple text values -- add ephemera for "name, type"
		ret = newTextRule(i.(rt.TextEval))
	}
	return
}
