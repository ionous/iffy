package assembly

import "database/sql"

func AssembleStory(db *sql.DB, baseKind string, reporter IssueReport) (err error) {
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
	asm := NewAssemblerReporter(db, reporter)
	if e := AssembleAncestry(asm, baseKind); e != nil || asm.IssueCount > 0 {
		err = e
	} else if e := AssembleAspects(asm); e != nil || asm.IssueCount > 0 {
		err = e
	} else if e := AssembleFields(asm); e != nil || asm.IssueCount > 0 {
		err = e
	} else if e := AssembleRelations(asm); e != nil || asm.IssueCount > 0 {
		err = e
	} else if e := AssembleDefaults(asm); e != nil || asm.IssueCount > 0 {
		err = e
	} else if e := AssembleNouns(asm); e != nil || asm.IssueCount > 0 {
		err = e
	} else if e := AssembleRelatives(asm); e != nil || asm.IssueCount > 0 {
		err = e
	} else if e := AssembleValues(asm); e != nil || asm.IssueCount > 0 {
		err = e
	} else if e := AssemblePatterns(asm); e != nil || asm.IssueCount > 0 {
		err = e
	} else if _, e := db.Exec(
		`insert into mdl_prog select type, prog as bytes from eph_prog;
		insert into mdl_check select * from asm_check`); e != nil {
		err = e
	}
	return
}
