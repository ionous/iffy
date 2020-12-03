/* generated using github.com/ionous/iffy/cmd/spec/spec.go */
const spec = [
  {
    "desc": "Assignments: Helper used when setting variables.",
    "name": "assignment",
    "uses": "slot"
  },
  {
    "desc": "Booleans: Statements which return true/false values.",
    "name": "bool_eval",
    "uses": "slot"
  },
  {
    "desc": "Comparison Types: Helper used when comparing two numbers, objects, pieces of text, etc.",
    "name": "comparator",
    "uses": "slot"
  },
  {
    "desc": "Action: Run a series of statements.",
    "name": "execute",
    "uses": "slot"
  },
  {
    "desc": "Number List: Statements which return a list of numbers.",
    "name": "num_list_eval",
    "uses": "slot"
  },
  {
    "desc": "Numbers: Statements which return a number.",
    "name": "number_eval",
    "uses": "slot"
  },
  {
    "desc": "Object: Statements which return an object.",
    "name": "object_eval",
    "uses": "slot"
  },
  {
    "desc": "Texts: Statements which return a record.",
    "name": "record_eval",
    "uses": "slot"
  },
  {
    "desc": "Record Lists:  Statements which return a list of records.",
    "name": "record_list_eval",
    "uses": "slot"
  },
  {
    "desc": "Texts: Statements which return text.",
    "name": "text_eval",
    "uses": "slot"
  },
  {
    "desc": "Text Lists: Statements which return a list of text.",
    "name": "text_list_eval",
    "uses": "slot"
  },
  {
    "group": [
      "hidden"
    ],
    "name": "activity",
    "spec": "{exe*execute}",
    "uses": "run",
    "with": {
      "slots": [
        "execute"
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
    "desc": "Always: returns true always.",
    "group": [
      "logic"
    ],
    "name": "always",
    "uses": "run",
    "with": {
      "params": {},
      "slots": [
        "bool_eval"
      ],
      "tokens": [
        "always"
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
    "group": [
      "patterns"
    ],
    "name": "argument",
    "spec": "its {name:variable_name} is {from:assignment}",
    "uses": "run",
    "with": {}
  },
  {
    "group": [
      "patterns"
    ],
    "name": "arguments",
    "spec": " when {arguments%args+argument}",
    "uses": "run",
    "with": {}
  },
  {
    "desc": "Assignment: Sets a variable to a value.",
    "group": [
      "variables"
    ],
    "name": "assign",
    "spec": "let {name:variable_name} be {from:assignment}",
    "uses": "run",
    "with": {
      "slots": [
        "execute"
      ]
    }
  },
  {
    "desc": "Assign Boolean: Assigns the passed boolean value.",
    "group": [
      "variables"
    ],
    "name": "assign_bool",
    "uses": "run",
    "with": {
      "params": {
        "$VAL": {
          "label": "val",
          "type": "bool_eval"
        }
      },
      "slots": [
        "assignment"
      ],
      "tokens": [
        "from bool",
        "$VAL"
      ]
    }
  },
  {
    "desc": "Assign Number: Assigns the passed number.",
    "group": [
      "variables"
    ],
    "name": "assign_num",
    "spec": "{val:number_eval}",
    "uses": "run",
    "with": {
      "slots": [
        "assignment"
      ]
    }
  },
  {
    "desc": "Assign Number List: Assigns the passed number list.",
    "group": [
      "variables"
    ],
    "name": "assign_num_list",
    "uses": "run",
    "with": {
      "params": {
        "$VALS": {
          "label": "vals",
          "type": "num_list_eval"
        }
      },
      "slots": [
        "assignment"
      ],
      "tokens": [
        "from num list",
        "$VALS"
      ]
    }
  },
  {
    "desc": "Assign Object: Assigns the passed object",
    "group": [
      "variables"
    ],
    "name": "assign_object",
    "uses": "run",
    "with": {
      "params": {
        "$VAL": {
          "label": "val",
          "type": "object_eval"
        }
      },
      "slots": [
        "assignment"
      ],
      "tokens": [
        "from object",
        "$VAL"
      ]
    }
  },
  {
    "desc": "Assign Record: Assigns the passed record.",
    "group": [
      "variables"
    ],
    "name": "assign_record",
    "uses": "run",
    "with": {
      "params": {
        "$VAL": {
          "label": "val",
          "type": "record_eval"
        }
      },
      "slots": [
        "assignment"
      ],
      "tokens": [
        "from record",
        "$VAL"
      ]
    }
  },
  {
    "desc": "Assign Record List: Assigns the passed record list.",
    "group": [
      "variables"
    ],
    "name": "assign_record_list",
    "uses": "run",
    "with": {
      "params": {
        "$VALS": {
          "label": "vals",
          "type": "record_list_eval"
        }
      },
      "slots": [
        "assignment"
      ],
      "tokens": [
        "from record list",
        "$VALS"
      ]
    }
  },
  {
    "desc": "Assign Text: Assigns the passed piece of text.",
    "group": [
      "variables"
    ],
    "name": "assign_text",
    "uses": "run",
    "with": {
      "params": {
        "$VAL": {
          "label": "val",
          "type": "text_eval"
        }
      },
      "slots": [
        "assignment"
      ],
      "tokens": [
        "from text",
        "$VAL"
      ]
    }
  },
  {
    "desc": "Assign Text List: Assigns the passed text list.",
    "group": [
      "variables"
    ],
    "name": "assign_text_list",
    "uses": "run",
    "with": {
      "params": {
        "$VALS": {
          "label": "vals",
          "type": "text_list_eval"
        }
      },
      "slots": [
        "assignment"
      ],
      "tokens": [
        "from text list",
        "$VALS"
      ]
    }
  },
  {
    "desc": "Greater Than or Equal To: The first value is larger than the second value.",
    "group": [
      "comparison"
    ],
    "name": "at_least",
    "spec": "\u003e=",
    "uses": "run",
    "with": {
      "slots": [
        "comparator"
      ]
    }
  },
  {
    "desc": "Less Than or Equal To: The first value is larger than the second value.",
    "group": [
      "comparison"
    ],
    "name": "at_most",
    "spec": "\u003c=",
    "uses": "run",
    "with": {
      "slots": [
        "comparator"
      ]
    }
  },
  {
    "desc": "Bool Value: specify an explicit true or false value.",
    "group": [
      "literals"
    ],
    "name": "bool_value",
    "spec": "{bool|quote}",
    "uses": "run",
    "with": {
      "slots": [
        "bool_eval"
      ]
    }
  },
  {
    "desc": "Bracket text: Sandwiches text printed during a block and puts them inside parenthesis '()'.",
    "group": [
      "printing"
    ],
    "name": "bracket_text",
    "uses": "run",
    "with": {
      "params": {
        "$GO": {
          "label": "go",
          "type": "activity"
        }
      },
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "bracket",
        "$GO"
      ]
    }
  },
  {
    "group": [
      "printing"
    ],
    "name": "buffer_text",
    "uses": "run",
    "with": {
      "params": {
        "$GO": {
          "label": "go",
          "type": "activity"
        }
      },
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "buffer",
        "$GO"
      ]
    }
  },
  {
    "name": "choose",
    "spec": "if {choose%if:bool_eval} then: {true:activity} else: {false:activity}",
    "uses": "run",
    "with": {
      "slots": [
        "execute"
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
    "desc": "List text: Separates words with commas, and 'and'.",
    "group": [
      "printing"
    ],
    "name": "comma_text",
    "uses": "run",
    "with": {
      "params": {
        "$GO": {
          "label": "go",
          "type": "activity"
        }
      },
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "commas",
        "$GO"
      ]
    }
  },
  {
    "desc": "Compare Numbers: True if eq,ne,gt,lt,ge,le two numbers.",
    "group": [
      "logic"
    ],
    "name": "compare_num",
    "spec": "{a:number_eval} {is:comparator} {b:number_eval}",
    "uses": "run",
    "with": {
      "slots": [
        "bool_eval"
      ]
    }
  },
  {
    "desc": "Compare Text: True if eq,ne,gt,lt,ge,le two strings ( lexical. )",
    "group": [
      "logic"
    ],
    "name": "compare_text",
    "spec": "{a:text_eval} {is:comparator} {b:text_eval}",
    "uses": "run",
    "with": {
      "slots": [
        "bool_eval"
      ]
    }
  },
  {
    "desc": "Copy Variable: Copy the contents of one variable to another.",
    "group": [
      "variables"
    ],
    "name": "copy_from",
    "uses": "run",
    "with": {
      "params": {
        "$NAME": {
          "label": "name",
          "type": "text"
        }
      },
      "slots": [
        "assignment"
      ],
      "tokens": [
        "copy from",
        "$NAME"
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
        "$PARTS": {
          "label": "parts",
          "repeats": true,
          "type": "text_eval"
        }
      },
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "cycle text",
        "$PARTS"
      ]
    }
  },
  {
    "desc": "Determine an activity",
    "group": [
      "patterns"
    ],
    "name": "determine_act",
    "uses": "run",
    "with": {
      "params": {
        "$ARGUMENTS": {
          "label": "arguments",
          "type": "arguments"
        },
        "$PATTERN": {
          "label": "pattern",
          "type": "text"
        }
      },
      "slots": [
        "execute"
      ],
      "tokens": [
        "determine act",
        "$PATTERN",
        "$ARGUMENTS"
      ]
    }
  },
  {
    "desc": "Determine a true/false value",
    "group": [
      "patterns"
    ],
    "name": "determine_bool",
    "spec": "the {true/false pattern%name:pattern_name}{?arguments}",
    "uses": "run",
    "with": {
      "slots": [
        "bool_eval"
      ]
    }
  },
  {
    "desc": "Determine a number",
    "group": [
      "patterns"
    ],
    "name": "determine_num",
    "spec": "the {number pattern%name:pattern_name}{?arguments}",
    "uses": "run",
    "with": {
      "slots": [
        "number_eval"
      ]
    }
  },
  {
    "desc": "Determine a list of numbers",
    "group": [
      "patterns"
    ],
    "name": "determine_num_list",
    "spec": "the {number list pattern%name:pattern_name}{?arguments}",
    "uses": "run",
    "with": {
      "slots": [
        "num_list_eval"
      ]
    }
  },
  {
    "desc": "Determine some text",
    "group": [
      "patterns"
    ],
    "name": "determine_text",
    "spec": "the {text pattern%name:pattern_name}{?arguments}",
    "uses": "run",
    "with": {
      "slots": [
        "text_eval"
      ]
    }
  },
  {
    "desc": "Determine a list of text",
    "group": [
      "patterns"
    ],
    "name": "determine_text_list",
    "spec": "the {text list pattern%name:pattern_name}{?arguments}",
    "uses": "run",
    "with": {
      "slots": [
        "text_list_eval"
      ]
    }
  },
  {
    "desc": "Subtract Numbers: Subtract two numbers.",
    "group": [
      "math"
    ],
    "name": "diff_of",
    "spec": "( {a:number_eval} - {b:number_eval} )",
    "uses": "run",
    "with": {
      "slots": [
        "number_eval"
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
    "desc": "Equal: Two values exactly match.",
    "group": [
      "comparison"
    ],
    "name": "equal",
    "spec": "is",
    "uses": "run",
    "with": {
      "slots": [
        "comparator"
      ]
    }
  },
  {
    "desc": "For Each Number: Loops over the passed list of numbers, or runs the 'else' activity if empty.",
    "group": [
      "exec"
    ],
    "name": "for_each_num",
    "uses": "run",
    "with": {
      "params": {
        "$ELSE": {
          "label": "else",
          "type": "activity"
        },
        "$GO": {
          "label": "go",
          "type": "activity"
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
    "desc": "For Each Text: Loops over the passed list of text, or runs the 'else' activity if empty.",
    "group": [
      "exec"
    ],
    "name": "for_each_text",
    "uses": "run",
    "with": {
      "params": {
        "$ELSE": {
          "label": "else",
          "type": "activity"
        },
        "$GO": {
          "label": "go",
          "type": "activity"
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
    "desc": "Get Field: Return the value of the named object property.",
    "group": [
      "objects"
    ],
    "name": "get_field",
    "spec": "the {field:text_eval} of {object:object_eval}",
    "uses": "run",
    "with": {
      "slots": [
        "bool_eval",
        "number_eval",
        "text_eval",
        "record_eval",
        "object_eval",
        "num_list_eval",
        "text_list_eval",
        "record_list_eval",
        "assignment"
      ]
    }
  },
  {
    "desc": "Get Variable: Return the value of the named variable.",
    "group": [
      "variables"
    ],
    "name": "get_var",
    "spec": "the {name:text_eval}",
    "uses": "run",
    "with": {
      "slots": [
        "bool_eval",
        "number_eval",
        "text_eval",
        "object_eval",
        "num_list_eval",
        "text_list_eval",
        "record_list_eval",
        "assignment"
      ]
    }
  },
  {
    "desc": "Greater Than: The first value is larger than the second value.",
    "group": [
      "comparison"
    ],
    "name": "greater_than",
    "spec": "\u003e",
    "uses": "run",
    "with": {
      "slots": [
        "comparator"
      ]
    }
  },
  {
    "desc": "Has Trait: Return true if noun is currently in the requested state.",
    "group": [
      "objects"
    ],
    "name": "has_trait",
    "spec": "{object%object:object_eval} is {trait:text_eval}",
    "uses": "run",
    "with": {
      "slots": [
        "bool_eval"
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
        "$KIND": {
          "label": "kind",
          "type": "text_eval"
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
        "is exact kind of",
        "$OBJ",
        "$KIND"
      ]
    }
  },
  {
    "desc": "Is Kind Of: True if the object is compatible with the named kind.",
    "group": [
      "objects"
    ],
    "name": "is_kind_of",
    "spec": "Is {object%obj:object_eval} a kind of {kind:singular_kind}",
    "uses": "run",
    "with": {
      "slots": [
        "bool_eval"
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
        "$TEST": {
          "label": "test",
          "type": "bool_eval"
        }
      },
      "slots": [
        "bool_eval"
      ],
      "tokens": [
        "is not true",
        "$TEST"
      ]
    }
  },
  {
    "desc": "Is True: Transparently returns the result of a boolean expression.",
    "group": [
      "logic"
    ],
    "name": "is_true",
    "spec": "{test:bool_eval} is true",
    "uses": "run",
    "with": {
      "slots": [
        "bool_eval"
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
        "$PARTS": {
          "label": "parts",
          "repeats": true,
          "type": "text_eval"
        },
        "$SEP": {
          "label": "sep",
          "type": "text_eval"
        }
      },
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "join",
        "$SEP",
        "$PARTS"
      ]
    }
  },
  {
    "desc": "Kind Of: Friendly name of the object's kind.",
    "group": [
      "objects"
    ],
    "name": "kind_of",
    "spec": "kind of {object%obj:object_eval}",
    "uses": "run",
    "with": {
      "slots": [
        "text_eval"
      ]
    }
  },
  {
    "desc": "Less Than: The first value is less than the second value.",
    "group": [
      "comparison"
    ],
    "name": "less_than",
    "spec": "\u003c",
    "uses": "run",
    "with": {
      "slots": [
        "comparator"
      ]
    }
  },
  {
    "desc": "Lines Value: specify one or more lines of text.",
    "group": [
      "literals"
    ],
    "name": "lines_value",
    "spec": "{lines|quote}",
    "uses": "run",
    "with": {}
  },
  {
    "desc": "Value of List: Get a value from a list. The first element is is index 1.",
    "group": [
      "list"
    ],
    "name": "list_at",
    "spec": "{list:text} entry {index:number}",
    "uses": "run",
    "with": {
      "slots": [
        "number_eval",
        "text_eval",
        "object_eval"
      ]
    }
  },
  {
    "desc": "For each in list: Loops over the elements in the passed list, or runs the 'else' activity if empty.",
    "group": [
      "list"
    ],
    "name": "list_each",
    "uses": "run",
    "with": {
      "params": {
        "$ELSE": {
          "label": "else",
          "type": "activity"
        },
        "$GO": {
          "label": "go",
          "type": "activity"
        },
        "$LIST": {
          "label": "list",
          "type": "text"
        },
        "$WITH": {
          "label": "with",
          "type": "text"
        }
      },
      "slots": [
        "execute"
      ],
      "tokens": [
        "each",
        "$LIST",
        "$WITH",
        "$GO",
        "$ELSE"
      ]
    }
  },
  {
    "desc": "Length of List: Determines the number of values in a list.",
    "group": [
      "list"
    ],
    "name": "list_len",
    "spec": "length of {list:text}",
    "uses": "run",
    "with": {
      "slots": [
        "number_eval"
      ]
    }
  },
  {
    "desc": "Map List: Transform the values from one list and place the results in another list.\n\t\tThe named pattern is called with two records 'in' and 'out' from the source and output lists respectively.",
    "group": [
      "list"
    ],
    "name": "list_map",
    "uses": "run",
    "with": {
      "params": {
        "$FROM_LIST": {
          "label": "from list",
          "type": "text"
        },
        "$TO_LIST": {
          "label": "to list",
          "type": "text"
        },
        "$USING_PATTERN": {
          "label": "using pattern",
          "type": "text"
        }
      },
      "slots": [
        "execute"
      ],
      "tokens": [
        "map",
        "$FROM_LIST",
        "$TO_LIST",
        "$USING_PATTERN"
      ]
    }
  },
  {
    "desc": "Pop from list: Remove an element from the front or back of a list.\nRuns an activity with the popped value, or runs the 'else' activity if the list was empty.",
    "group": [
      "list"
    ],
    "name": "list_pop",
    "uses": "run",
    "with": {
      "params": {
        "$ELSE": {
          "label": "else",
          "type": "activity"
        },
        "$FRONT": {
          "label": "front",
          "type": "bool"
        },
        "$GO": {
          "label": "go",
          "type": "activity"
        },
        "$LIST": {
          "label": "list",
          "type": "text"
        },
        "$WITH": {
          "label": "with",
          "type": "text"
        }
      },
      "slots": [
        "execute"
      ],
      "tokens": [
        "pop",
        "$LIST",
        "$WITH",
        "$FRONT",
        "$GO",
        "$ELSE"
      ]
    }
  },
  {
    "desc": "Push into list: Add elements to the front or back of a list.\nReturns the new length of the list.",
    "group": [
      "list"
    ],
    "name": "list_push",
    "spec": "push {list:text} {front} {inserting%insert?assignment}",
    "uses": "run",
    "with": {
      "slots": [
        "execute",
        "number_eval"
      ]
    }
  },
  {
    "desc": "Set Value of List: Overwrite an existing value in a list.",
    "group": [
      "list"
    ],
    "name": "list_set",
    "uses": "run",
    "with": {
      "params": {
        "$ELEMENT": {
          "label": "element",
          "type": "assignment"
        },
        "$INDEX": {
          "label": "index",
          "type": "number_eval"
        },
        "$LIST": {
          "label": "list",
          "type": "text"
        }
      },
      "slots": [
        "execute"
      ],
      "tokens": [
        "set",
        "$LIST",
        "$INDEX",
        "$ELEMENT"
      ]
    }
  },
  {
    "desc": "Slice of List: Create a new list from a section of another list.",
    "group": [
      "list"
    ],
    "name": "list_slice",
    "spec": "slice {list:text} {from entry%start?number} {ending before entry%eend?number}",
    "uses": "run",
    "with": {
      "slots": [
        "execute",
        "num_list_eval",
        "text_list_eval",
        "record_list_eval"
      ]
    }
  },
  {
    "desc": "Sort list: rearrange the elements in the named list by using the designated pattern to test pairs of elements.",
    "group": [
      "list"
    ],
    "name": "list_sort",
    "uses": "run",
    "with": {
      "params": {
        "$LIST": {
          "label": "list",
          "type": "text"
        },
        "$PATTERN": {
          "label": "pattern",
          "type": "text"
        }
      },
      "slots": [
        "execute"
      ],
      "tokens": [
        "sort",
        "$LIST",
        "$PATTERN"
      ]
    }
  },
  {
    "desc": "Splice into list: Modify a list by adding and removing elements.\nNote: the type of the elements being added must match the type of the list. \nText cant be added to a list of numbers, numbers cant be added to a list of text, \nand true/false values can't be added to a list.",
    "group": [
      "list"
    ],
    "name": "list_splice",
    "spec": "splice into {list:text} {at entry%start?number} {removing%remove?number} {inserting%insert?assignment}",
    "uses": "run",
    "with": {
      "slots": [
        "execute",
        "num_list_eval",
        "text_list_eval",
        "record_list_eval"
      ]
    }
  },
  {
    "name": "make",
    "uses": "run",
    "with": {
      "params": {
        "$ARGUMENTS": {
          "label": "arguments",
          "type": "arguments"
        },
        "$NAME": {
          "label": "name",
          "type": "text"
        }
      },
      "slots": [
        "object_eval"
      ],
      "tokens": [
        "make",
        "$NAME",
        "$ARGUMENTS"
      ]
    }
  },
  {
    "desc": "Lowercase: returns new text, with every letter turned into lowercase. \n\t\tFor example, \"shout\" from \"SHOUT\".",
    "group": [
      "format"
    ],
    "name": "make_lowercase",
    "spec": "{text:text_eval} in lowercase",
    "uses": "run",
    "with": {
      "slots": [
        "text_eval"
      ]
    }
  },
  {
    "desc": "Reverse text: returns new text flipped back to front. \n\t\tFor example, \"elppA\" from \"Apple\", or \"noon\" from \"noon\".",
    "group": [
      "format"
    ],
    "name": "make_reversed",
    "spec": "{text:text_eval} in reverse",
    "uses": "run",
    "with": {
      "slots": [
        "text_eval"
      ]
    }
  },
  {
    "desc": "Sentence case: returns new text, start each sentence with a capital letter. \n\t\tFor example, \"Empire Apple.\" from \"Empire apple.\".",
    "group": [
      "format"
    ],
    "name": "make_sentence_case",
    "spec": "{text:text_eval} in sentence-case",
    "uses": "run",
    "with": {
      "slots": [
        "text_eval"
      ]
    }
  },
  {
    "desc": "Title case: returns new text, starting each word with a capital letter. \n\t\tFor example, \"Empire Apple\" from \"empire apple\".",
    "group": [
      "format"
    ],
    "name": "make_title_case",
    "spec": "{text:text_eval} in title-case",
    "uses": "run",
    "with": {
      "slots": [
        "text_eval"
      ]
    }
  },
  {
    "desc": "Uppercase: returns new text, with every letter turned into uppercase. \n\t\tFor example, \"APPLE\" from \"apple\".",
    "group": [
      "format"
    ],
    "name": "make_uppercase",
    "spec": "{text:text_eval} in uppercase",
    "uses": "run",
    "with": {
      "slots": [
        "text_eval"
      ]
    }
  },
  {
    "desc": "Like: Determine whether the specified text is similar to the specified pattern.\n\t\tMatching is case-insensitive ( meaning, \"A\" matches \"a\" ) and there are two symbols with special meaning. \n\t\tA percent sign (\"%\") in the pattern matches any series of zero or more characters in the original text, \n\t\twhile an underscore matches (\"_\") any one single character. ",
    "group": [
      "matching"
    ],
    "name": "match_like",
    "spec": "{text:text_eval} is like {pattern:text_eval}",
    "uses": "run",
    "with": {
      "slots": [
        "bool_eval"
      ]
    }
  },
  {
    "desc": "Matches: Determine whether the specified text is similar to the specified regular expression.",
    "group": [
      "matching"
    ],
    "name": "matches",
    "spec": "{text:text_eval} matches {pattern:text}",
    "uses": "run",
    "with": {
      "slots": [
        "bool_eval"
      ]
    }
  },
  {
    "desc": "Move Variable: Move the contents of one variable to another, leaving the first variable blank.",
    "group": [
      "variables"
    ],
    "name": "move_from",
    "uses": "run",
    "with": {
      "params": {
        "$NAME": {
          "label": "name",
          "type": "text"
        }
      },
      "slots": [
        "assignment"
      ],
      "tokens": [
        "move from",
        "$NAME"
      ]
    }
  },
  {
    "desc": "Name Of: Full name of the object.",
    "group": [
      "objects"
    ],
    "name": "name_of",
    "spec": "name of {object%obj:object_eval}",
    "uses": "run",
    "with": {
      "slots": [
        "text_eval"
      ]
    }
  },
  {
    "desc": "Number Value: Specify a particular number.",
    "group": [
      "literals"
    ],
    "name": "num_value",
    "spec": "{num:number}",
    "uses": "run",
    "with": {
      "slots": [
        "number_eval"
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
          "type": "number"
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
    "desc": "Object Exists: Returns whether there is a noun of the specified name.",
    "group": [
      "objects"
    ],
    "name": "object_exists",
    "spec": "object named {name:text_eval}",
    "uses": "run",
    "with": {
      "slots": [
        "bool_eval"
      ]
    }
  },
  {
    "desc": "Pack: Puts a value into a record.",
    "group": [
      "variables"
    ],
    "name": "pack",
    "uses": "run",
    "with": {
      "params": {
        "$FIELD": {
          "label": "field",
          "type": "text"
        },
        "$FROM": {
          "label": "from",
          "type": "assignment"
        },
        "$RECORD": {
          "label": "record",
          "type": "record_eval"
        }
      },
      "slots": [
        "execute"
      ],
      "tokens": [
        "pack",
        "$RECORD",
        "$FIELD",
        "$FROM"
      ]
    }
  },
  {
    "desc": "Pluralize: Returns the plural form of a singular word. (ex.  apples for apple. )",
    "group": [
      "format"
    ],
    "name": "pluralize",
    "spec": "the plural of {text:text_eval}",
    "uses": "run",
    "with": {
      "slots": [
        "text_eval"
      ]
    }
  },
  {
    "desc": "A number as text: Writes a number using numerals, eg. '1'.",
    "group": [
      "printing"
    ],
    "name": "print_num",
    "spec": "as text {num:number_eval}",
    "uses": "run",
    "with": {
      "slots": [
        "text_eval"
      ]
    }
  },
  {
    "desc": "A number in words: Writes a number in plain english: eg. 'one'",
    "group": [
      "printing"
    ],
    "name": "print_num_word",
    "uses": "run",
    "with": {
      "params": {
        "$NUM": {
          "label": "num",
          "type": "number_eval"
        }
      },
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "print num word",
        "$NUM"
      ]
    }
  },
  {
    "desc": "Multiply Numbers: Multiply two numbers.",
    "group": [
      "math"
    ],
    "name": "product_of",
    "spec": "( {a:number_eval} * {b:number_eval} )",
    "uses": "run",
    "with": {
      "slots": [
        "number_eval"
      ]
    }
  },
  {
    "desc": "Divide Numbers: Divide one number by another.",
    "group": [
      "math"
    ],
    "name": "quotient_of",
    "spec": "( {a:number_eval} / {b:number_eval} )",
    "uses": "run",
    "with": {
      "slots": [
        "number_eval"
      ]
    }
  },
  {
    "desc": "Range of numbers: Generates a series of numbers.",
    "group": [
      "flow"
    ],
    "name": "range_over",
    "uses": "run",
    "with": {
      "params": {
        "$START": {
          "label": "start",
          "type": "number_eval"
        },
        "$STEP": {
          "label": "step",
          "type": "number_eval"
        },
        "$STOP": {
          "label": "stop",
          "type": "number_eval"
        }
      },
      "slots": [
        "num_list_eval"
      ],
      "tokens": [
        "range",
        "$START",
        "$STOP",
        "$STEP"
      ]
    }
  },
  {
    "desc": "Modulus Numbers: Divide one number by another, and return the remainder.",
    "group": [
      "math"
    ],
    "name": "remainder_of",
    "spec": "( {a:number_eval} % {b:number_eval} )",
    "uses": "run",
    "with": {
      "slots": [
        "number_eval"
      ]
    }
  },
  {
    "desc": "Render Template: Parse text using iffy templates. See: https://github.com/ionous/iffy/wiki/Templates",
    "group": [
      "format"
    ],
    "name": "render_template",
    "spec": "the template {lines%template:lines|quote}",
    "uses": "run",
    "with": {
      "slots": [
        "text_eval"
      ]
    }
  },
  {
    "desc": "Say: print some bit of text to the player.",
    "group": [
      "printing"
    ],
    "name": "say_text",
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
    "desc": "Set Field: Sets the named field to the assigned value.",
    "group": [
      "objects"
    ],
    "name": "set_field",
    "uses": "run",
    "with": {
      "params": {
        "$FIELD": {
          "label": "field",
          "type": "text"
        },
        "$FROM": {
          "label": "from",
          "type": "assignment"
        },
        "$OBJ": {
          "label": "obj",
          "type": "object_eval"
        }
      },
      "slots": [
        "execute"
      ],
      "tokens": [
        "set field",
        "$OBJ",
        "$FIELD",
        "$FROM"
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
        "$PARTS": {
          "label": "parts",
          "repeats": true,
          "type": "text_eval"
        }
      },
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "shuffle text",
        "$PARTS"
      ]
    }
  },
  {
    "desc": "Singularize: Returns the singular form of a plural word. (ex. apple for apples )",
    "group": [
      "format"
    ],
    "name": "singularize",
    "spec": "the singular {text:text_eval}",
    "uses": "run",
    "with": {
      "slots": [
        "text_eval"
      ]
    }
  },
  {
    "desc": "Slash text: Separates words with left-leaning slashes '/'.",
    "group": [
      "printing"
    ],
    "name": "slash_text",
    "uses": "run",
    "with": {
      "params": {
        "$GO": {
          "label": "go",
          "type": "activity"
        }
      },
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "slash",
        "$GO"
      ]
    }
  },
  {
    "desc": "Span Text: Writes text with spaces between words.",
    "group": [
      "printing"
    ],
    "name": "span_text",
    "uses": "run",
    "with": {
      "params": {
        "$GO": {
          "label": "go",
          "type": "activity"
        }
      },
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "span",
        "$GO"
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
        "$PARTS": {
          "label": "parts",
          "repeats": true,
          "type": "text_eval"
        }
      },
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "stopping text",
        "$PARTS"
      ]
    }
  },
  {
    "desc": "Add Numbers: Add two numbers.",
    "group": [
      "math"
    ],
    "name": "sum_of",
    "spec": "( {a:number_eval} + {b:number_eval} )",
    "uses": "run",
    "with": {
      "slots": [
        "number_eval"
      ]
    }
  },
  {
    "desc": "Text Value: specify a small bit of text.",
    "group": [
      "literals"
    ],
    "name": "text_value",
    "spec": "{text}",
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
          "type": "text"
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
    "desc": "Not Equal To: Two values don't match exactly.",
    "group": [
      "comparison"
    ],
    "name": "unequal",
    "spec": "is not",
    "uses": "run",
    "with": {
      "slots": [
        "comparator"
      ]
    }
  },
  {
    "desc": "Unpack: Get a value from a record.",
    "group": [
      "variables"
    ],
    "name": "unpack",
    "spec": "unpack {field:text_eval} from {record:record_eval}",
    "uses": "run",
    "with": {
      "slots": [
        "bool_eval",
        "number_eval",
        "text_eval",
        "record_eval",
        "object_eval",
        "num_list_eval",
        "text_list_eval",
        "record_list_eval",
        "assignment"
      ]
    }
  },
  {
    "name": "comparison",
    "uses": "group"
  },
  {
    "name": "cycle",
    "uses": "group"
  },
  {
    "name": "exec",
    "uses": "group"
  },
  {
    "name": "flow",
    "uses": "group"
  },
  {
    "name": "format",
    "uses": "group"
  },
  {
    "name": "hidden",
    "uses": "group"
  },
  {
    "name": "list",
    "uses": "group"
  },
  {
    "name": "literals",
    "uses": "group"
  },
  {
    "name": "logic",
    "uses": "group"
  },
  {
    "name": "matching",
    "uses": "group"
  },
  {
    "name": "math",
    "uses": "group"
  },
  {
    "name": "objects",
    "uses": "group"
  },
  {
    "name": "patterns",
    "uses": "group"
  },
  {
    "name": "printing",
    "uses": "group"
  },
  {
    "name": "strings",
    "uses": "group"
  },
  {
    "name": "variables",
    "uses": "group"
  }
]
