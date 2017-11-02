package chart

import (
	"github.com/ionous/errutil"
	"github.com/ionous/iffy/template/postfix"
	"github.com/kr/pretty"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestCall(t *testing.T) {
	assert := testify.New(t)
	x := true
	x = x && assert.NoError(testCall(t, "fun!", "fun"))
	x = x && assert.NoError(testCall(t, "args: a b", "args", "a", "b"))
	x = x && assert.NoError(testCall(t, "quest?", "quest"))
}

func testCall(t *testing.T, str string, name string, args ...string) (err error) {
	p := MakeCallParser(&MatchFactory{args})
	t.Logf("parsing: '%s'", str)
	if e := parse(&p, str); e != nil {
		err = e
	} else if res, e := p.GetExpression(); e != nil {
		err = e
	} else {
		var want postfix.Expression
		for _, str := range args {
			want = append(want, Quote(str))
		}
		want = append(want, Command{name, len(args)})
		//
		t.Log("got", res)
		if diff := pretty.Diff(res, want); len(diff) > 0 {
			err = errutil.New("want", want)
		}
	}
	return
}
