package main

import (
	"database/sql"
	"log"
	"os"
)

func ExampleKindData() {
	str := func(s string) sql.NullString {
		return sql.NullString{String: s, Valid: true}
	}
	templates.ExecuteTemplate(os.Stdout, "kindList", []Kind{
		{Name: "things", Path: "", Spec: "The things.",
			Props: []Prop{
				{"doe", "a deer", str("A female deer.")},
				{"flightiness", "flightless", str("An example aspect.")},
				{Name: "ray", Value: "5"},
			},
			Nouns: []string{"something", "someone"}},
	})

	// Output:
	// <h1>Kinds</h1>
	// <a href="#things">Things</a>.
	//
	// <h2 id="things">Things</h2>
	// <span>Parent kind: none.</span> <span class="spec">The things.</span>
	//
	// <h3>Properties</h3>
	// <dl>
	// 	<dt>Doe: <span>a deer.</span></dt><dd>A female deer.</dd>
	// 	<dt>Flightiness: <span>flightless.</span></dt><dd>An example aspect.</dd>
	// 	<dt>Ray: <span>5.</span></dt>
	// </dl>
	//
	// <h3>Nouns</h3>
	// 	<a href="/atlas/nouns#something">Something</a>,
	// 	<a href="/atlas/nouns#someone">Someone</a>.
}

func ExampleKindDB() {
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
	// <span>Parent kind: none.</span> <span class="spec">From inform: &#39;Represents anything interactive in the world. People, pieces of scenery, furniture, doors and mislaid umbrellas might all be examples, and so might more surprising things like the sound of birdsong or a shaft of sunlight.&#39;</span>
	//
	// <h3>Properties</h3>
	// <dl>
	// 	<dt>Brief: <span>&#34;&#34;.</span></dt><dd></dd>
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
