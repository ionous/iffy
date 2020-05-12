package qna

import (
	"bytes"
	"encoding/gob"
	r "reflect"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/pattern"
	"github.com/ionous/iffy/tables"
)

func newRule(rd pattern.Rule) r.Value {
	rtype := r.TypeOf(rd).Elem()
	return r.New(rtype)
}
func newRuleSet(rd pattern.Rule) r.Value {
	rtype := r.TypeOf(rd.RuleDesc().RuleSet).Elem()
	return r.New(rtype).Elem() // rulesets are arrays, not pointers to arrays
}

func findRuleByTypeName(ruleType string) (ret pattern.Rule, err error) {
	found := false
	for _, rule := range pattern.Rules {
		if desc := rule.RuleDesc(); desc.Name == ruleType {
			ret, found = rule, true
			break
		}
	}
	if !found {
		err = errutil.New("unknown rule type", ruleType)
	}
	return
}

func (n *Fields) cacheRules(key keyType, pattern, patternType string) (ret interface{}, err error) {
	if rd, e := findRuleByTypeName(patternType); e != nil {
		err = e
	} else {
		var prog []byte
		if rows, e := n.db.Query(
			`select bytes 
					from mdl_rule mr
					join mdl_prog mp
					on (mr.idProg = mp.rowid)
					where mr.pattern = ?
					and mp.type = ?`,
			pattern, // object is the name of the patter
			patternType); e != nil {
			err = e
		} else {
			rs := newRuleSet(rd)
			if e := tables.ScanAll(rows, func() (err error) {
				rl := newRule(rd)
				dec := gob.NewDecoder(bytes.NewBuffer(prog))
				if e := dec.DecodeValue(rl); e != nil {
					err = e
				} else {
					rs = r.Append(rs, rl)
				}
				return
			}, &prog); e != nil {
				err = e
			} else {
				rules := rs.Interface()
				n.pairs[key] = rules
				ret = rules
			}
		}
	}
	return
}
