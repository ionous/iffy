package assembly

import (
	"bytes"
	"database/sql"
	"encoding/gob"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/affine"
	"github.com/ionous/iffy/dl/pattern"
	"github.com/ionous/iffy/dl/term"
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
	patternName string          // name of the pattern
	patternType string          // "return" type of the pattern
	prologue    []term.Preparer // list of all parameters sent to the pattern
	locals      []term.Preparer // ...
	returns     term.Preparer
}

func (pat *patternEntry) AddParam(cat string, param term.Preparer) (err error) {
	switch cat {
	case tables.NAMED_PARAMETER:
		pat.prologue = append(pat.prologue, param)
	case tables.NAMED_LOCAL:
		pat.locals = append(pat.locals, param)
	case tables.NAMED_RETURN:
		// fix? check for multiple / different sets?
		pat.returns = param
	default:
		err = errutil.New("unknown category", cat)
	}
	return
}

type patternCache map[string]*patternEntry

func (cache patternCache) init(name, patternType string) (ret *pattern.Pattern, okay bool) {
	if c, ok := (cache)[name]; ok && c.patternType == patternType {
		ret = &pattern.Pattern{
			Name:    name,
			Params:  c.prologue,
			Locals:  c.locals,
			Returns: c.returns,
		}
		okay = true
	}
	return
}

// read pattern declarations from the ephemera db
func buildPatternCache(db *sql.DB) (ret patternCache, err error) {
	// build the pattern cache
	out := make(patternCache)
	var patternName, paramName, category, typeName string
	var kind, affinity sql.NullString
	var prog []byte
	var last *patternEntry
	if e := tables.QueryAll(db,
		`select ap.pattern, ap.param, ap.cat, ap.type, ap.affinity, ap.kind, ep.prog
		from asm_pattern_decl ap
		left join eph_prog ep
		on (ep.rowid = ap.idProg)`,
		func() (err error) {
			// fix: need to handle conflicting prog definitions
			// fix: should probably watch for locals which shadow parameter names
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
				// it also seems a bad time to be camelizing things.
				paramName := lang.Breakcase(paramName)

				// locals have simple type names, parameters are still using _eval.
				var p term.Preparer
				switch typeName {
				case "text_eval", "text":
					p = &term.Text{Name: paramName}
				case "number_eval", "number":
					p = &term.Number{Name: paramName}
				case "bool_eval", "bool":
					p = &term.Bool{Name: paramName}
				case "text_list", "text_list_eval":
					p = &term.TextList{Name: paramName}
				case "num_list", "num_list_eval":
					p = &term.NumList{Name: paramName}
				default:
					// the type might be some sort of kind...
					if kind := kind.String; len(kind) > 0 {
						switch aff := affinity.String; aff {
						case string(affine.Object):
							p = &term.Object{Name: paramName, Kind: kind}
						case string(affine.Record):
							p = &term.Record{Name: paramName, Kind: kind}
						case string(affine.RecordList):
							p = &term.RecordList{Name: paramName, Kind: kind}
						}
					}
				}
				if p != nil {
					err = last.AddParam(category, p)
				} else {
					err = errutil.Fmt("pattern %q parameter %q has unknown type %q(%s)",
						patternName, paramName, typeName, affinity.String)
				}
			}
			return
		},
		&patternName, &paramName, &category, &typeName, &affinity, &kind, &prog); e != nil {
		err = e
	} else {
		ret = out
	}
	return
}

func decode(prog []byte, out interface{}) (err error) {
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

func buildPatternRules(asm *Assembler, patterns patternCache) error {
	var name string
	var prog []byte
	rule := BuildRule{
		Query: `select pattern, prog from asm_rule where type='rule'`,
		NewContainer: func(name string) (ret interface{}) {
			if c, ok := patterns.init(name, "execute"); ok {
				ret = c
			}
			return
		},
		NewEl: func(c interface{}) interface{} {
			pat := c.(*pattern.Pattern)
			pat.Rules = append(pat.Rules, &pattern.Rule{})
			return pat.Rules[len(pat.Rules)-1]
		},
	}
	return rule.buildFromRule(asm, &name, &prog)
}
