package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/ionous/iffy/tables"
)

func ExamplePairData() {
	templates.ExecuteTemplate(os.Stdout, "pairList", Pairing{
		Rel: &Relation{
			Name:        "authorship",
			Kind:        "people",
			Cardinality: tables.ONE_TO_MANY,
			OtherKind:   "books",
			Spec:        "Writers and their works.",
		},
		Pairs: []*Pair{
			{First: "n. k. jemisin", Second: "the city we became"},
			{First: "n. k. jemisin", Second: "how long til black future month"},
			{First: "rudy rucker", Second: "realware"},
		},
	})

	// Output:
	// <h1>Authorship</h1>
	// Relates <a href="/atlas/kinds#people">People</a> to many <a href="/atlas/kinds#books">Books</a>.
	//  Writers and their works.
	// <table>
	// <tr>
	//   <td><a href="/atlas/nouns#n.-k.-jemisin">N. K. Jemisin</a></td>
	//   <td><a href="/atlas/nouns#the-city-we-became">The City We Became</a></td>
	// </tr>
	// <tr>
	//   <td></td>
	//   <td><a href="/atlas/nouns#how-long-til-black-future-month">How Long Til Black Future Month</a></td>
	// </tr>
	// <tr>
	//   <td><a href="/atlas/nouns#rudy-rucker">Rudy Rucker</a></td>
	//   <td><a href="/atlas/nouns#realware">Realware</a></td>
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
	// Relates <a href="/atlas/kinds#vehicles">Vehicles</a> to many <a href="/atlas/kinds#people">People</a>.
	//  The outside of insides.
	// <table>
	// <tr>
	//   <td><a href="/atlas/nouns#dune-buggy">Dune Buggy</a></td>
	//   <td><a href="/atlas/nouns#picard">Picard</a></td>
	// </tr>
	// <tr>
	//   <td><a href="/atlas/nouns#enterprise">Enterprise</a></td>
	//   <td><a href="/atlas/nouns#riker">Riker</a></td>
	// </tr>
	// </table>
}
