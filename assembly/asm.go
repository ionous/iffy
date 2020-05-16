package assembly

import "database/sql"

func AssembleModel(db *sql.DB) (err error) {
	// [-] adds relations between kinds
	// [-] creates instances
	// [-] sets instance properties
	// [-] relates instances
	// [] makes action handlers
	// [] makes event listeners
	// [] computes aliases
	// [] sets up printed name property
	// - backtracing to source:
	// ex. each "important" table entry gets an separate entry pointing back to original source
	//
	m := NewModeler(db)
	if e := AssembleAncestry(m, db, "things"); e != nil {
		err = e
	} else if e := AssembleAspects(m, db); e != nil {
		err = e
	} else if e := AssembleFields(m, db); e != nil {
		err = e
	} else if e := AssembleRelations(m, db); e != nil {
		err = e
	} else if e := AssembleDefaults(m, db); e != nil {
		err = e
	} else if e := AssembleNouns(m, db); e != nil {
		err = e
	} else if e := AssembleRelatives(m, db); e != nil {
		err = e
	} else if e := AssembleValues(m, db); e != nil {
		err = e
	} else if e := AssemblePatterns(m, db); e != nil {
		err = e
	} else if _, e := db.Exec("insert into mdl_prog select type, bytes from eph_prog;" +
		"insert into mdl_check select * from asm_check"); e != nil {
		err = e
	}
	return
}
