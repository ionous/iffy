package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/ionous/iffy/tables"
)

func ExampleRelData() {
	templates.ExecuteTemplate(os.Stdout, "relList", []Relation{
		{"containing", "containers", tables.ONE_TO_MANY, "things", "Containers contain stuff."},
		{"driving", "people", tables.ONE_TO_ONE, "cars", "No backseat drivers please."},
	})

	// Output:
	// <h1>Relations</h1>
	// <dl>
	//   <dt><a href="/atlas/relations/containing">Containing</a></dt>
	//    <dd>Relates <a href="/atlas/kinds#containers">Containers</a> to many <a href="/atlas/kinds#things">Things</a>.
	//  Containers contain stuff.</dd>
	//   <dt><a href="/atlas/relations/driving">Driving</a></dt>
	//    <dd>Relates <a href="/atlas/kinds#people">People</a> to <a href="/atlas/kinds#cars">Cars</a>.
	//  No backseat drivers please.</dd>
	// </dl>
}

func ExampleRelDB() {
	const memory = "file:ExampleRelDB.db?cache=shared&mode=memory"
	if db, e := sql.Open(tables.DefaultDriver, memory); e != nil {
		log.Fatalln("couldnt open db ", e)
	} else if e := createTestData(db); e != nil {
		log.Fatal("couldnt create test data ", e)
	} else if e := CreateAtlas(db); e != nil {
		log.Fatal("couldnt create atlas tables ", e)
	} else if e := listOfRelations(os.Stdout, db); e != nil {
		log.Fatal("couldnt process relations ", e)
	}

	// Output:
	// <h1>Relations</h1>
	// <dl>
	//   <dt><a href="/atlas/relations/containing">Containing</a></dt>
	//    <dd>Relates <a href="/atlas/kinds#vehicles">Vehicles</a> to many <a href="/atlas/kinds#people">People</a>.
	//  The outside of insides.</dd>
	// </dl>
}
