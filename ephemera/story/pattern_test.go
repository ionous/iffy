package story

import (
	"strings"
	"testing"

	"github.com/ionous/iffy/tables"
)

// import an object type description
func TestObjectType(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if n, e := imp_object_type(k, _object_type); e != nil {
		t.Fatal(e)
	} else if n.String() != "animals" {
		t.Fatal(n)
	}
}

// import a variable type description
func TestVariableTypePrimitive(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if varType, e := imp_variable_type(k, _variable_type); e != nil {
		t.Fatal(e)
	} else if varType.String() != "text_eval" {
		t.Fatal(varType)
	}
}

// import a variable declaration
func TestVariableDeclObject(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if varName, typeName, e := imp_variable_decl(k, _variable_decl); e != nil {
		t.Fatal(e)
	} else if varName.String() != "pet" {
		t.Fatal(varName)
	} else if typeName.String() != "animals" {
		t.Fatal(typeName)
	}
}

func TestPatternVariablesDecl(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if e := imp_pattern_variables_decl(k, _pattern_variables_decl); e != nil {
		t.Fatal(e)
	} else {
		var buf strings.Builder
		tables.WriteCsv(db, &buf, "select name, category from eph_named", 2)
		tables.WriteCsv(db, &buf, "select * from eph_pattern", 4)
		if have, want := buf.String(), lines(
			"corral,pattern_name",  // 1
			"pet,variable_name",    // 2
			"animals,plural_kinds", // 3
			"1,2,3,1",              // NewPatternDecl
		); have != want {
			t.Fatal("mismatch", have)
		} else {
			t.Log("ok")
		}
	}
}

func TestPrimitiveType(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if typ, e := imp_primitive_type(k, _primitive_type); e != nil {
		t.Fatal(e)
	} else if typ.String() != "bool_eval" {
		t.Fatal(typ)
	}
}

func TestPatternType_Activity(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if typ, e := imp_pattern_type(k, _pattern_type_activity); e != nil {
		t.Fatal(e)
	} else if typ.String() != "execute" {
		t.Fatal(typ)
	}
}

func TestPatternType_Primitive(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if typ, e := imp_pattern_type(k, _pattern_type_primitive); e != nil {
		t.Fatal(e)
	} else if typ.String() != "bool_eval" {
		t.Fatal(typ)
	}
}

func TestPatternName(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if n, e := imp_pattern_name(k, _pattern_name); e != nil {
		t.Fatal(e)
	} else if n.String() != "corral" {
		t.Fatal(n)
	}
}

func TestPatternDecl(t *testing.T) {
	k, db := newTestImporter(t)
	defer db.Close()
	if e := imp_pattern_decl(k, _pattern_decl); e != nil {
		t.Fatal(e)
	} else {
		var buf strings.Builder
		tables.WriteCsv(db, &buf, "select name, category from eph_named", 2)
		tables.WriteCsv(db, &buf, "select * from eph_pattern", 4)
		if have, want := buf.String(), lines(
			"corral,pattern_name", // 1
			"execute,type",        // 2
			"1,1,2,1",             // NewPatternDecl
		); have != want {
			t.Fatal("mismatch", have)
		} else {
			t.Log("ok")
		}
	}
}

var _pattern_decl = map[string]interface{}{
	"type": "pattern_decl",
	"value": map[string]interface{}{
		"$NAME": _pattern_name,
		"$TYPE": _pattern_type_activity,
	},
}

var _pattern_name = map[string]interface{}{
	"type":  "pattern_name",
	"value": "corral",
}

var _pattern_type_activity = map[string]interface{}{
	"type": "pattern_type",
	"value": map[string]interface{}{
		"$ACTIVITY": map[string]interface{}{
			"type":  "patterned_activity",
			"value": "$ACTIVITY",
		},
	},
}

var _pattern_type_primitive = map[string]interface{}{
	"type": "pattern_type",
	"value": map[string]interface{}{
		"$VALUE": map[string]interface{}{
			"type": "variable_type",
			"value": map[string]interface{}{
				"$PRIMITIVE": _primitive_type,
			},
		},
	},
}

var _primitive_type = map[string]interface{}{
	"type":  "primitive_type",
	"value": "$BOOL",
}

var _pattern_variables_decl = map[string]interface{}{
	"type": "pattern_variables_decl",
	"value": map[string]interface{}{
		"$PATTERN_NAME": map[string]interface{}{
			"type":  "pattern_name",
			"value": "corral",
		},
		"$VARIABLE_DECL": []interface{}{
			_variable_decl,
		},
	},
}

var _variable_decl = map[string]interface{}{
	"type": "variable_decl",
	"value": map[string]interface{}{
		"$TYPE": map[string]interface{}{
			"type": "variable_type",
			"value": map[string]interface{}{
				"$OBJECT": _object_type,
			},
		},
		"$NAME": map[string]interface{}{
			"type":  "variable_name",
			"value": "pet",
		},
	},
}

var _variable_type = map[string]interface{}{
	"type": "variable_type",
	"value": map[string]interface{}{
		"$PRIMITIVE": map[string]interface{}{
			"type":  "primitive_type",
			"value": "$TEXT",
		},
	},
}

var _object_type = map[string]interface{}{
	"type": "object_type",
	"value": map[string]interface{}{
		"$AN": map[string]interface{}{
			"type":  "an",
			"value": "$AN",
		},
		"$KINDS": map[string]interface{}{
			"type":  "plural_kinds",
			"value": "animals",
		},
	},
}
