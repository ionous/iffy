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
}

func testChart(t *testing.T, str, want string) (err error) {
	t.Log("test:", str)
	if res, e := Parse(str); e != nil {
		err = e
	} else {
		got := Blocks{res}.String()
		if got != want {
			err = errutil.New(want, "mismatched results", got)
		} else {
			t.Log("ok:", res)
		}
	}
	return
}
