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

func xExamplePairDB() {
	const memory = "file:ExampleKindDB.db?cache=shared&mode=memory"
	if db, e := sql.Open("sqlite3", memory); e != nil {
		log.Fatalln("couldnt open db ", e)
	} else if e := createTestData(db); e != nil {
		log.Fatal("couldnt create test data ", e)
	} else if e := CreateAtlas(db); e != nil {
		log.Fatal("couldnt create atlas tables ", e)
	} else if e := listOfKinds(os.Stdout, db); e != nil {
		log.Fatal("couldnt process kinds ", e)
	}

	// Output:
	// <h1>Kinds</h1>
	// <a href="#things">Things</a>, <a href="#people">People</a>, <a href="#vehicles">Vehicles</a>, <a href="#cars">Cars</a>.
	//
	// <h2 id="things">Things</h2>
	// <span>Parent kind: none.</span> <span class="spec">From inform: 'Represents anything interactive in the world. People, pieces of scenery, furniture, doors and mislaid umbrellas might all be examples, and so might more surprising things like the sound of birdsong or a shaft of sunlight.'</span>
	//
	// <h3>Properties</h3>
	// <dl>
	// 	<dt>Brief: <span>"".</span></dt><dd></dd>
	// </dl>
	// <h2 id="people">People</h2>
	// <span>Parent kind: <a href="#things">Things</a>.</span> <span class="spec"></span>
	//
	// <h3>Nouns</h3>
	// 	<a href="/atlas/nouns#picard">Picard</a>,
	// 	<a href="/atlas/nouns#riker">Riker</a>.
	//
	// <h2 id="vehicles">Vehicles</h2>
	// <span>Parent kind: <a href="#things">Things</a>.</span> <span class="spec"></span>
	//
	// <h3>Properties</h3>
	// <dl>
	// 	<dt>Flightiness: <span>flight worthy.</span></dt><dd></dd>
	// </dl>
	//
	// <h3>Nouns</h3>
	// 	<a href="/atlas/nouns#enterprise">Enterprise</a>.
	//
	// <h2 id="cars">Cars</h2>
	// <span>Parent kind: <a href="#vehicles">Vehicles</a>.</span> <span class="spec"></span>
	//
	// <h3>Properties</h3>
	// <dl>
	// 	<dt>Flightiness: <span>flightless.</span></dt>
	// 	<dt>Num Wheels: <span>4.</span></dt><dd>Not all cars are created equal, or even even.</dd>
	// </dl>
	//
	// <h3>Nouns</h3>
	// 	<a href="/atlas/nouns#dune-buggy">Dune Buggy</a>.
}
