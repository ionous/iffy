package assembly

import (
	"database/sql"
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/tables"
	"github.com/ionous/iffy/test/testdb"
)

func newAssemblyTest(t *testing.T, path string) (ret *assemblyTest, err error) {
	var source string
	if len(path) > 0 {
		source = path
	} else if p, e := testdb.PathFromName(t.Name()); e != nil {
		err = e
	} else {
		source = p
	}
	if err == nil {
		if db, e := sql.Open(SqlCustomDriver, source); e != nil {
			err = errutil.New(e, "for", source)
		} else if e := tables.CreateEphemera(db); e != nil {
			err = errutil.New(e, "for", source)
		} else if e := tables.CreateAssembly(db); e != nil {
			err = errutil.New(e, "for", source)
		} else if e := tables.CreateModel(db); e != nil {
			err = errutil.New(e, "for", source)
		} else {
			ds := new(Dilemmas)
			rec := ephemera.NewRecorder(t.Name(), db)
			mdl := NewAssemblerReporter(db, func(pos reader.Position, msg string) {
				t.Log("report:", pos, msg)
				ds.Add(pos, msg)
			})
			ret = &assemblyTest{
				T:         t,
				db:        db,
				rec:       rec,
				assembler: mdl,
				dilemmas:  ds,
			}
		}
	}
	return
}

type assemblyTest struct {
	*testing.T
	db        *sql.DB
	rec       *ephemera.Recorder
	assembler *Assembler
	dilemmas  *Dilemmas
}
