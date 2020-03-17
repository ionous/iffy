/* generated using github.com/ionous/iffy/cmd/spec/spec.go */
const spec = [
  {
    "desc": "Execute: Run a series of statements.",
    "name": "execute",
    "uses": "slot"
  },
  {
    "desc": "Booleans: Statements which return true/false values.",
    "name": "bool_eval",
    "uses": "slot"
  },
  {
    "desc": "Numbers: Statements which return a number.",
    "name": "number_eval",
    "uses": "slot"
  },
  {
    "desc": "Texts: Statements which return text.",
    "name": "text_eval",
    "uses": "slot"
  },
  {
    "desc": "Objects: Statements which return an existing object.",
    "name": "object_eval",
    "uses": "slot"
  },
  {
    "desc": "Number List: Statements which return a list of numbers.",
    "name": "num_list_eval",
    "uses": "slot"
  },
  {
    "desc": "Text Lists: Statements which return a list of text.",
    "name": "text_list_eval",
    "uses": "slot"
  },
  {
    "desc": "Object Lists: Statements which return a list of existing objects.",
    "name": "obj_list_eval",
    "uses": "slot"
  },
  {
    "desc": "Comparison Types: Helper used when comparing two numbers, objects, pieces of text, etc.",
    "name": "compare_to",
    "uses": "slot"
  },
  {
    "desc": "Add Numbers: Add two numbers.",
    "group": [
      "math"
    ],
    "name": "sum_of",
    "uses": "run",
    "with": {
      "params": {
        "$A": {
          "label": "a",
          "type": "number_eval"
        },
        "$B": {
          "label": "b",
          "type": "number_eval"
        }
      },
      "slots": [
        "number_eval"
      ],
      "tokens": [
        "(",
        "$A",
        "+",
        "$B",
        ")"
      ]
    }
  },
  {
    "desc": "All True: returns true if all of the evaluations are true.",
    "group": [
      "logic"
    ],
    "name": "all_true",
    "uses": "run",
    "with": {
      "params": {
        "$TEST": {
          "label": "test",
          "repeats": true,
          "type": "bool_eval"
        }
      },
      "slots": [
        "bool_eval"
      ],
      "tokens": [
        "all true",
        "$TEST"
      ]
    }
  },
  {
    "desc": "Any True: returns true if any of the evaluations are true.",
    "group": [
      "logic"
    ],
    "name": "any_true",
    "uses": "run",
    "with": {
      "params": {
        "$TEST": {
          "label": "test",
          "repeats": true,
          "type": "bool_eval"
        }
      },
      "slots": [
        "bool_eval"
      ],
      "tokens": [
        "any true",
        "$TEST"
      ]
    }
  },
  {
    "desc": "Bool Value: specifies an explicit true/false value.",
    "group": [
      "literals"
    ],
    "name": "bool_value",
    "uses": "run",
    "spec": "{bool:bool_eval}",
    "with": {
      "params": {
        "$BOOL": {
          "label": "bool",
          "type": "bool"
        }
      },
      "slots": [
        "bool_eval"
      ],
      "tokens": [
        "bool value",
        "$BOOL"
      ]
    }
  },
  {
    "desc": "Branch: execute a single block of statements based on a boolean test.",
    "group": [
      "exec"
    ],
    "name": "choose",
    "uses": "run",
    "spec": "if {choose%if:bool_eval} then: {true*execute|ghost} else: {false*execute|ghost}",
    "with": {
      "params": {
        "$FALSE": {
          "label": "false",
          "repeats": true,
          "type": "execute"
        },
        "$IF": {
          "label": "if",
          "type": "bool_eval"
        },
        "$TRUE": {
          "label": "true",
          "repeats": true,
          "type": "execute"
        }
      },
      "slots": [
        "execute"
      ],
      "tokens": [
        "choose",
        "$IF",
        "$TRUE",
        "$FALSE"
      ]
    }
  },
  {
    "desc": "Choose Number: Pick one of two numbers based on a boolean test.",
    "group": [
      "math"
    ],
    "name": "choose_num",
    "uses": "run",
    "with": {
      "params": {
        "$FALSE": {
          "label": "false",
          "type": "number_eval"
        },
        "$IF": {
          "label": "if",
          "type": "bool_eval"
        },
        "$TRUE": {
          "label": "true",
          "type": "number_eval"
        }
      },
      "slots": [
        "number_eval"
      ],
      "tokens": [
        "choose num",
        "$IF",
        "$TRUE",
        "$FALSE"
      ]
    }
  },
  {
    "desc": "Choose Object: Pick one of two objects based on a boolean test.",
    "group": [
      "objects"
    ],
    "name": "choose_obj",
    "uses": "run",
    "with": {
      "params": {
        "$FALSE": {
          "label": "false",
          "type": "object_eval"
        },
        "$IF": {
          "label": "if",
          "type": "bool_eval"
        },
        "$TRUE": {
          "label": "true",
          "type": "object_eval"
        }
      },
      "slots": [
        "object_eval"
      ],
      "tokens": [
        "choose obj",
        "$IF",
        "$TRUE",
        "$FALSE"
      ]
    }
  },
  {
    "desc": "Choose Text: Pick one of two strings based on a boolean test.",
    "group": [
      "format"
    ],
    "name": "choose_text",
    "uses": "run",
    "with": {
      "params": {
        "$FALSE": {
          "label": "false",
          "type": "text_eval"
        },
        "$IF": {
          "label": "if",
          "type": "bool_eval"
        },
        "$TRUE": {
          "label": "true",
          "type": "text_eval"
        }
      },
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "choose text",
        "$IF",
        "$TRUE",
        "$FALSE"
      ]
    }
  },
  {
    "desc": "Compare Numbers: True if eq,ne,gt,lt,ge,le two numbers.",
    "group": [
      "logic"
    ],
    "name": "compare_num",
    "uses": "run",
    "with": {
      "params": {
        "$A": {
          "label": "a",
          "type": "number_eval"
        },
        "$B": {
          "label": "b",
          "type": "number_eval"
        },
        "$IS": {
          "label": "is",
          "type": "compare_to"
        }
      },
      "slots": [
        "bool_eval"
      ],
      "tokens": [
        "compare num",
        "$A",
        "$IS",
        "$B"
      ]
    }
  },
  {
    "desc": "Compare Objects",
    "name": "compare_obj",
    "uses": "run",
    "with": {
      "params": {
        "$A": {
          "label": "a",
          "type": "object_eval"
        },
        "$B": {
          "label": "b",
          "type": "object_eval"
        },
        "$IS": {
          "label": "is",
          "type": "compare_to"
        }
      },
      "slots": [
        "bool_eval"
      ],
      "tokens": [
        "compare obj",
        "$A",
        "$IS",
        "$B"
      ]
    }
  },
  {
    "desc": "Compare Text: True if eq,ne,gt,lt,ge,le two strings ( lexical. )",
    "group": [
      "logic"
    ],
    "name": "compare_text",
    "uses": "run",
    "with": {
      "params": {
        "$A": {
          "label": "a",
          "type": "text_eval"
        },
        "$B": {
          "label": "b",
          "type": "text_eval"
        },
        "$IS": {
          "label": "is",
          "type": "compare_to"
        }
      },
      "slots": [
        "bool_eval"
      ],
      "tokens": [
        "compare text",
        "$A",
        "$IS",
        "$B"
      ]
    }
  },
  {
    "desc": "Cycle Text: When called multiple times, returns each of its inputs in turn.",
    "group": [
      "cycle"
    ],
    "name": "cycle_text",
    "uses": "run",
    "with": {
      "params": {
        "$ID": {
          "label": "id",
          "type": "text"
        },
        "$VALUES": {
          "label": "values",
          "repeats": true,
          "type": "text_eval"
        }
      },
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "cycle text",
        "$ID",
        "$VALUES"
      ]
    }
  },
  {
    "desc": "Divide Numbers: Divide one number by another.",
    "group": [
      "math"
    ],
    "name": "quotient_of",
    "uses": "run",
    "with": {
      "params": {
        "$A": {
          "label": "a",
          "type": "number_eval"
        },
        "$B": {
          "label": "b",
          "type": "number_eval"
        }
      },
      "slots": [
        "number_eval"
      ],
      "tokens": [
        "(",
        "$A",
        "/",
        "$B",
        ")"
      ]
    }
  },
  {
    "desc": "Do Nothing: Statement which does nothing.",
    "group": [
      "exec"
    ],
    "name": "do_nothing",
    "uses": "run",
    "with": {
      "params": {},
      "slots": [
        "execute"
      ],
      "tokens": [
        "do nothing"
      ]
    }
  },
  {
    "desc": "Filter Object List: A list of objects which pass the evaluation.",
    "group": [
      "objects"
    ],
    "name": "filter",
    "uses": "run",
    "with": {
      "params": {
        "$ACCEPT": {
          "label": "accept",
          "type": "bool_eval"
        },
        "$LIST": {
          "label": "list",
          "type": "obj_list_eval"
        }
      },
      "slots": [
        "obj_list_eval"
      ],
      "tokens": [
        "filter",
        "$LIST",
        "$ACCEPT"
      ]
    }
  },
  {
    "desc": "For Each Number: Loops over the passed list of numbers, or runs the 'else' statement if empty.",
    "group": [
      "exec"
    ],
    "name": "for_each_num",
    "uses": "run",
    "with": {
      "params": {
        "$ELSE": {
          "label": "else",
          "repeats": true,
          "type": "execute"
        },
        "$GO": {
          "label": "go",
          "repeats": true,
          "type": "execute"
        },
        "$IN": {
          "label": "in",
          "type": "num_list_eval"
        }
      },
      "slots": [
        "execute"
      ],
      "tokens": [
        "for each num",
        "$IN",
        "$GO",
        "$ELSE"
      ]
    }
  },
  {
    "desc": "For Each Object: Loops over the passed list of objects, or runs the 'else' statement if empty.",
    "group": [
      "exec"
    ],
    "name": "for_each_obj",
    "uses": "run",
    "with": {
      "params": {
        "$ELSE": {
          "label": "else",
          "repeats": true,
          "type": "execute"
        },
        "$GO": {
          "label": "go",
          "repeats": true,
          "type": "execute"
        },
        "$IN": {
          "label": "in",
          "type": "obj_list_eval"
        }
      },
      "slots": [
        "execute"
      ],
      "tokens": [
        "for each obj",
        "$IN",
        "$GO",
        "$ELSE"
      ]
    }
  },
  {
    "desc": "For Each Text: Loops over the passed list of text, or runs the 'else' statement if empty.",
    "group": [
      "exec"
    ],
    "name": "for_each_text",
    "uses": "run",
    "with": {
      "params": {
        "$ELSE": {
          "label": "else",
          "repeats": true,
          "type": "execute"
        },
        "$GO": {
          "label": "go",
          "repeats": true,
          "type": "execute"
        },
        "$IN": {
          "label": "in",
          "type": "text_list_eval"
        }
      },
      "slots": [
        "execute"
      ],
      "tokens": [
        "for each text",
        "$IN",
        "$GO",
        "$ELSE"
      ]
    }
  },
  {
    "desc": "Get Property: Return the value of an object's property.",
    "group": [
      "objects"
    ],
    "name": "get",
    "uses": "run",
    "with": {
      "params": {
        "$OBJ": {
          "label": "obj",
          "type": "object_eval"
        },
        "$PROP": {
          "label": "prop",
          "type": "text"
        }
      },
      "slots": [
        "bool_eval",
        "number_eval",
        "text_eval",
        "object_eval",
        "num_list_eval",
        "text_list_eval",
        "obj_list_eval"
      ],
      "tokens": [
        "Get",
        "$PROP",
        "of",
        "$OBJ"
      ]
    }
  },
  {
    "desc": "Includes Text: True if text contains text.",
    "group": [
      "strings"
    ],
    "name": "includes",
    "uses": "run",
    "with": {
      "params": {
        "$PART": {
          "label": "part",
          "type": "text_eval"
        },
        "$TEXT": {
          "label": "text",
          "type": "text_eval"
        }
      },
      "slots": [
        "bool_eval"
      ],
      "tokens": [
        "includes",
        "$TEXT",
        "$PART"
      ]
    }
  },
  {
    "desc": "Is Kind Of: True if the object is compatible with the named kind.",
    "group": [
      "objects"
    ],
    "name": "is_class",
    "uses": "run",
    "with": {
      "params": {
        "$CLASS": {
          "label": "class",
          "type": "text"
        },
        "$OBJ": {
          "label": "obj",
          "type": "object_eval"
        }
      },
      "slots": [
        "bool_eval"
      ],
      "tokens": [
        "Is",
        "$OBJ",
        "a",
        "kind",
        "of",
        "$CLASS"
      ]
    }
  },
  {
    "desc": "Is Empty: True if the text is empty.",
    "group": [
      "strings"
    ],
    "name": "is_empty",
    "uses": "run",
    "with": {
      "params": {
        "$TEXT": {
          "label": "text",
          "type": "text_eval"
        }
      },
      "slots": [
        "bool_eval"
      ],
      "tokens": [
        "is empty",
        "$TEXT"
      ]
    }
  },
  {
    "desc": "Is Exact Kind: True if the object is exactly the named kind.",
    "group": [
      "objects"
    ],
    "name": "is_exact_class",
    "uses": "run",
    "with": {
      "params": {
        "$CLASS": {
          "label": "class",
          "type": "text"
        },
        "$OBJ": {
          "label": "obj",
          "type": "object_eval"
        }
      },
      "slots": [
        "bool_eval"
      ],
      "tokens": [
        "is exact class",
        "$OBJ",
        "$CLASS"
      ]
    }
  },
  {
    "desc": "Is Not: Returns the opposite value.",
    "group": [
      "logic"
    ],
    "name": "is_not",
    "uses": "run",
    "with": {
      "params": {
        "$BOOL_EVAL": {
          "label": "bool eval",
          "type": "bool_eval"
        }
      },
      "slots": [
        "bool_eval"
      ],
      "tokens": [
        "is not",
        "$BOOL_EVAL"
      ]
    }
  },
  {
    "desc": "Join Strings: Returns multiple pieces of text as a single new piece of text.",
    "group": [
      "strings"
    ],
    "name": "join",
    "uses": "run",
    "with": {
      "params": {
        "$TEXT": {
          "label": "text",
          "repeats": true,
          "type": "text_eval"
        }
      },
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "join",
        "$TEXT"
      ]
    }
  },
  {
    "desc": "Length of Object List: Number of objects.",
    "group": [
      "objects"
    ],
    "name": "len",
    "uses": "run",
    "with": {
      "params": {
        "$LIST": {
          "label": "list",
          "type": "obj_list_eval"
        }
      },
      "slots": [
        "number_eval"
      ],
      "tokens": [
        "len",
        "$LIST"
      ]
    }
  },
  {
    "desc": "List Up: Generates a list of objects.",
    "group": [
      "objects"
    ],
    "name": "list_up",
    "uses": "run",
    "with": {
      "params": {
        "$ALLOW_DUPLICATES": {
          "label": "allow duplicates",
          "type": "bool"
        },
        "$MAX_OBJECTS": {
          "label": "max objects",
          "type": "number"
        },
        "$NEXT": {
          "label": "next",
          "type": "object_eval"
        },
        "$SOURCE": {
          "label": "source",
          "type": "object_eval"
        }
      },
      "slots": [
        "obj_list_eval"
      ],
      "tokens": [
        "list up",
        "$SOURCE",
        "$NEXT",
        "$ALLOW_DUPLICATES",
        "$MAX_OBJECTS"
      ]
    }
  },
  {
    "desc": "Modulus Numbers: Divide one number by another, and return the remainder.",
    "group": [
      "math"
    ],
    "name": "remainder_of",
    "uses": "run",
    "with": {
      "params": {
        "$A": {
          "label": "a",
          "type": "number_eval"
        },
        "$B": {
          "label": "b",
          "type": "number_eval"
        }
      },
      "slots": [
        "number_eval"
      ],
      "tokens": [
        "(",
        "$A",
        "%",
        "$B",
        ")"
      ]
    }
  },
  {
    "desc": "Multiply Numbers: Multiply two numbers.",
    "group": [
      "math"
    ],
    "name": "product_of",
    "uses": "run",
    "with": {
      "params": {
        "$A": {
          "label": "a",
          "type": "number_eval"
        },
        "$B": {
          "label": "b",
          "type": "number_eval"
        }
      },
      "slots": [
        "number_eval"
      ],
      "tokens": [
        "(",
        "$A",
        "*",
        "$B",
        ")"
      ]
    }
  },
  {
    "desc": "Number Value: Specify a particular number.",
    "group": [
      "literals"
    ],
    "name": "num_value",
    "uses": "run",
    "with": {
      "params": {
        "$NUM": {
          "label": "num",
          "type": "number"
        }
      },
      "slots": [
        "number_eval"
      ],
      "tokens": [
        "num value",
        "$NUM"
      ]
    }
  },
  {
    "desc": "Number List: Specify a list of multiple numbers.",
    "group": [
      "literals"
    ],
    "name": "numbers",
    "uses": "run",
    "with": {
      "params": {
        "$VALUES": {
          "label": "values",
          "repeats": true,
          "type": "float64"
        }
      },
      "slots": [
        "num_list_eval"
      ],
      "tokens": [
        "numbers",
        "$VALUES"
      ]
    }
  },
  {
    "desc": "Named Object: Searches through the scope for a matching name.",
    "group": [
      "objects"
    ],
    "name": "object_name",
    "uses": "run",
    "with": {
      "params": {
        "$NAME": {
          "label": "name",
          "type": "text"
        }
      },
      "slots": [
        "object_eval"
      ],
      "tokens": [
        "object name",
        "$NAME"
      ]
    }
  },
  {
    "desc": "Object List: Searches through the scope for matching names.",
    "group": [
      "objects"
    ],
    "name": "object_names",
    "uses": "run",
    "with": {
      "params": {
        "$NAMES": {
          "label": "names",
          "repeats": true,
          "type": "string"
        }
      },
      "slots": [
        "obj_list_eval"
      ],
      "tokens": [
        "object names",
        "$NAMES"
      ]
    }
  },
  {
    "desc": "Pluralize: Creates plural text from the passed (presumably singular) text.",
    "group": [
      "format"
    ],
    "name": "pluralize",
    "uses": "run",
    "with": {
      "params": {
        "$TEXT": {
          "label": "text",
          "type": "text_eval"
        }
      },
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "pluralize",
        "$TEXT"
      ]
    }
  },
  {
    "desc": "Say List: Writes words separated with commas, ending with an 'and'.",
    "group": [
      "format"
    ],
    "name": "print_list",
    "uses": "run",
    "with": {
      "params": {
        "$BLOCK": {
          "label": "block",
          "repeats": true,
          "type": "execute"
        }
      },
      "slots": [
        "execute",
        "text_eval"
      ],
      "tokens": [
        "print list",
        "$BLOCK"
      ]
    }
  },
  {
    "desc": "Say Number: Writes a number using numerals, eg. '1'.",
    "group": [
      "format"
    ],
    "name": "print_num",
    "uses": "run",
    "with": {
      "params": {
        "$NUM": {
          "label": "num",
          "type": "number_eval"
        }
      },
      "slots": [
        "execute",
        "text_eval"
      ],
      "tokens": [
        "print num",
        "$NUM"
      ]
    }
  },
  {
    "desc": "Say Span: Writes text with spaces between words.",
    "group": [
      "format"
    ],
    "name": "print_span",
    "uses": "run",
    "with": {
      "params": {
        "$BLOCK": {
          "label": "block",
          "repeats": true,
          "type": "execute"
        }
      },
      "slots": [
        "execute",
        "text_eval"
      ],
      "tokens": [
        "print span",
        "$BLOCK"
      ]
    }
  },
  {
    "desc": "Range of Numbers: Generates a series of numbers.",
    "group": [
      "flow"
    ],
    "name": "range",
    "uses": "run",
    "with": {
      "params": {
        "$END": {
          "label": "end",
          "type": "number"
        },
        "$START": {
          "label": "start",
          "type": "number"
        },
        "$STEP": {
          "label": "step",
          "type": "number"
        }
      },
      "slots": [
        "num_list_eval"
      ],
      "tokens": [
        "range",
        "$START",
        "$END",
        "$STEP"
      ]
    }
  },
  {
    "desc": "Related List: Returns a stream of objects related to the requested object.",
    "group": [
      "objects"
    ],
    "name": "related_list",
    "uses": "run",
    "with": {
      "params": {
        "$OBJECT": {
          "label": "object",
          "type": "object_eval"
        },
        "$RELATION": {
          "label": "relation",
          "type": "text"
        }
      },
      "slots": [
        "obj_list_eval"
      ],
      "tokens": [
        "related list",
        "$RELATION",
        "$OBJECT"
      ]
    }
  },
  {
    "desc": "Is Relation Empty: Returns true if the requested object has no related objects.",
    "group": [
      "objects"
    ],
    "name": "relation_empty",
    "uses": "run",
    "with": {
      "params": {
        "$OBJECT": {
          "label": "object",
          "type": "object_eval"
        },
        "$RELATION": {
          "label": "relation",
          "type": "text"
        }
      },
      "slots": [
        "bool_eval"
      ],
      "tokens": [
        "relation empty",
        "$RELATION",
        "$OBJECT"
      ]
    }
  },
  {
    "desc": "Reverse Object List: returns the listed objects, last first.",
    "group": [
      "objects"
    ],
    "name": "reverse",
    "uses": "run",
    "with": {
      "params": {
        "$LIST": {
          "label": "list",
          "type": "obj_list_eval"
        }
      },
      "slots": [
        "obj_list_eval"
      ],
      "tokens": [
        "reverse",
        "$LIST"
      ]
    }
  },
  {
    "desc": "Say: writes a piece of text.",
    "group": [
      "format"
    ],
    "name": "say",
    "uses": "run",
    "with": {
      "params": {
        "$TEXT": {
          "label": "text",
          "type": "text_eval"
        }
      },
      "slots": [
        "execute"
      ],
      "tokens": [
        "say",
        "$TEXT"
      ]
    }
  },
  {
    "desc": "Set Bool: Sets the named property to the passed boolean value.",
    "group": [
      "objects"
    ],
    "name": "set_bool",
    "uses": "run",
    "with": {
      "params": {
        "$OBJ": {
          "label": "obj",
          "type": "object_eval"
        },
        "$PROP": {
          "label": "prop",
          "type": "text"
        },
        "$VAL": {
          "label": "val",
          "type": "bool_eval"
        }
      },
      "slots": [
        "execute"
      ],
      "tokens": [
        "set bool",
        "$OBJ",
        "$PROP",
        "$VAL"
      ]
    }
  },
  {
    "desc": "Set Number: sets the named property to the passed number.",
    "group": [
      "objects"
    ],
    "name": "set_num",
    "uses": "run",
    "with": {
      "params": {
        "$OBJ": {
          "label": "obj",
          "type": "object_eval"
        },
        "$PROP": {
          "label": "prop",
          "type": "text"
        },
        "$VAL": {
          "label": "val",
          "type": "number_eval"
        }
      },
      "slots": [
        "execute"
      ],
      "tokens": [
        "set num",
        "$OBJ",
        "$PROP",
        "$VAL"
      ]
    }
  },
  {
    "desc": "Set Object: Sets the named property to the passed object ( reference. )",
    "group": [
      "objects"
    ],
    "name": "set_obj",
    "uses": "run",
    "with": {
      "params": {
        "$OBJ": {
          "label": "obj",
          "type": "object_eval"
        },
        "$PROP": {
          "label": "prop",
          "type": "text"
        },
        "$VAL": {
          "label": "val",
          "type": "object_eval"
        }
      },
      "slots": [
        "execute"
      ],
      "tokens": [
        "set obj",
        "$OBJ",
        "$PROP",
        "$VAL"
      ]
    }
  },
  {
    "desc": "Set State: Sets the object to the passed state.",
    "group": [
      "objects"
    ],
    "name": "set_state",
    "uses": "run",
    "with": {
      "params": {
        "$REF": {
          "label": "ref",
          "type": "object_eval"
        },
        "$STATE": {
          "label": "state",
          "type": "text"
        }
      },
      "slots": [
        "execute"
      ],
      "tokens": [
        "set state",
        "$REF",
        "$STATE"
      ]
    }
  },
  {
    "desc": "Set Text: Sets the named property to the passed string.",
    "group": [
      "objects"
    ],
    "name": "set_text",
    "uses": "run",
    "with": {
      "params": {
        "$OBJ": {
          "label": "obj",
          "type": "object_eval"
        },
        "$PROP": {
          "label": "prop",
          "type": "text"
        },
        "$VAL": {
          "label": "val",
          "type": "text_eval"
        }
      },
      "slots": [
        "execute"
      ],
      "tokens": [
        "set text",
        "$OBJ",
        "$PROP",
        "$VAL"
      ]
    }
  },
  {
    "desc": "Shuffle Text: When called multiple times returns its inputs at random.",
    "group": [
      "format"
    ],
    "name": "shuffle_text",
    "uses": "run",
    "with": {
      "params": {
        "$ID": {
          "label": "id",
          "type": "text"
        },
        "$VALUES": {
          "label": "values",
          "repeats": true,
          "type": "text_eval"
        }
      },
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "shuffle text",
        "$ID",
        "$VALUES"
      ]
    }
  },
  {
    "desc": "Stopping Text: When called multiple times returns each of its inputs in turn, sticking to the last one.",
    "group": [
      "format"
    ],
    "name": "stopping_text",
    "uses": "run",
    "with": {
      "params": {
        "$ID": {
          "label": "id",
          "type": "text"
        },
        "$VALUES": {
          "label": "values",
          "repeats": true,
          "type": "text_eval"
        }
      },
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "stopping text",
        "$ID",
        "$VALUES"
      ]
    }
  },
  {
    "desc": "Subtract Numbers: Subtract two numbers.",
    "group": [
      "math"
    ],
    "name": "diff_of",
    "uses": "run",
    "with": {
      "params": {
        "$A": {
          "label": "a",
          "type": "number_eval"
        },
        "$B": {
          "label": "b",
          "type": "number_eval"
        }
      },
      "slots": [
        "number_eval"
      ],
      "tokens": [
        "(",
        "$A",
        "-",
        "$B",
        ")"
      ]
    }
  },
  {
    "desc": "Text Value: specifies a string value.",
    "group": [
      "literals"
    ],
    "name": "text_value",
    "spec": "{text:lines|quote}",
    "uses": "run",
    "with": {
      "slots": [
        "text_eval"
      ]
    }
  },
  {
    "desc": "Text List: specifies multiple string values.",
    "group": [
      "literals"
    ],
    "name": "texts",
    "uses": "run",
    "with": {
      "params": {
        "$VALUES": {
          "label": "values",
          "repeats": true,
          "type": "string"
        }
      },
      "slots": [
        "text_list_eval"
      ],
      "tokens": [
        "texts",
        "$VALUES"
      ]
    }
  },
  {
    "name": "logic",
    "uses": "group"
  },
  {
    "name": "objects",
    "uses": "group"
  },
  {
    "name": "format",
    "uses": "group"
  },
  {
    "name": "math",
    "uses": "group"
  },
  {
    "name": "literals",
    "uses": "group"
  },
  {
    "name": "exec",
    "uses": "group"
  },
  {
    "name": "cycle",
    "uses": "group"
  },
  {
    "name": "strings",
    "uses": "group"
  },
  {
    "name": "flow",
    "uses": "group"
  }
]
