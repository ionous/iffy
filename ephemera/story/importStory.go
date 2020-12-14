package story

import (
	"database/sql"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy"
	"github.com/ionous/iffy/dl/core"
	"github.com/ionous/iffy/ephemera/decode"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/kr/pretty"
)

func ImportStory(src string, db *sql.DB, m reader.Map, reporter decode.IssueReport) (err error) {
	return ImportStories(src, db, []reader.Map{m}, reporter)
}

func ImportStories(src string, db *sql.DB, ms []reader.Map, reporter decode.IssueReport) (err error) {
	iffy.RegisterGobs()
	dec := decode.NewDecoderReporter(src, reporter)
	k := NewImporterDecoder(src, db, dec)
	//
	for _, slats := range iffy.AllSlats {
		dec.AddDefaultCallbacks(slats)
	}
	dec.AddDefaultCallbacks(core.Slats)
	// add ops from iffy.js, including golang generated stubs via stubs.js
	// anything that implements ImportStub() will get processed during ReadSpec.
	k.AddModel(Model)
	//
	for _, m := range ms {
		if i, e := dec.ReadSpec(m); e != nil {
			err = e
			break
		} else {
			pretty.Println(i)
			if story, ok := i.(*Story); !ok {
				err = errutil.Fmt("imported spec wasn't a story %T", i)
				break
			} else if e := story.ImportStory(k); e != nil {
				err = e
				break
			}
		}
	}
	return
}
