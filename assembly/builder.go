package assembly

import (
	"bytes"
	"database/sql"
	"encoding/gob"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/lang"
	"github.com/ionous/iffy/tables"
)

type BuildRule struct {
	Query        string
	NewContainer func(name string) interface{}
	NewEl        func(c interface{}) interface{}
}

// map name to pattern interface
type patternEntry struct {
	patternName string              // name of the pattern
	patternType string              // "return" type of the pattern
	prologue    []pattern.Parameter // list of all parameters sent to the pattern
	locals      []pattern.Parameter // ...
}

func (pat *patternEntry) AddParam(param pattern.Parameter) {
	pat.prologue = append(pat.prologue, param)
}

type patternCache map[string]*patternEntry

func (cache patternCache) init(name, patternType string) (ret pattern.CommonPattern, okay bool) {
	if c, ok := (cache)[name]; ok && c.patternType == patternType {
		ret.Name = name
		ret.Prologue = c.prologue
		ret.Locals = c.locals
		okay = true
	}
	return
}

// read pattern declarations from the ephemera db
func buildPatternCache(db *sql.DB) (ret patternCache, err error) {
	// build the pattern cache
	out := make(patternCache)
	var patternName, paramName, typeName string
	var kind sql.NullString
	var prog []byte
	var last *patternEntry
	if e := tables.QueryAll(db,
		`select ap.pattern, ap.param, ap.type, ap.kind, ep.prog 
		from asm_pattern_decl ap
		left join eph_prog ep
		on (ep.rowid = ap.idProg)`,
		func() (err error) {
			// fix: need to handle conflicting prog definitions
			if last == nil || last.patternName != patternName {
				if patternName != paramName {
					err = errutil.New("expected the first param should be the pattern return type", patternName, paramName, typeName)
				} else {
					last = &patternEntry{patternName: patternName, patternType: typeName}
					out[patternName] = last
				}
			}
			if err == nil && paramName != patternName {
				// fix: these should probably be tables.PRIM_ names
				// ie. "text" not "text_eval" -- tests and other things have to be adjusted
				paramName := lang.Camelize(paramName)
				//
				switch typeName {
				case "text_eval":
					p := pattern.TextParam{Name: paramName}
					if e := decode(prog, &p.Init); e != nil {
						err = errutil.New("couldnt decode", patternName, paramName, e)
					} else {
						last.AddParam(&p)
					}
				case "number_eval":
					p := pattern.NumParam{Name: paramName}
					if e := decode(prog, &p.Init); e != nil {
						err = errutil.New("couldnt decode", patternName, paramName, e)
					} else {
						last.AddParam(&p)
					}
				case "bool_eval":
					p := pattern.BoolParam{Name: paramName}
					if e := decode(prog, &p.Init); e != nil {
						err = errutil.New("couldnt decode", patternName, paramName, e)
					} else {
						last.AddParam(&p)
					}
				default:
					// the type might be some sort of kind...
					if kind := kind.String; len(kind) > 0 {
						p := pattern.ObjectParam{Name: paramName, Kind: kind}
						if e := decode(prog, &p.Init); e != nil {
							err = errutil.New("couldnt decode", patternName, paramName, e)
						} else {
							last.AddParam(&p)
						}
					} else {
						err = errutil.Fmt("pattern %q parameter %q has unknown type %q ( expected an eval .)",
							patternName, paramName, typeName)
					}
				}
			}
			return
		},
		&patternName, &paramName, &typeName, &kind, &prog); e != nil {
		err = e
	} else {
		ret = out
	}
	return
}

func decode( prog []byte, out interface{}) (err error) {
	if len(prog) > 0 {
		dec := gob.NewDecoder(bytes.NewBuffer(prog))
		err = dec.Decode(out)
	}
	return
}

// collect the rules of all the various patterns and write them into the assembly
func buildPatterns(asm *Assembler) (err error) {
	if patterns, e := buildPatternCache(asm.cache.DB()); e != nil {
		err = e
	} else {
		err = buildPatternRules(asm, patterns)
	}
	return
}

func buildPatternRules(asm *Assembler, patterns patternCache) (err error) {
	var rules = []BuildRule{{
		Query: `select pattern, prog from asm_rule where type='bool_rule'`,
		NewContainer: func(name string) (ret interface{}) {
			if c, ok := patterns.init(name, "bool_eval"); ok {
				ret = &pattern.BoolPattern{CommonPattern: c}
			}
			return
		},
		NewEl: func(c interface{}) interface{} {
			pat := c.(*pattern.BoolPattern)
			pat.Rules = append(pat.Rules, &pattern.BoolRule{})
			return pat.Rules[len(pat.Rules)-1]
		},
	}, {
		Query: `select pattern, prog from asm_rule where type='number_rule'`,
		NewContainer: func(name string) (ret interface{}) {
			if c, ok := patterns.init(name, "number_eval"); ok {
				ret = &pattern.NumberPattern{CommonPattern: c}
			}
			return
		},
		NewEl: func(c interface{}) interface{} {
			pat := c.(*pattern.NumberPattern)
			pat.Rules = append(pat.Rules, &pattern.NumberRule{})
			return pat.Rules[len(pat.Rules)-1]
		},
	}, {
		Query: `select pattern, prog from asm_rule where type='text_rule'`,
		NewContainer: func(name string) (ret interface{}) {
			if c, ok := patterns.init(name, "text_eval"); ok {
				ret = &pattern.TextPattern{CommonPattern: c}
			}
			return
		},
		NewEl: func(c interface{}) interface{} {
			pat := c.(*pattern.TextPattern)
			pat.Rules = append(pat.Rules, &pattern.TextRule{})
			return pat.Rules[len(pat.Rules)-1]
		},
	}, {
		Query: `select pattern, prog from asm_rule where type='num_list_rule'`,
		NewContainer: func(name string) (ret interface{}) {
			if c, ok := patterns.init(name, "num_list_eval"); ok {
				ret = &pattern.NumListPattern{CommonPattern: c}
			}
			return
		},
		NewEl: func(c interface{}) interface{} {
			pat := c.(*pattern.NumListPattern)
			pat.Rules = append(pat.Rules, &pattern.NumListRule{})
			return pat.Rules[len(pat.Rules)-1]
		},
	}, {
		Query: `select pattern, prog from asm_rule where type='text_list_rule'`,
		NewContainer: func(name string) (ret interface{}) {
			if c, ok := patterns.init(name, "text_list_eval"); ok {
				ret = &pattern.TextListPattern{CommonPattern: c}
			}
			return
		},
		NewEl: func(c interface{}) interface{} {
			pat := c.(*pattern.TextListPattern)
			pat.Rules = append(pat.Rules, &pattern.TextListRule{})
			return pat.Rules[len(pat.Rules)-1]
		},
	}, {
		Query: `select pattern, prog from asm_rule where type='execute_rule'`,
		NewContainer: func(name string) (ret interface{}) {
			if c, ok := patterns.init(name, "execute"); ok {
				ret = &pattern.ActivityPattern{CommonPattern: c}
			}
			return
		},
		NewEl: func(c interface{}) interface{} {
			pat := c.(*pattern.ActivityPattern)
			pat.Rules = append(pat.Rules, &pattern.ExecuteRule{})
			return pat.Rules[len(pat.Rules)-1]
		},
	}}
	for _, rule := range rules {
		var name string
		var prog []byte
		if e := rule.buildFromRule(asm, &name, &prog); e != nil {
			err = e
			break
		}
	}
	return
}
