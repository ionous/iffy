package chart

import (
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestPipe(t *testing.T) {
	assert := testify.New(t)
	x := true
	x = x && assert.NoError(testPipe(t, "", "")) // arguments are optional.
	x = x && assert.NoError(testPipe(t, "world", "[world]"))
	x = x && assert.NoError(testPipe(t, "up!|up!", "UP/0UP/1"))
	x = x && assert.NoError(testPipe(t, "hello!", "HELLO/0"))
	x = x && assert.NoError(testPipe(t, "world|hello!", "[world]HELLO/1"))
	x = x && assert.NoError(testPipe(t, "world | hello! ", "[world]HELLO/1"))
	x = x && assert.NoError(testPipe(t, "world|hello! there", "[there][world]HELLO/2"))
	x = x && assert.NoError(testPipe(t, "world|capitalize!|hello: there", "[there][world]CAPITALIZE/1HELLO/2"))
}

func testPipe(t *testing.T, str string, want string) (err error) {
	var p PipeParser
	t.Logf("parsing: '%s'", str)
	if e := parse(&p, str); e != nil {
		err = e
	} else if exp, e := p.GetExpression(); e != nil {
		err = e
	} else {
		res := exp.String()
		t.Logf("got: '%s'", res)
		if diff := pretty.Diff(res, want); len(diff) > 0 {
			err = errutil.Fmt("wanted: '%s'", want)
		}
	}
	return
}
