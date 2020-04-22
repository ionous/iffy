package tables

import (
	"io"
	"strings"
)

// ex. WriteCsv(db, os.Stdout, "select col1, col2 from table", 2)
func WriteCsv(db Query, w io.Writer, q string, cols int) (err error) {
	if rows, e := db.Query(q); e != nil {
		err = e
	} else {
		c := make([]string, cols)
		cp := make([]interface{}, cols)
		for i := 0; i < cols; i++ {
			cp[i] = &c[i]
		}
		err = ScanAll(rows, func() (err error) {
			io.WriteString(w, strings.Join(c, ","))
			io.WriteString(w, "\n")
			return
		}, cp...)
	}
	return
}
