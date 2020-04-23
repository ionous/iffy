package internal

import (
	"database/sql"
	"os/user"
	"path"
	"strings"
	"testing"

	"github.com/ionous/iffy/ephemera"
	"github.com/ionous/iffy/tables"
	"github.com/kr/pretty"
)

// import an object type description
func TestXObjectType(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if n, e := imp_object_type(k, objectTypeData); e != nil {
		t.Fatal(e)
	} else if n.String() != "animals" {
		t.Fatal(n)
	}
}

// import a variable type description
func TestXVariableTypePrimitive(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if e := imp_variable_type(k, variableTypeData,
		func(prim string) {
			if prim != tables.EVAL_TEXT {
				t.Fatal(prim)
			}
		},
		nil); e != nil {
		t.Fatal(e)
	}
}

// import a variable declaration
func TestXVariableDeclObject(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if e := imp_variable_decl(k, variableDeclData, nil,
		func(varName ephemera.Named, typeName ephemera.Named) {
			if varName.String() != "pet" {
				t.Fatal(varName)
			} else if typeName.String() != "animals" {
				t.Fatal(typeName)
			}
		}); e != nil {
		t.Fatal(e)
	}
}

func TestXPatternVariablesDecl(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if e := imp_pattern_variables_decl(k, imp_pattern_variables_declData); e != nil {
		t.Fatal(e)
	} else {
		var b strings.Builder
		tables.WriteCsv(db, &b, "select name, category from eph_named", 2)
		tables.WriteCsv(db, &b, "select * from eph_pattern_eval", 3)
		tables.WriteCsv(db, &b, "select * from eph_pattern_kind", 3)
		if diff := pretty.Diff(b.String(), lines(
			"corral,pattern_name",
			"pet,variable_name",
			"animals,plural_kinds",
			"1,2,3",
		)); len(diff) > 0 {
			t.Fatal("mismatch", diff)
		} else {
			t.Log("ok")
		}
	}
}

func lines(s ...string) string {
	return strings.Join(s, "\n") + "\n"
}

func newTestImporter(t *testing.T) (ret *Importer, retDB *sql.DB) {
	const path = "file:test.db?cache=shared&mode=memory"
	// if path, e := getPath(t.Name() + ".db"); e != nil {
	// 	t.Fatal(e)
	// } else
	if db, e := sql.Open("sqlite3", path); e != nil {
		t.Fatal("db open", e)
	} else {
		if e := tables.CreateEphemera(db); e != nil {
			t.Fatal("create ephemera", e)
		} else {
			ret, retDB = NewImporter(t.Name(), db), db
		}
	}
	return
}

func getPath(file string) (ret string, err error) {
	if user, e := user.Current(); e != nil {
		err = e
	} else {
		ret = path.Join(user.HomeDir, file)
	}
	return
}

var imp_pattern_variables_declData = map[string]interface{}{
	"id":   "id-1719a47c939-7",
	"type": "pattern_variables_decl",
	"value": map[string]interface{}{
		"$PATTERN_NAME": map[string]interface{}{
			"id":    "id-1719a47c939-3",
			"type":  "pattern_name",
			"value": "corral",
		},
		"$VARIABLE_DECL": []interface{}{
			variableDeclData,
		},
	},
}

var variableDeclData = map[string]interface{}{
	"id":   "id-1719a47c939-11",
	"type": "variable_decl",
	"value": map[string]interface{}{
		"$TYPE": map[string]interface{}{
			"id":   "id-1719a47c939-9",
			"type": "variable_type",
			"value": map[string]interface{}{
				"$OBJECT": objectTypeData,
			},
		},
		"$NAME": map[string]interface{}{
			"id":    "id-1719a47c939-10",
			"type":  "variable_name",
			"value": "pet",
		},
	},
}
var variableTypeData = map[string]interface{}{
	"id":   "id-1719a47c939-4",
	"type": "variable_type",
	"value": map[string]interface{}{
		"$PRIMITIVE": map[string]interface{}{
			"id":    "id-1719a47c939-8",
			"type":  "primitive_type",
			"value": "$TEXT",
		},
	},
}
var objectTypeData = map[string]interface{}{
	"id":   "id-1719a47c939-14",
	"type": "object_type",
	"value": map[string]interface{}{
		"$AN": map[string]interface{}{
			"id":    "id-1719a47c939-12",
			"type":  "an",
			"value": "$AN",
		},
		"$KINDS": map[string]interface{}{
			"id":    "id-1719a47c939-13",
			"type":  "plural_kinds",
			"value": "animals",
		},
	},
}
