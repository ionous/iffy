package qna

import (
	"github.com/ionous/iffy/pattern"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/tables"
)

type programMap map[string]tables.Prog

var programs = []tables.Prog{{
	"bool_eval",
	(*rt.BoolEval)(nil),
}, {
	"number_eval",
	(*rt.NumberEval)(nil),
}, {
	"text_eval",
	(*rt.TextEval)(nil),
}, {
	"num_list_eval",
	(*rt.NumListEval)(nil),
}, {
	"text_list_eval",
	(*rt.TextListEval)(nil),
}, {
	"execute",
	(*rt.Execute)(nil),
}, {
	"bool_rule",
	(*pattern.BoolRule)(nil),
}, {
	"number_rule",
	(*pattern.NumberRule)(nil),
}, {
	"text_rule",
	(*pattern.TextRule)(nil),
}, {
	"num_list_rule",
	(*pattern.NumListRule)(nil),
}, {
	"text_list_rule",
	(*pattern.TextListRule)(nil),
}, {
	"execute_rule",
	(*pattern.ExecuteRule)(nil),
}}

// given a name ( ex. text_rule ") return its program fragment info
func findProgByName(progType string) (ret *tables.Prog, okay bool) {
	if x := len(progType); x > 1 {
		progType := progType[1:]
		for _, p := range programs {
			if p.Type == progType {
				ret, okay = &p, true
				break
			}
		}
	}
	return
}

func (n *Fields) getAggregatedProg(key keyType, p *tables.Prog) (ret interface{}, err error) {
	if rows, e := n.progBytes.Query(key.owner, p.Type); e != nil {
		err = e
	} else if prog, e := p.Aggregate(rows); e != nil {
		err = e
	} else {
		n.pairs[key] = prog
		ret = prog
		// pretty.Println(prog)
	}
	return
}
