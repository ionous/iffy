package assembly

import _ "github.com/mattn/go-sqlite3"

// 	porterstemmer "github.com/reiver/go-porterstemmer"

// func stem(s string) string {
// 	return porterstemmer.StemString(s)
// }

// func init() {
// 	sql.Register("sqlite3_custom", &sqlite.SQLiteDriver{
// 		ConnectHook: func(conn *sqlite.SQLiteConn) (err error) {
// 			if e := conn.RegisterFunc("stem", stem, true); e != nil {
// 				err = e
// 			}
// 			return
// 		},
// 	})
// }
