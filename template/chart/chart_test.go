package chart

import (
	"github.com/ionous/errutil"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestChart(t *testing.T) {
	assert, x := testify.New(t), true
	x = x && assert.NoError(testChart(t, "hello world",
		"hello world"))
	x = x && assert.NoError(testChart(t, `hello {player}`,
		"hello {player:}")) // our tricky case: a word is a key.
	x = x && assert.NoError(testChart(t, `{player}, hello`,
		"{player:}, hello"))
	x = x && assert.NoError(testChart(t, `{player!}`,
		"{PLAYER/0}"))
	x = x && assert.NoError(testChart(t,
		`hello {player!} to the {"world"|cap?}.`,
		"hello {PLAYER/0} to the {`world` CAP/1}."))
	//
	x = x && assert.NoError(testChart(t, `{if player!}`,
		"{if:PLAYER/0}"))
	x = x && assert.NoError(testChart(t, `{if player}`, "{if:player}"))
}

func testChart(t *testing.T, str, want string) (err error) {
	t.Log("test:", str)
	if ds, e := Parse(str); e != nil {
		err = e
	} else if got := String(ds); want == ignoreResult {
		t.Log("got", got)
	} else if got != want {
		err = errutil.Fmt("want(%d): %s; != got(%d): %s.", len(want), want, len(got), got)
	} else {
		t.Log("ok", got)
	}
	return
}
