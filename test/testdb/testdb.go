package testdb

import (
	"database/sql"
	"io"
	"os/user"
	"path"
	"strings"

	"github.com/ionous/errutil"
	"github.com/ionous/iffy/tables"
	"github.com/ionous/sliceOf"
)

const Memory = "file:test.db?cache=shared&mode=memory"

// from testing.T.Name() return a local db testing path
func PathFromName(name string) (ret string, err error) {
	rest := strings.Replace(name, "/", ".", -1) + ".db"
	if user, e := user.Current(); e != nil {
		err = errutil.New(e, "for", name)
	} else {
		ret = path.Join(user.HomeDir, rest)
	}
	return
}

// given the table name, and column names return a super-secret helper
func TableCols(table_cols ...string) []string {
	return table_cols
}

// compatible with sql.DB for use with caches, etc.
type Executer interface {
	Exec(q string, args ...interface{}) (sql.Result, error)
}

// compatible with sql.DB for use with caches, etc.
type Querier interface {
	Query(q string, args ...interface{}) (*sql.Rows, error)
}

func Ins(db Executer, tablecols []string, els ...interface{}) (err error) {
	ins, width := tables.Insert(tablecols[0], tablecols[1:]...), len(tablecols)-1
	for i, cnt := 0, len(els); i < cnt; i += width {
		if _, e := db.Exec(ins, els[i:i+width]...); e != nil {
			err = e
			break
		}
	}
	return
}

func WriteCsv(db Querier, w io.Writer, tablecols []string, where string) (err error) {
	table, cols := tablecols[0], strings.Join(tablecols[1:], ", ")
	q := strings.Join(sliceOf.String("select", cols, "from", table, where, "order by", cols), " ")
	return tables.WriteCsv(db, w, q, len(tablecols)-1)
}
