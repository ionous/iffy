package story

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/tables"
)

func imp_determine_act(k *Importer, m reader.Map) (ret interface{}, err error) {
	var from pattern.DetermineAct
	if e := fromPattern(k, m, &from, (*pattern.FromPattern)(&from)); e != nil {
		err = e
	} else {
		ret = &from
	}
	return
}
func imp_determine_num(k *Importer, m reader.Map) (ret interface{}, err error) {
	var from pattern.DetermineNum
	if e := fromPattern(k, m, &from, (*pattern.FromPattern)(&from)); e != nil {
		err = e
	} else {
		ret = &from
	}
	return
}
func imp_determine_text(k *Importer, m reader.Map) (ret interface{}, err error) {
	var from pattern.DetermineText
	if e := fromPattern(k, m, &from, (*pattern.FromPattern)(&from)); e != nil {
		err = e
	} else {
		ret = &from
	}
	return
}
func imp_determine_bool(k *Importer, m reader.Map) (ret interface{}, err error) {
	var from pattern.DetermineBool
	if e := fromPattern(k, m, &from, (*pattern.FromPattern)(&from)); e != nil {
		err = e
	} else {
		ret = &from
	}
	return
}
func imp_determine_num_list(k *Importer, m reader.Map) (ret interface{}, err error) {
	var from pattern.DetermineNumList
	if e := fromPattern(k, m, &from, (*pattern.FromPattern)(&from)); e != nil {
		err = e
	} else {
		ret = &from
	}
	return
}
func imp_determine_text_list(k *Importer, m reader.Map) (ret interface{}, err error) {
	var from pattern.DetermineTextList
	if e := fromPattern(k, m, &from, (*pattern.FromPattern)(&from)); e != nil {
		err = e
	} else {
		ret = &from
	}
	return
}

// used by "determine_num", etc.
// from and spec are the same object, but go-type unpacking from interfaces isnt always easy.
func fromPattern(k *Importer,
	m reader.Map,
	spec composer.Slat,
	from *pattern.FromPattern) (err error) {
	typeName := spec.Compose().Name
	if evalType, e := slotName(spec); e != nil {
		err = e
	} else if m, e := reader.Unpack(m, typeName); e != nil {
		err = e
	} else if patternName, e := imp_pattern_name(k, m.MapOf("$PATTERN")); e != nil {
		err = e
	} else if ps, e := imp_parameters(k, patternName, m.MapOf("$PARAMETERS")); e != nil {
		err = e
	} else {
		// fix: object type names will need adaption of some sort re plural_kinds
		patternType := k.NewName(evalType, tables.NAMED_TYPE, reader.At(m))
		k.NewPatternRef(patternName, patternName, patternType)
		// assign results
		from.Pattern = patternName.String()
		from.Parameters = ps
	}
	return
}

func imp_parameters(k *Importer, pid ephemera.Named, r reader.Map) (ret *pattern.Parameters, err error) {
	if m, e := reader.Unpack(r, "parameters"); e != nil {
		err = e
	} else {
		var ps []*pattern.Parameter
		// record the referenced parameters: names, types, pairings
		if e := reader.Repeats(m.SliceOf("$PARAMS"),
			func(m reader.Map) (err error) {
				if p, e := imp_parameter(k, pid, m); e != nil {
					err = e
				} else {
					ps = append(ps, p)
				}
				return
			}); e != nil {
			err = e
		} else {
			ret = &pattern.Parameters{ps}
		}
	}
	return
}

func imp_parameter(k *Importer, patternName ephemera.Named, r reader.Map) (ret *pattern.Parameter, err error) {
	if m, e := reader.Unpack(r, "parameter"); e != nil {
		err = e
	} else if paramName, e := imp_variable_name(k, m.MapOf("$NAME")); e != nil {
		err = e
	} else if a, e := imp_assignment(k, m.MapOf("$FROM")); e != nil {
		err = e
	} else if slotName, e := slotName(a.GetEval()); e != nil {
		err = e // ^ its possible this should be the type that the assignment contains.
	} else {
		paramType := k.NewName(slotName, tables.NAMED_TYPE, reader.At(m))
		k.NewPatternRef(patternName, paramName, paramType)
		ret = &pattern.Parameter{Name: paramName.String(), From: a}
	}
	return
}

func imp_assignment(k *Importer, m reader.Map) (ret core.Assignment, err error) {
	err = k.DecodeSlot(m, "assignment", &ret)
	return
}

func imp_text(k *Importer, m reader.Map) (ret string, err error) {
	return reader.String(m, "text")
}
