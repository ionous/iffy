package main

import (
	"database/sql"
	"log"
	"os"
)

func ExampleAspectData() {
	templates.ExecuteTemplate(os.Stdout, "aspectList", []Aspect{{
		Name:  "flightiness",
		Kinds: []string{"vehicles", "birds"},
		Traits: []Trait{{
			Name: "flightless",
			Spec: "grounded.",
		}, {
			Name: "glide worthy",
			Spec: "Better at landing than taking off.",
		}, {
			Name: "flight worthy",
		}},
	}, {
		Name:  "bool",
		Spec:  "example binary aspect.",
		Kinds: []string{"things"},
		Traits: []Trait{{
			Name: "true",
		}, {
			Name: "false",
		}},
	}})

	// Output:
	// <h1>Aspects</h1>
	//   <a href="#flightiness">Flightiness</a>,
	//   <a href="#bool">Bool</a>.
	// <h2 id="flightiness">Flightiness</h2>
	// <h3>Kinds</h3>
	//   <a href="/atlas/kinds#vehicles">Vehicles</a>,
	//   <a href="/atlas/kinds#birds">Birds</a>.
	// <h3>Traits</h3>
	// <dl>
	//   <dt>Flightless</dt>
	//    <dd>grounded.</dd>
	//   <dt>Glide Worthy</dt>
	//    <dd>Better at landing than taking off.</dd>
	//   <dt>Flight Worthy</dt>
	// </dl>
	// <h2 id="bool">Bool</h2>
	// example binary aspect.
	// <h3>Kinds</h3>
	//   <a href="/atlas/kinds#things">Things</a>.
	// <h3>Traits</h3>
	// <dl>
	//   <dt>True</dt>
	//   <dt>False</dt>
	// </dl>
}

func ExampleAspectDB() {
	const memory = "file:ExampleAspectDB.db?cache=shared&mode=memory"
	if db, e := sql.Open("sqlite3", memory); e != nil {
		log.Fatalln("couldnt open db ", e)
	} else if e := createTestData(db); e != nil {
		log.Fatal("couldnt create test data ", e)
	} else if e := CreateAtlas(db); e != nil {
		log.Fatal("couldnt create atlas tables ", e)
	} else if e := listOfAspects(os.Stdout, db); e != nil {
		log.Fatal("couldnt process aspects ", e)
	}

	// Output:
	// <h1>Aspects</h1>
	//   <a href="#flightiness">Flightiness</a>.
	// <h2 id="flightiness">Flightiness</h2>
	// The flight worthiness of vehicles, an example of an aspect with several traits.
	// <h3>Kinds</h3>
	//   <a href="/atlas/kinds#vehicles">Vehicles</a>.
	// <h3>Traits</h3>
	// <dl>
	//   <dt>Flight Worthy</dt>
	//   <dt>Flightless</dt>
	//   <dt>Glide Worthy</dt>
	//    <dd>Better at landing than taking off.</dd>
	// </dl>
}
