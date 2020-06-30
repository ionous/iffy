package tables

import (
	"database/sql"
	"io"
)

// ex. WriteCsv(db, os.Stdout, "select col1, col2 from table", 2)
func WriteCsv(db Query, w io.Writer, q string, cols int) (err error) {
	if rows, e := db.Query(q); e != nil {
		err = e
	} else {
		c := make([]sql.NullString, cols)
		cp := make([]interface{}, cols)
		for i := 0; i < cols; i++ {
			cp[i] = &c[i]
		}
		err = ScanAll(rows, func() (err error) {
			for i, col := range c {
				if i > 0 {
					io.WriteString(w, ",")
				}
				if !col.Valid {
					io.WriteString(w, "NULL")
				} else {
					io.WriteString(w, col.String)
				}
			}
			io.WriteString(w, "\n")
			return
		}, cp...) // pass the pointers to the column strings
	}
	return
}
