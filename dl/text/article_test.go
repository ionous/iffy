package text_test

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/dl/text"
	"github.com/ionous/iffy/ops"
	"github.com/ionous/iffy/ref"
	"github.com/ionous/iffy/ref/unique"
	"github.com/ionous/iffy/rt"
	"github.com/ionous/iffy/rtm"
	. "github.com/ionous/iffy/tests"
	"github.com/ionous/sliceOf"
	"github.com/stretchr/testify/suite"
	"testing"
)

var kinds = []interface{}{
	&Kind{Name: "Lamp-post"},
	&Kind{Name: "Soldiers", IndefiniteArticle: "some"},
	&Kind{Name: "trevor", CommonProper: ProperNamed},
}

func TestArticles(t *testing.T) {
	suite.Run(t, new(ArticleSuite))
}

// Regular expression to select test suites specified command-line argument "-run". Regular expression to select the methods of test suites specified command-line argument "-m"
type ArticleSuite struct {
	suite.Suite
	ops   *ops.Ops
	test  *testing.T
	run   rt.Runtime
	lines rtm.LineWriter
}

func (t *ArticleSuite) Lines() (ret []string) {
	ret = t.lines.Lines()
	t.lines = rtm.LineWriter{}
	return
}

func (t *ArticleSuite) SetupTest() {
	errutil.Panic = false
	t.ops = ops.NewOps(
		(*text.Commands)(nil),
		(*core.Commands)(nil),
	)
	t.test = t.T()

	/// need a model finder
	cls := make(ref.Classes)
	mm := unique.PanicRegistry(cls)
	unique.RegisterType(mm, (*Kind)(nil))
	//
	if objs, e := cls.MakeModel(kinds); e != nil {
		panic(e)
	} else {
		t.run = rtm.NewRtm(cls, objs, ref.MakeRelations(cls))
		t.run.PushWriter(&t.lines)
	}
}
func (t *ArticleSuite) match(expected string, run func(c *ops.Builder)) {
	var root struct{ Eval rt.Execute }
	if c, ok := t.ops.NewBuilder(&root); ok {
		run(c)

		if _, e := c.Build(); t.NoError(e) {
			if e := root.Eval.Execute(t.run); t.NoError(e) {
				lines := t.Lines()
				t.Equal(sliceOf.String(expected), lines)
			}
		}
	}
}

// lower a/n
func (t *ArticleSuite) TestATrailingLampPost() {
	t.match("You can only just make out a lamp-post.",
		func(c *ops.Builder) {
			if c.Cmd("print line").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", "You can only just make out")
					c.Cmd("print text", c.Cmd("lower a/n", "lamp post"))
					c.Cmd("print text", ".")
					c.End()
				}
				c.End()
			}
		})
}

func (t *ArticleSuite) TestATrailingTrevor() {
	t.match("You can only just make out Trevor.",
		func(c *ops.Builder) {
			if c.Cmd("print line").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", "You can only just make out")
					c.Cmd("print text", c.Cmd("lower a/n", "trevor"))
					c.Cmd("print text", ".")
					c.End()
				}
				c.End()
			}
		})
}

func (t *ArticleSuite) TestATrailingSoldiers() {
	t.match("You can only just make out some soldiers.",
		func(c *ops.Builder) {
			if c.Cmd("print line").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", "You can only just make out")
					c.Cmd("print text", c.Cmd("lower a/n", "soldiers"))
					c.Cmd("print text", ".")
					c.End()
				}
				c.End()
			}
		})
}

// upper a/n
func (t *ArticleSuite) TestALeadingLampPost() {
	t.match("A lamp-post can be made out in the mist.",
		func(c *ops.Builder) {
			if c.Cmd("print line").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", c.Cmd("upper a/n", "lamp post"))
					c.Cmd("print text", "can be made out in the mist.")
					c.End()
				}
				c.End()
			}
		})
}

func (t *ArticleSuite) TestALeadingTrevor() {
	t.match("Trevor can be made out in the mist.",
		func(c *ops.Builder) {
			if c.Cmd("print line").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", c.Cmd("upper a/n", "trevor"))
					c.Cmd("print text", "can be made out in the mist.")
					c.End()
				}
				c.End()
			}
		})
}

func (t *ArticleSuite) TestALeadingSoldiers() {
	t.match("Some soldiers can be made out in the mist.",
		func(c *ops.Builder) {
			if c.Cmd("print line").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", c.Cmd("upper a/n", "soldiers"))
					c.Cmd("print text", "can be made out in the mist.")
					c.End()
				}
				c.End()
			}
		})
}

// lower-the
func (t *ArticleSuite) TestTheTrailingLampPost() {
	t.match("You can only just make out the lamp-post.",
		func(c *ops.Builder) {
			if c.Cmd("print line").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", "You can only just make out")
					c.Cmd("print text", c.Cmd("lower the", "lamp post"))
					c.Cmd("print text", ".")
					c.End()
				}
				c.End()
			}
		})
}

func (t *ArticleSuite) TestTheTrailingTrevor() {
	t.match("You can only just make out Trevor.",
		func(c *ops.Builder) {
			if c.Cmd("print line").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", "You can only just make out")
					c.Cmd("print text", c.Cmd("lower the", "trevor"))
					c.Cmd("print text", ".")
					c.End()
				}
				c.End()
			}
		})
}

func (t *ArticleSuite) TestTheTrailingSoldiers() {
	t.match("You can only just make out the soldiers.",
		func(c *ops.Builder) {
			if c.Cmd("print line").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", "You can only just make out")
					c.Cmd("print text", c.Cmd("lower the", "soldiers"))
					c.Cmd("print text", ".")
					c.End()
				}
				c.End()
			}
		})
}

// uppe the
func (t *ArticleSuite) TestTheLeadingLampPost() {
	t.match("The lamp-post may be a trick of the mist.",
		func(c *ops.Builder) {
			if c.Cmd("print line").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", c.Cmd("upper the", "lamp post"))
					c.Cmd("print text", "may be a trick of the mist.")
					c.End()
				}
				c.End()
			}
		})
}

func (t *ArticleSuite) TestTheLeadingTrevor() {
	t.match("Trevor may be a trick of the mist.",
		func(c *ops.Builder) {
			if c.Cmd("print line").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", c.Cmd("upper the", "trevor"))
					c.Cmd("print text", "may be a trick of the mist.")
					c.End()
				}
				c.End()
			}
		})
}

func (t *ArticleSuite) TestTheLeadingSoldiers() {
	t.match("The soldiers may be a trick of the mist.",
		func(c *ops.Builder) {
			if c.Cmd("print line").Begin() {
				if c.Cmds().Begin() {
					c.Cmd("print text", c.Cmd("upper the", "soldiers"))
					c.Cmd("print text", "may be a trick of the mist.")
					c.End()
				}
				c.End()
			}
		})
}
