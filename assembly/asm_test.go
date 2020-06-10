package assembly

import (
	"database/sql"
	"os/user"
	"path"
	"strings"
	"testing"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/tables"
)

func getPath(file string) (ret string, err error) {
	if user, e := user.Current(); e != nil {
		err = e
	} else {
		ret = path.Join(user.HomeDir, file)
	}
	return
}

const memory = "file:test.db?cache=shared&mode=memory"

func newAssemblyTest(t *testing.T, path string) (ret *assemblyTest, err error) {
	var source string
	if len(path) > 0 {
		source = path
	} else {
		base := strings.Replace(t.Name(), "/", ".", -1) + ".db"
		if p, e := getPath(base); e != nil {
			err = errutil.New(e, "for", base)
		} else {
			source = p
		}
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
