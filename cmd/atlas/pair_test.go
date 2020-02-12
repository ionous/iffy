package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/ionous/iffy/tables"
)

func ExamplePairData() {
	pairTemplate.Execute(os.Stdout, Pairing{
		Rel: &Relation{
			Name:        "authorship",
			Kind:        "people",
			Cardinality: tables.ONE_TO_MANY,
			OtherKind:   "books",
			Spec:        "Writers and their works.",
		},
		Pairs: []*Pair{
			{First: "N.K. Jemisin", Second: "The City We Became"},
			{First: "N.K. Jemisin", Second: "How Long 'Til Black Future Month"},
			{First: "Ted Chiang", Second: "Exhalation"},
		},
	})

	// Output:
	// <h1>Authorship</h1>
	// Relates people to many books. Writers and their works.
	// <table>
	// <tr>
	//   <td>N.K. Jemisin</td>
	//   <td>The City We Became</td>
	// </tr>
	// <tr>
	//   <td></td>
	//   <td>How Long 'Til Black Future Month</td>
	// </tr>
	// <tr>
	//   <td>Ted Chiang</td>
	//   <td>Exhalation</td>
	// </tr>
	// </table>
}

func ExamplePairDB() {
	const memory = "file:ExamplePairDB.db?cache=shared&mode=memory"
	if db, e := sql.Open("sqlite3", memory); e != nil {
		log.Fatalln("couldnt open db ", e)
	} else if e := createTestData(db); e != nil {
		log.Fatal("couldnt create test data ", e)
	} else if e := CreateAtlas(db); e != nil {
		log.Fatal("couldnt create atlas tables ", e)
	} else if e := listOfPairs(os.Stdout, "containing", db); e != nil {
		log.Fatal("couldnt process pairs ", e)
	}

	// Output:
	// <h1>Containing</h1>
	// Relates vehicles to many people.
	// <table>
	// <tr>
	//   <td>Dune Buggy</td>
	//   <td>Picard</td>
	// </tr>
	// <tr>
	//   <td>Enterprise</td>
	//   <td>Riker</td>
	// </tr>
	// </table>
}
