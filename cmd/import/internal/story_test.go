package internal

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/ionous/iffy/ephemera/debug"
	"github.com/ionous/iffy/ephemera/reader"
	"github.com/ionous/iffy/tables"
)

func TestImportStory(t *testing.T) {
	const memory = "file:test.db?cache=shared&mode=memory"
	if db, e := sql.Open("sqlite3", memory); e != nil {
		t.Fatal("db open", e)
	} else {
		defer db.Close()
		if e := tables.CreateEphemera(db); e != nil {
			t.Fatal("create ephemera", e)
		}
		var in reader.Map
		if e := json.Unmarshal([]byte(debug.Blob), &in); e != nil {
			t.Fatal("read json", e)
		} else if ok := in.Has(itemType, "story"); !ok {
			t.Fatal("read story")
		} else {
			k := NewImporter(t.Name(), db, generators)
			//
			if e := k.parseItem(in); e != nil {
				t.Fatal(e)
			}
			// if b, e := json.MarshalIndent(genq.Tables, "", "  "); e != nil {
			// 	t.Fatal(e)
			// } else {
			// 	pretty.Println("tables", string(b))
			// }
		}
	}
}

// tables {
//   "eph_named": [
//     {
//       "category": "kind",
//       "idSource": 1,
//       "name": "things",
//       "offset": "noun"
//     },
//     {
//       "category": "aspect",
//       "idSource": 1,
//       "name": "nounType",
//       "offset": "noun"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "common",
//       "offset": "noun"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "proper",
//       "offset": "noun"
//     },
//     {
//       "category": "noun",
//       "idSource": 1,
//       "name": "cabin",
//       "offset": "id7"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "common",
//       "offset": "id8"
//     },
//     {
//       "category": "kind",
//       "idSource": 1,
//       "name": "room",
//       "offset": "id11"
//     },
//     {
//       "category": "kind",
//       "idSource": 1,
//       "name": "things",
//       "offset": "summary"
//     },
//     {
//       "category": "field",
//       "idSource": 1,
//       "name": "appearance",
//       "offset": "summary"
//     },
//     {
//       "category": "field",
//       "idSource": 1,
//       "name": "appearance",
//       "offset": "id12"
//     },
//     {
//       "category": "expr",
//       "idSource": 1,
//       "name": "The front of the small cabin is entirely occupied with navigational instruments, a radar display, and radios for calling back to shore. Along each side runs a bench with faded blue vinyl cushions[if the compartment is closed], which can be lifted to reveal the storage space underneath[otherwise], one of which is currently lifted to allow access to the storage compartment within[end if]. A glass case against the wall contains several fishing rods.\n\nScratched windows offer a view of the surrounding bay, and there is a door south to the deck. A sign taped to one wall announces the menu of tours offered by the Yakutat Charter Boat Company.",
//       "offset": "id13"
//     },
//     {
//       "category": "noun",
//       "idSource": 1,
//       "name": "cabin",
//       "offset": "id20"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "common",
//       "offset": "id18"
//     },
//     {
//       "category": "verb",
//       "idSource": 1,
//       "name": "contains",
//       "offset": "id23"
//     },
//     {
//       "category": "noun",
//       "idSource": 1,
//       "name": "glass case",
//       "offset": "id27"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "common",
//       "offset": "id25"
//     },
//     {
//       "category": "verb",
//       "idSource": 1,
//       "name": "In",
//       "offset": "id-16f156a8166-0"
//     },
//     {
//       "category": "noun",
//       "idSource": 1,
//       "name": "glass case",
//       "offset": "id-16f156a8166-6"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "common",
//       "offset": "id-16f156a8166-7"
//     },
//     {
//       "category": "noun",
//       "idSource": 1,
//       "name": "collection of fishing rods",
//       "offset": "id-16f156a8166-9"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "common",
//       "offset": "id-16f156a8166-10"
//     },
//     {
//       "category": "noun",
//       "idSource": 1,
//       "name": "case",
//       "offset": "id-16f156a8166-17"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "common",
//       "offset": "id-16f156a8166-18"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "closed",
//       "offset": "id-16f156a8166-20"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "transparent",
//       "offset": "id-16f156a8166-22"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "lockable",
//       "offset": "id-16f156a8166-23"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "locked",
//       "offset": "id-16f156a8166-24"
//     },
//     {
//       "category": "noun",
//       "idSource": 1,
//       "name": "case",
//       "offset": "id-16f156a8166-31"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "common",
//       "offset": "id-16f156a8166-32"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "scenery",
//       "offset": "id-16f156a8166-34"
//     },
//     {
//       "category": "noun",
//       "idSource": 1,
//       "name": "small silver key",
//       "offset": "id-16f156a8166-42"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "common",
//       "offset": "id-16f156a8166-43"
//     },
//     {
//       "category": "verb",
//       "idSource": 1,
//       "name": "unlocks",
//       "offset": "id-16f156a8166-53"
//     },
//     {
//       "category": "noun",
//       "idSource": 1,
//       "name": "case",
//       "offset": "id-16f156a8166-57"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "common",
//       "offset": "id-16f156a8166-58"
//     },
//     {
//       "category": "noun",
//       "idSource": 1,
//       "name": "bench",
//       "offset": "id-16f156a8166-70"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "common",
//       "offset": "id-16f156a8166-71"
//     },
//     {
//       "category": "verb",
//       "idSource": 1,
//       "name": "in",
//       "offset": "id-16f156a8166-72"
//     },
//     {
//       "category": "noun",
//       "idSource": 1,
//       "name": "cabin",
//       "offset": "id-16f156a8166-76"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "common",
//       "offset": "id-16f156a8166-77"
//     },
//     {
//       "category": "verb",
//       "idSource": 1,
//       "name": "On",
//       "offset": "id-16f165d636d-1"
//     },
//     {
//       "category": "noun",
//       "idSource": 1,
//       "name": "bench",
//       "offset": "id-16f165d636d-7"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "common",
//       "offset": "id-16f165d636d-8"
//     },
//     {
//       "category": "noun",
//       "idSource": 1,
//       "name": "blue vinyl cushions",
//       "offset": "id-16f165d636d-10"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "common",
//       "offset": "id-16f165d636d-11"
//     },
//     {
//       "category": "field",
//       "idSource": 1,
//       "name": "indefinite article",
//       "offset": "id-16f165d636d-11"
//     },
//     {
//       "category": "field",
//       "idSource": 1,
//       "name": "indefinite article",
//       "offset": "common_noun"
//     },
//     {
//       "category": "kind",
//       "idSource": 1,
//       "name": "things",
//       "offset": "common_noun"
//     },
//     {
//       "category": "noun",
//       "idSource": 1,
//       "name": "bench",
//       "offset": "id-16f165d636d-18"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "common",
//       "offset": "id-16f165d636d-19"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "enterable",
//       "offset": "id-16f165d636d-21"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "scenery",
//       "offset": "id-16f165d636d-23"
//     },
//     {
//       "category": "noun",
//       "idSource": 1,
//       "name": "cushions",
//       "offset": "id-16f165d636d-35"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "common",
//       "offset": "id-16f165d636d-36"
//     },
//     {
//       "category": "trait",
//       "idSource": 1,
//       "name": "scenery",
//       "offset": "id-16f165d636d-38"
//     }
//   ],
//   "eph_noun": [
//     {
//       "idNamedKind": "room",
//       "idNamedNoun": "cabin"
//     }
//   ],
//   "eph_field": [
//     {
//       "idNamedField": "nounType",
//       "idNamedKind": "things",
//       "primType": "aspect"
//     },
//     {
//       "idNamedField": "appearance",
//       "idNamedKind": "things",
//       "primType": "expr"
//     },
//     {
//       "idNamedField": "indefinite article",
//       "idNamedKind": "things",
//       "primType": "text"
//     }
//   ],
//   "eph_relative": [
//     {
//       "idNamedDependent": "glass case",
//       "idNamedHead": "cabin",
//       "idNamedStem": "contains"
//     },
//     {
//       "idNamedDependent": "collection of fishing rods",
//       "idNamedHead": "glass case",
//       "idNamedStem": "In"
//     },
//     {
//       "idNamedDependent": "case",
//       "idNamedHead": "small silver key",
//       "idNamedStem": "unlocks"
//     },
//     {
//       "idNamedDependent": "cabin",
//       "idNamedHead": "bench",
//       "idNamedStem": "in"
//     },
//     {
//       "idNamedDependent": "blue vinyl cushions",
//       "idNamedHead": "bench",
//       "idNamedStem": "On"
//     }
//   ],
//   "eph_source": [
//     {
//       "src": "blob"
//     }
//   ],
//   "eph_trait": [
//     {
//       "idNamedAspect": "nounType",
//       "idNamedTrait": "common",
//       "rank": 0
//     },
//     {
//       "idNamedAspect": "nounType",
//       "idNamedTrait": "proper",
//       "rank": 1
//     }
//   ],
//   "eph_value": [
//     {
//       "idNamedNoun": "cabin",
//       "idNamedProp": "common",
//       "value": true
//     },
//     {
//       "idNamedNoun": "cabin",
//       "idNamedProp": "appearance",
//       "value": "The front of the small cabin is entirely occupied with navigational instruments, a radar display, and radios for calling back to shore. Along each side runs a bench with faded blue vinyl cushions[if the compartment is closed], which can be lifted to reveal the storage space underneath[otherwise], one of which is currently lifted to allow access to the storage compartment within[end if]. A glass case against the wall contains several fishing rods.\n\nScratched windows offer a view of the surrounding bay, and there is a door south to the deck. A sign taped to one wall announces the menu of tours offered by the Yakutat Charter Boat Company."
//     },
//     {
//       "idNamedNoun": "cabin",
//       "idNamedProp": "common",
//       "value": true
//     },
//     {
//       "idNamedNoun": "glass case",
//       "idNamedProp": "common",
//       "value": true
//     },
//     {
//       "idNamedNoun": "glass case",
//       "idNamedProp": "common",
//       "value": true
//     },
//     {
//       "idNamedNoun": "collection of fishing rods",
//       "idNamedProp": "common",
//       "value": true
//     },
//     {
//       "idNamedNoun": "case",
//       "idNamedProp": "common",
//       "value": true
//     },
//     {
//       "idNamedNoun": "case",
//       "idNamedProp": "closed",
//       "value": true
//     },
//     {
//       "idNamedNoun": "case",
//       "idNamedProp": "transparent",
//       "value": true
//     },
//     {
//       "idNamedNoun": "case",
//       "idNamedProp": "lockable",
//       "value": true
//     },
//     {
//       "idNamedNoun": "case",
//       "idNamedProp": "locked",
//       "value": true
//     },
//     {
//       "idNamedNoun": "case",
//       "idNamedProp": "common",
//       "value": true
//     },
//     {
//       "idNamedNoun": "case",
//       "idNamedProp": "scenery",
//       "value": true
//     },
//     {
//       "idNamedNoun": "small silver key",
//       "idNamedProp": "common",
//       "value": true
//     },
//     {
//       "idNamedNoun": "case",
//       "idNamedProp": "common",
//       "value": true
//     },
//     {
//       "idNamedNoun": "bench",
//       "idNamedProp": "common",
//       "value": true
//     },
//     {
//       "idNamedNoun": "cabin",
//       "idNamedProp": "common",
//       "value": true
//     },
//     {
//       "idNamedNoun": "bench",
//       "idNamedProp": "common",
//       "value": true
//     },
//     {
//       "idNamedNoun": "blue vinyl cushions",
//       "idNamedProp": "common",
//       "value": true
//     },
//     {
//       "idNamedNoun": "blue vinyl cushions",
//       "idNamedProp": "indefinite article",
//       "value": "some"
//     },
//     {
//       "idNamedNoun": "bench",
//       "idNamedProp": "common",
//       "value": true
//     },
//     {
//       "idNamedNoun": "bench",
//       "idNamedProp": "enterable",
//       "value": true
//     },
//     {
//       "idNamedNoun": "bench",
//       "idNamedProp": "scenery",
//       "value": true
//     },
//     {
//       "idNamedNoun": "cushions",
//       "idNamedProp": "common",
//       "value": true
//     },
//     {
//       "idNamedNoun": "cushions",
//       "idNamedProp": "scenery",
//       "value": true
//     }
//   ]
// }

// func getUserFile() (ret string, err error) {
//   if user, e := user.Current(); e != nil {
//     err = e
//   } else {
//     ret = path.Join(user.HomeDir, "iffyTest.db")
//   }
//   return
// }
