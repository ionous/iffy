package qna

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/ionous/iffy/assembly"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/object"
	"github.com/ionous/iffy/pattern"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/tables"
)

func SayIt(s string) rt.TextEval {
	return &core.Text{s}
}

type MatchNumber int

func (m MatchNumber) GetBool(run rt.Runtime) (okay bool, err error) {
	if v, e := run.GetVariable("num"); e != nil {
		err = e
	} else {
		n := int(v.(float64))
		okay = n == int(m)
	}
	return
}

// manually assemble a pattern database, and test that it works as expected.
func TestSayMe(t *testing.T) {
	gob.Register((*pattern.TextRule)(nil))
	gob.Register((*core.Text)(nil))
	gob.Register((*MatchNumber)(nil))

	db := newQnaDB(t, memory)
	defer db.Close()
	if e := tables.CreateModel(db); e != nil {
		t.Fatal(e)
	}
	m := assembly.NewAssembler(db)
	if e := WriteRules(m, "sayMe", pattern.TextRules{
		{nil, SayIt("Not between 1 and 3.")},
		{MatchNumber(3), SayIt("San!")},
		{MatchNumber(3), SayIt("Three!")},
		{MatchNumber(2), SayIt("Two!")},
		{MatchNumber(1), SayIt("One!")}}); e != nil {
		t.Fatal(e)
	}
	//
	run := NewRuntime(db)
	if p, e := run.GetField("sayMe", object.TextRule); e != nil {
		t.Fatal(e)
	} else if _, ok := p.(pattern.TextRules); !ok {
		t.Fatalf("not %T", p)
	} else {
		for i, expect := range []string{"One!", "Two!", "Three!", "Not between 1 and 3."} {
			det := core.DetermineText{
				"sayMe", &core.Parameters{[]*core.Parameter{{
					"num", &core.FromNum{
						&core.Number{float64(i + 1)},
					},
				}}},
			}
			if text, e := rt.GetText(run, &det); e != nil {
				t.Fatal(e)
			} else if expect != text {
				t.Fatal(i, text)
			} else {
				t.Log(text)
			}
		}
	}
}

func WriteRules(m *assembly.Assembler, name string, rules pattern.TextRules) (err error) {
	for _, rl := range rules {
		if e := WriteRule(m, name, rl); e != nil {
			err = e
			break
		}
	}
	return
}

func WriteRule(m *assembly.Assembler, name string, rule pattern.Rule) (err error) {
	d := rule.RuleDesc()
	if pid, e := WriteProg(m, d.Name, rule); e != nil {
		err = e
	} else {
		err = m.WriteRule(name, pid)
	}
	return

}

func WriteProg(m *assembly.Assembler, typeName string, prog interface{}) (ret int64, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if e := enc.Encode(prog); e != nil {
		err = e
	} else {
		ret, err = m.WriteProg(typeName, buf.Bytes())
	}
	return
}
