package assembly

import (
	"database/sql"

	"github.com/ionous/iffy"
)

func AssembleStory(db *sql.DB, baseKind string, reporter IssueReport) (err error) {
	iffy.RegisterGobs()
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
	if e := AssemblePlurals(asm); e != nil || asm.IssueCount > 0 {
		err = e
	} else if e := AssembleAncestry(asm, baseKind); e != nil || asm.IssueCount > 0 {
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
	} else if e := AssembleTests(asm); e != nil || asm.IssueCount > 0 {
		err = e
	}
	return
}
