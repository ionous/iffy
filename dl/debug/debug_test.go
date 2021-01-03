package debug

import (
	"log"
	"strings"
	"testing"

	"github.com/ionous/iffy/dl/core"
)

func TestLog(t *testing.T) {
	w := log.Writer()
	defer log.SetOutput(w)
	var b strings.Builder
	log.SetOutput(&b)
	//
	lo := Log{Level: Warning, Value: &core.FromText{&core.Text{"hello"}}}
	if e := lo.Execute(nil); e != nil {
		t.Fatal(e)
	} else if got := b.String(); !strings.HasSuffix(got, " ### Warning: hello\n") {
		t.Fatalf("got %q", got)
	}
}
