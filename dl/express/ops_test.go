package express

import (
	"bytes"
	"io"
	r "reflect"
	"testing"

	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/rules"
	"github.com/ionous/iffy/ident"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/spec/ops"
)

// TestOps to verify iffy statements.
// . text expressions: tested elsewhere
// . number, etc. expressions
// . direct calls: "{Story.score|testScore!}"
// . pattern determinations: c.Cmd("set text", "story", "status left", "{playerSurroundings!}")
// . pattern evaluations:  c.Cmd("{print banner text!}"
//
func TestOps(t *testing.T) {
	classes := make(unique.Types)
	cmds := ops.NewOps(classes)
	patterns := unique.NewStack(cmds.ShadowTypes)

	// a sample command.
	type TestScore struct {
		rt.NumberEval
	}
	type Pluralize struct {
		rt.TextEval
	}
	// a sample, empty, pattern.
	type TestPrint struct {
		Story ident.Id `if:"cls:kind"`
	}

	unique.PanicTypes(cmds,
		(*TestScore)(nil),
		(*Pluralize)(nil),
	)
	unique.PanicBlocks(cmds,
		(*core.Commands)(nil),
		(*rules.RuntimeCmds)(nil), // determine.
		(*Commands)(nil))

	unique.PanicTypes(patterns,
		(*TestPrint)(nil))

	xform := NewTransform(cmds, nil)

	t.Run("property", func(t *testing.T) {
		var root struct{ rt.NumberEval }
		c := cmds.NewBuilder(&root, xform)
		//
		// TODO: document how this works.
		//
		c.Cmd("test score", "{5 + 5}")
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else {
			testEqual(t, &TestScore{
				NumberEval: &core.SumOf{
					A: &core.NumValue{Num: 5},
					B: &core.NumValue{Num: 5},
				},
			}, root.NumberEval)
		}
	})
	//
	t.Run("global", func(t *testing.T) {
		var root struct{ rt.NumberEval }
		c := cmds.NewBuilder(&root, xform)
		c.Cmd("test score", "{Story.score}")
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else {
			//
			// FIX: since the numeval works above, why do we need render? why not just GetAt?
			//
			testEqual(t, &TestScore{&Render{
				&core.ObjectName{"Story"}, "score"},
			}, root.NumberEval)
		}
	})
	//
	t.Run("run pipe", func(t *testing.T) {
		var root struct{ rt.NumberEval }
		c := cmds.NewBuilder(&root, xform)
		c.Cmd("{Story.score|testScore!}")
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else {
			testEqual(t, &TestScore{&Render{
				&core.ObjectName{"Story"}, "score"},
			}, root.NumberEval)
		}
	})
	t.Run("run params", func(t *testing.T) {
		var root struct{ rt.NumberEval }
		c := cmds.NewBuilder(&root, xform)
		c.Cmd("{testScore! Story.score}")
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else {
			testEqual(t, &TestScore{&Render{
				&core.ObjectName{"Story"}, "score"},
			}, root.NumberEval)
		}
	})
	t.Run("determine", func(t *testing.T) {
		var root struct{ rt.TextEval }
		c := cmds.NewBuilder(&root, xform)
		c.Cmd("{testPrint! story}")
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else {
			testStrings(t, root.TextEval,
				`Determine{&TestPrint{GetAt{"story"}}}`)
		}
	})
	t.Run("dpipe", func(t *testing.T) {
		var root struct{ rt.TextEval }
		c := cmds.NewBuilder(&root, xform)
		c.Cmd("{testPrint: story|buffer:|pluralize:}")
		if e := c.Build(); e != nil {
			t.Fatal(e)
		} else {
			testStrings(t, root.TextEval,
				`Pluralize{Buffer{Determine{&TestPrint{GetAt{"story"}}}}}`)
		}
	})
}

func testStrings(t *testing.T, res interface{}, want string) {
	// got := pretty.Sprintf("%#v", res)
	got := dump(res)
	t.Log("got :", got)
	if got != want {
		t.Fatalf("want: %s", want)
	}
}

// a very hacky string representation of an iffy call tree
// kr/pretty doesnt handle shadow class well.
// you have to add a custom formatter as a *class* member :\
// causing new dependencies on fmt,io,pretty,and more.
func dump(i interface{}) string {
	var b bytes.Buffer
	dummy(&b, r.ValueOf(i).Elem())
	return b.String()
}

func dummy(w io.Writer, v ops.Target) {
	vt := v.Type()
	io.WriteString(w, vt.Name())
	io.WriteString(w, "{")
	for i := 0; i < v.NumField(); i++ {
		if i > 0 {
			io.WriteString(w, ",")
		}
		if fv := v.Field(i); fv.IsValid() {
			field := vt.Field(i)
			switch fv.Kind() {
			case r.Slice:
				if field.Type.Elem().Kind() == r.Interface {
					for j := 0; j < fv.Len(); j++ {
						dummy(w, fv.Index(j).Elem().Elem())
					}
				} else {
					io.WriteString(w, "[]")

				}
			case r.Interface:
				if el := fv.Elem(); el.Type() == r.TypeOf((*ops.ShadowClass)(nil)) {
					shade := el.Interface().(*ops.ShadowClass)
					io.WriteString(w, "&")
					dummy(w, shade)
				} else {
					dummy(w, el.Elem())
				}

			case r.String:
				io.WriteString(w, "\""+fv.String()+"\"")
			default:
				io.WriteString(w, fv.Kind().String())
			}
		}
	}
	io.WriteString(w, "}")
}
