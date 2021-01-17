/* generated using github.com/ionous/iffy/cmd/spec/spec.go */
const spec = [
  {
    "desc": "debug level: Debug level.",
    "group": [
      "list"
    ],
    "name": "debug_level",
    "uses": "str",
    "with": {
      "params": {
        "$FIX": "fix",
        "$NOTE": "note",
        "$TO_DO": {
          "label": "toDo",
          "value": "to do"
        },
        "$WARNING": "warning"
      },
      "tokens": [
        "$NOTE",
        " or ",
        "$TO_DO",
        " or ",
        "$WARNING",
        " or ",
        "$FIX"
      ]
    }
  },
  {
    "desc": "List Case: When sorting, treat uppercase and lowercase versions of letters the same.",
    "group": [
      "list"
    ],
    "name": "list_case",
    "uses": "str",
    "with": {
      "params": {
        "$FALSE": {
          "label": "includeCase",
          "value": "include_case"
        },
        "$TRUE": {
          "label": "ignoreCase",
          "value": "ignore_case"
        }
      },
      "tokens": [
        "$FALSE",
        " or ",
        "$TRUE"
      ]
    }
  },
  {
    "desc": "List Edge: Put elements at the front or back of a list.",
    "group": [
      "list"
    ],
    "name": "list_edge",
    "uses": "str",
    "with": {
      "params": {
        "$FALSE": {
          "label": "atBack",
          "value": "at_back"
        },
        "$TRUE": {
          "label": "atFront",
          "value": "at_front"
        }
      },
      "tokens": [
        "$FALSE",
        " or ",
        "$TRUE"
      ]
    }
  },
  {
    "desc": "List Order: Sort larger values towards the end of a list.",
    "group": [
      "list"
    ],
    "name": "list_order",
    "uses": "str",
    "with": {
      "params": {
        "$FALSE": "ascending",
        "$TRUE": "descending"
      },
      "tokens": [
        "$FALSE",
        " or ",
        "$TRUE"
      ]
    }
  },
  {
    "name": "relation_name",
    "uses": "str",
    "with": {
      "params": {
        "$RELATION_NAME": "relation name"
      },
      "tokens": [
        "$RELATION_NAME"
      ]
    }
  },
  {
    "name": "variable_name",
    "uses": "str",
    "with": {
      "params": {
        "$VARIABLE_NAME": "variable name"
      },
      "tokens": [
        "$VARIABLE_NAME"
      ]
    }
  },
  {
    "desc": "Assignments: Helper for setting variables.",
    "name": "assignment",
    "uses": "slot"
  },
  {
    "desc": "Booleans: Statements which return true/false values.",
    "name": "bool_eval",
    "uses": "slot"
  },
  {
    "desc": "brancher: Helper for choose action.",
    "name": "brancher",
    "uses": "slot"
  },
  {
    "desc": "Comparison Types: Helper for comparing values.",
    "name": "comparator",
    "uses": "slot"
  },
  {
    "desc": "Action: Run a series of statements.",
    "name": "execute",
    "uses": "slot"
  },
  {
    "desc": "fields: Helper for setting fields.",
    "name": "fields",
    "uses": "slot"
  },
  {
    "desc": "List iterator: Helper for accessing lists.",
    "name": "list_iterator",
    "uses": "slot"
  },
  {
    "desc": "List source: Helper for accessing lists.",
    "name": "list_source",
    "uses": "slot"
  },
  {
    "desc": "List target: Helper for accessing lists.",
    "name": "list_target",
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
    "name": "comparison",
    "uses": "group"
  },
  {
    "name": "cycle",
    "uses": "group"
  },
  {
    "name": "debug",
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
    "name": "relations",
    "uses": "group"
  },
  {
    "name": "strings",
    "uses": "group"
  },
  {
    "name": "variables",
    "uses": "group"
  },
  {
    "group": [
      "hidden"
    ],
    "name": "activity",
    "spec": "{exe*execute}",
    "uses": "flow",
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
    "uses": "flow",
    "with": {
      "params": {
        "$TEST": {
          "label": "test",
          "repeats": true,
          "type": "bool_eval"
        }
      },
      "roles": "QZSZK",
      "slots": [
        "bool_eval"
      ],
      "tokens": [
        "allTrue",
        " ",
        "test",
        ": ",
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
    "uses": "flow",
    "with": {
      "params": {},
      "roles": "E",
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
    "uses": "flow",
    "with": {
      "params": {
        "$TEST": {
          "label": "test",
          "repeats": true,
          "type": "bool_eval"
        }
      },
      "roles": "QZSZK",
      "slots": [
        "bool_eval"
      ],
      "tokens": [
        "anyTrue",
        " ",
        "test",
        ": ",
        "$TEST"
      ]
    }
  },
  {
    "group": [
      "patterns"
    ],
    "name": "argument",
    "spec": " {name:variable_name}: {from:assignment}",
    "uses": "flow"
  },
  {
    "group": [
      "patterns"
    ],
    "name": "arguments",
    "spec": " {arguments%args+argument|comma-and}",
    "uses": "flow"
  },
  {
    "desc": "AsNum: Define the name of a number variable.",
    "name": "as_num",
    "uses": "flow",
    "with": {
      "params": {
        "$VAR": {
          "label": "var",
          "type": "variable_name"
        }
      },
      "roles": "SZK",
      "slots": [
        "list_iterator"
      ],
      "tokens": [
        "asNum",
        ": ",
        "$VAR"
      ]
    }
  },
  {
    "desc": "AsRec: Define the name of a record variable.",
    "name": "as_rec",
    "uses": "flow",
    "with": {
      "params": {
        "$VAR": {
          "label": "var",
          "type": "variable_name"
        }
      },
      "roles": "SZK",
      "slots": [
        "list_iterator"
      ],
      "tokens": [
        "asRec",
        ": ",
        "$VAR"
      ]
    }
  },
  {
    "desc": "AsTxt: Define the name of a text variable.",
    "name": "as_txt",
    "uses": "flow",
    "with": {
      "params": {
        "$VAR": {
          "label": "var",
          "type": "variable_name"
        }
      },
      "roles": "SZK",
      "slots": [
        "list_iterator"
      ],
      "tokens": [
        "asTxt",
        ": ",
        "$VAR"
      ]
    }
  },
  {
    "desc": "Assignment: Sets a variable to a value.",
    "group": [
      "variables"
    ],
    "name": "assign",
    "uses": "flow",
    "with": {
      "params": {
        "$FROM": {
          "label": "be",
          "type": "assignment"
        },
        "$VAR": {
          "label": "var",
          "type": "variable_name"
        }
      },
      "roles": "CZKZSZKT",
      "slots": [
        "execute"
      ],
      "tokens": [
        "let",
        ": ",
        "$VAR",
        ", ",
        "be",
        ": ",
        "$FROM",
        "."
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
    "uses": "flow",
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
    "uses": "flow",
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
    "spec": "{bool}",
    "uses": "flow",
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
    "uses": "flow",
    "with": {
      "params": {
        "$GO": {
          "label": "go",
          "optional": true,
          "type": "activity"
        }
      },
      "roles": "QZSZK",
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "bracket",
        " ",
        "go",
        ": ",
        "$GO"
      ]
    }
  },
  {
    "desc": "Break: In a repeating loop, exit the loop.",
    "group": [
      "flow"
    ],
    "name": "break",
    "uses": "flow",
    "with": {
      "params": {},
      "roles": "E",
      "slots": [
        "execute"
      ],
      "tokens": [
        "break"
      ]
    }
  },
  {
    "group": [
      "printing"
    ],
    "name": "buffer_text",
    "uses": "flow",
    "with": {
      "params": {
        "$GO": {
          "label": "go",
          "optional": true,
          "type": "activity"
        }
      },
      "roles": "QZSZK",
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "buffer",
        " ",
        "go",
        ": ",
        "$GO"
      ]
    }
  },
  {
    "name": "choose_action",
    "uses": "flow",
    "with": {
      "params": {
        "$DO": {
          "label": "do",
          "type": "activity"
        },
        "$ELSE": {
          "label": "else",
          "optional": true,
          "type": "brancher"
        },
        "$IF": {
          "label": "choose",
          "type": "bool_eval"
        }
      },
      "roles": "CZKZSZKZKT",
      "slots": [
        "execute"
      ],
      "tokens": [
        "if",
        ": ",
        "$IF",
        ", ",
        "do",
        ": ",
        "$DO",
        ", ",
        "$ELSE",
        "."
      ]
    }
  },
  {
    "name": "choose_more",
    "uses": "flow",
    "with": {
      "params": {
        "$DO": {
          "label": "do",
          "type": "activity"
        },
        "$ELSE": {
          "label": "else",
          "optional": true,
          "type": "brancher"
        },
        "$IF": {
          "label": "choose",
          "type": "bool_eval"
        }
      },
      "roles": "SZKZSZKZK",
      "slots": [
        "brancher"
      ],
      "tokens": [
        "elseIf",
        ": ",
        "$IF",
        ", ",
        "do",
        ": ",
        "$DO",
        ", ",
        "$ELSE"
      ]
    }
  },
  {
    "name": "choose_nothing_else",
    "uses": "flow",
    "with": {
      "params": {
        "$DO": {
          "label": "do",
          "type": "activity"
        }
      },
      "roles": "SZK",
      "slots": [
        "brancher"
      ],
      "tokens": [
        "elseDo",
        ": ",
        "$DO"
      ]
    }
  },
  {
    "desc": "Choose Number: Pick one of two numbers based on a boolean test.",
    "group": [
      "math"
    ],
    "name": "choose_num",
    "uses": "flow",
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
      "roles": "QZSZKZSZKZSZK",
      "slots": [
        "number_eval"
      ],
      "tokens": [
        "chooseNum",
        " ",
        "if",
        ": ",
        "$IF",
        ", ",
        "true",
        ": ",
        "$TRUE",
        ", ",
        "false",
        ": ",
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
    "uses": "flow",
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
      "roles": "QZSZKZSZKZSZK",
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "chooseText",
        " ",
        "if",
        ": ",
        "$IF",
        ", ",
        "true",
        ": ",
        "$TRUE",
        ", ",
        "false",
        ": ",
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
    "uses": "flow",
    "with": {
      "params": {
        "$GO": {
          "label": "go",
          "optional": true,
          "type": "activity"
        }
      },
      "roles": "QZSZK",
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "commas",
        " ",
        "go",
        ": ",
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
    "uses": "flow",
    "with": {
      "params": {
        "$A": {
          "label": "num",
          "type": "number_eval"
        },
        "$B": {
          "label": "b",
          "type": "number_eval"
        },
        "$IS": {
          "label": "is",
          "type": "comparator"
        }
      },
      "roles": "FZSZKZKZK",
      "slots": [
        "bool_eval"
      ],
      "tokens": [
        "is",
        " ",
        "num",
        ": ",
        "$A",
        ", ",
        "$IS",
        ", ",
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
    "uses": "flow",
    "with": {
      "params": {
        "$A": {
          "label": "txt",
          "type": "text_eval"
        },
        "$B": {
          "label": "b",
          "type": "text_eval"
        },
        "$IS": {
          "label": "is",
          "type": "comparator"
        }
      },
      "roles": "FZSZKZKZK",
      "slots": [
        "bool_eval"
      ],
      "tokens": [
        "is",
        " ",
        "txt",
        ": ",
        "$A",
        ", ",
        "$IS",
        ", ",
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
    "uses": "flow",
    "with": {
      "params": {
        "$PARTS": {
          "label": "parts",
          "repeats": true,
          "type": "text_eval"
        }
      },
      "roles": "QZSZK",
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "cycleText",
        " ",
        "parts",
        ": ",
        "$PARTS"
      ]
    }
  },
  {
    "group": [
      "debug"
    ],
    "name": "debug_log",
    "uses": "flow",
    "with": {
      "params": {
        "$LEVEL": {
          "label": "level",
          "type": "number"
        },
        "$VALUE": {
          "label": "value",
          "type": "assignment"
        }
      },
      "roles": "CZKZKT",
      "slots": [
        "execute"
      ],
      "tokens": [
        "log",
        ": ",
        "$VALUE",
        ", ",
        "$LEVEL",
        "."
      ]
    }
  },
  {
    "desc": "determine act: Determine an activity",
    "group": [
      "patterns"
    ],
    "name": "determine_act",
    "spec": "determine {activity%name:pattern_name}{?arguments}",
    "uses": "flow",
    "with": {
      "slots": [
        "execute"
      ]
    }
  },
  {
    "desc": "determine bool: Determine a true/false value",
    "group": [
      "patterns"
    ],
    "name": "determine_bool",
    "spec": "{true/false pattern%name:pattern_name}{?arguments}",
    "uses": "flow",
    "with": {
      "slots": [
        "bool_eval"
      ]
    }
  },
  {
    "desc": "determine num: Determine a number",
    "group": [
      "patterns"
    ],
    "name": "determine_num",
    "spec": "{number pattern%name:pattern_name}{?arguments}",
    "uses": "flow",
    "with": {
      "slots": [
        "number_eval"
      ]
    }
  },
  {
    "desc": "determine num list: Determine a list of numbers",
    "group": [
      "patterns"
    ],
    "name": "determine_num_list",
    "spec": "{number list pattern%name:pattern_name}{?arguments}",
    "uses": "flow",
    "with": {
      "slots": [
        "num_list_eval"
      ]
    }
  },
  {
    "desc": "determine text: Determine some text",
    "group": [
      "patterns"
    ],
    "name": "determine_text",
    "spec": "{text pattern%name:pattern_name}{?arguments}",
    "uses": "flow",
    "with": {
      "slots": [
        "text_eval"
      ]
    }
  },
  {
    "desc": "determine text list: Determine a list of text",
    "group": [
      "patterns"
    ],
    "name": "determine_text_list",
    "spec": "{text list pattern%name:pattern_name}{?arguments}",
    "uses": "flow",
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
    "uses": "flow",
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
    "uses": "flow",
    "with": {
      "params": {
        "$REASON": {
          "label": "reason",
          "type": "text"
        }
      },
      "roles": "QZSZK",
      "slots": [
        "execute"
      ],
      "tokens": [
        "doNothing",
        " ",
        "reason",
        ": ",
        "$REASON"
      ]
    }
  },
  {
    "desc": "Equal: Two values exactly match.",
    "group": [
      "comparison"
    ],
    "name": "equal",
    "spec": "==",
    "uses": "flow",
    "with": {
      "slots": [
        "comparator"
      ]
    }
  },
  {
    "desc": "Erase: Remove one or more values from a list",
    "name": "erase_edge",
    "uses": "flow",
    "with": {
      "params": {
        "$AT_EDGE": {
          "label": "atEdge",
          "type": "list_edge"
        },
        "$FROM": {
          "label": "from",
          "type": "list_source"
        }
      },
      "roles": "CZKZKT",
      "slots": [
        "execute"
      ],
      "tokens": [
        "erase",
        ": ",
        "$FROM",
        ", ",
        "$AT_EDGE",
        "."
      ]
    }
  },
  {
    "desc": "Erase: remove one or more values from a list",
    "name": "erase_index",
    "uses": "flow",
    "with": {
      "params": {
        "$AT_INDEX": {
          "label": "atIndex",
          "type": "number_eval"
        },
        "$COUNT": {
          "label": "count",
          "type": "number_eval"
        },
        "$FROM": {
          "label": "from",
          "type": "list_source"
        }
      },
      "roles": "CZKZKZSZKT",
      "slots": [
        "execute"
      ],
      "tokens": [
        "erase",
        ": ",
        "$COUNT",
        ", ",
        "$FROM",
        ", ",
        "atIndex",
        ": ",
        "$AT_INDEX",
        "."
      ]
    }
  },
  {
    "desc": "Erasing from list: Erase elements from the front or back of a list.\nRuns an activity with a list containing the erased values; the list can be empty if nothing was erased.",
    "group": [
      "list"
    ],
    "name": "erasing",
    "uses": "flow",
    "with": {
      "params": {
        "$AS": {
          "label": "as",
          "type": "text"
        },
        "$AT_INDEX": {
          "label": "atIndex",
          "type": "number_eval"
        },
        "$COUNT": {
          "label": "count",
          "type": "number_eval"
        },
        "$DO": {
          "label": "do",
          "type": "activity"
        },
        "$FROM": {
          "label": "from",
          "type": "list_source"
        }
      },
      "roles": "CZKZKZSZKZSZKZSZKT",
      "slots": [
        "execute"
      ],
      "tokens": [
        "erasing",
        ": ",
        "$COUNT",
        ", ",
        "$FROM",
        ", ",
        "atIndex",
        ": ",
        "$AT_INDEX",
        ", ",
        "as",
        ": ",
        "$AS",
        ", ",
        "do",
        ": ",
        "$DO",
        "."
      ]
    }
  },
  {
    "desc": "From Bool: Assigns the calculated boolean value.",
    "group": [
      "variables"
    ],
    "name": "from_bool",
    "uses": "flow",
    "with": {
      "params": {
        "$VAL": {
          "label": "val",
          "type": "bool_eval"
        }
      },
      "roles": "FZK",
      "slots": [
        "assignment"
      ],
      "tokens": [
        "fromBool",
        ": ",
        "$VAL"
      ]
    }
  },
  {
    "desc": "From Name: Assigns the calculated piece of name.",
    "group": [
      "variables"
    ],
    "name": "from_name",
    "uses": "flow",
    "with": {
      "params": {
        "$VAL": {
          "label": "val",
          "type": "text_eval"
        }
      },
      "roles": "FZK",
      "slots": [
        "assignment"
      ],
      "tokens": [
        "fromName",
        ": ",
        "$VAL"
      ]
    }
  },
  {
    "desc": "From Number: Assigns the calculated number.",
    "group": [
      "variables"
    ],
    "name": "from_num",
    "spec": "{val:number_eval}",
    "uses": "flow",
    "with": {
      "slots": [
        "assignment"
      ]
    }
  },
  {
    "desc": "FromNumList: Uses a list of numbers",
    "name": "from_num_list",
    "uses": "flow",
    "with": {
      "params": {
        "$VAR": {
          "label": "var",
          "type": "variable_name"
        }
      },
      "roles": "SZK",
      "slots": [
        "list_source"
      ],
      "tokens": [
        "fromNumList",
        ": ",
        "$VAR"
      ]
    }
  },
  {
    "desc": "From Numbers: Assigns the calculated numbers.",
    "group": [
      "variables"
    ],
    "name": "from_nums",
    "uses": "flow",
    "with": {
      "params": {
        "$VALS": {
          "label": "vals",
          "type": "num_list_eval"
        }
      },
      "roles": "FZK",
      "slots": [
        "assignment"
      ],
      "tokens": [
        "fromNumbers",
        ": ",
        "$VALS"
      ]
    }
  },
  {
    "desc": "From Object: Assigns the calculated object",
    "group": [
      "variables"
    ],
    "name": "from_object",
    "uses": "flow",
    "with": {
      "params": {
        "$VAL": {
          "label": "val",
          "type": "object_eval"
        }
      },
      "roles": "FZK",
      "slots": [
        "assignment"
      ],
      "tokens": [
        "fromObject",
        ": ",
        "$VAL"
      ]
    }
  },
  {
    "desc": "FromRecList: Uses a list of records",
    "name": "from_rec_list",
    "uses": "flow",
    "with": {
      "params": {
        "$VAR": {
          "label": "var",
          "type": "variable_name"
        }
      },
      "roles": "SZK",
      "slots": [
        "list_source"
      ],
      "tokens": [
        "fromRecList",
        ": ",
        "$VAR"
      ]
    }
  },
  {
    "desc": "From Record: Assigns the calculated record.",
    "group": [
      "variables"
    ],
    "name": "from_record",
    "uses": "flow",
    "with": {
      "params": {
        "$VAL": {
          "label": "val",
          "type": "record_eval"
        }
      },
      "roles": "FZK",
      "slots": [
        "assignment"
      ],
      "tokens": [
        "fromRecord",
        ": ",
        "$VAL"
      ]
    }
  },
  {
    "desc": "From Records: Assigns the calculated records.",
    "group": [
      "variables"
    ],
    "name": "from_records",
    "uses": "flow",
    "with": {
      "params": {
        "$VALS": {
          "label": "vals",
          "type": "record_list_eval"
        }
      },
      "roles": "FZK",
      "slots": [
        "assignment"
      ],
      "tokens": [
        "fromRecords",
        ": ",
        "$VALS"
      ]
    }
  },
  {
    "desc": "From Text: Assigns the calculated piece of text.",
    "group": [
      "variables"
    ],
    "name": "from_text",
    "uses": "flow",
    "with": {
      "params": {
        "$VAL": {
          "label": "val",
          "type": "text_eval"
        }
      },
      "roles": "FZK",
      "slots": [
        "assignment"
      ],
      "tokens": [
        "fromText",
        ": ",
        "$VAL"
      ]
    }
  },
  {
    "desc": "From Texts: Assigns the calculated texts.",
    "group": [
      "variables"
    ],
    "name": "from_texts",
    "uses": "flow",
    "with": {
      "params": {
        "$VALS": {
          "label": "vals",
          "type": "text_list_eval"
        }
      },
      "roles": "FZK",
      "slots": [
        "assignment"
      ],
      "tokens": [
        "fromTexts",
        ": ",
        "$VALS"
      ]
    }
  },
  {
    "desc": "FromTxtList: Uses a list of text",
    "name": "from_txt_list",
    "uses": "flow",
    "with": {
      "params": {
        "$VAR": {
          "label": "var",
          "type": "variable_name"
        }
      },
      "roles": "SZK",
      "slots": [
        "list_source"
      ],
      "tokens": [
        "fromTxtList",
        ": ",
        "$VAR"
      ]
    }
  },
  {
    "desc": "Get Field: Return the value of the named object property.",
    "group": [
      "objects"
    ],
    "name": "get_field",
    "spec": "the {field:text} of {object:object_eval}",
    "uses": "flow",
    "with": {
      "slots": [
        "bool_eval",
        "number_eval",
        "text_eval",
        "record_eval",
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
    "uses": "flow",
    "with": {
      "params": {
        "$NAME": {
          "label": "name",
          "type": "text"
        }
      },
      "roles": "FZK",
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
      ],
      "tokens": [
        "var",
        ": ",
        "$NAME"
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
    "uses": "flow",
    "with": {
      "slots": [
        "comparator"
      ]
    }
  },
  {
    "group": [
      "logic"
    ],
    "name": "has_dominion",
    "uses": "flow",
    "with": {
      "params": {
        "$NAME": {
          "label": "name",
          "type": "text"
        }
      },
      "roles": "QZSZK",
      "slots": [
        "bool_eval"
      ],
      "tokens": [
        "hasDominion",
        " ",
        "name",
        ": ",
        "$NAME"
      ]
    }
  },
  {
    "desc": "Has Trait: Return true if noun is currently in the requested state.",
    "group": [
      "objects"
    ],
    "name": "has_trait",
    "spec": "{object:object_eval} is {trait:text_eval}",
    "uses": "flow",
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
    "uses": "flow",
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
      "roles": "QZSZKZSZK",
      "slots": [
        "bool_eval"
      ],
      "tokens": [
        "includes",
        " ",
        "text",
        ": ",
        "$TEXT",
        ", ",
        "part",
        ": ",
        "$PART"
      ]
    }
  },
  {
    "desc": "IntoNumList: Targets a list of numbers",
    "name": "into_num_list",
    "uses": "flow",
    "with": {
      "params": {
        "$VAR": {
          "label": "var",
          "type": "variable_name"
        }
      },
      "roles": "SZK",
      "slots": [
        "list_target"
      ],
      "tokens": [
        "intoNumList",
        ": ",
        "$VAR"
      ]
    }
  },
  {
    "desc": "IntoObj: Targets an object with a predetermined name",
    "name": "into_obj",
    "uses": "flow",
    "with": {
      "params": {
        "$VAR": {
          "label": "var",
          "type": "variable_name"
        }
      },
      "roles": "SZK",
      "slots": [
        "fields"
      ],
      "tokens": [
        "intoObj",
        ": ",
        "$VAR"
      ]
    }
  },
  {
    "desc": "IntoObjNamed: Targets an object with a computed name",
    "name": "into_obj_named",
    "uses": "flow",
    "with": {
      "params": {
        "$OBJ_NAME": {
          "label": "objName",
          "type": "text_eval"
        }
      },
      "roles": "SZK",
      "slots": [
        "fields"
      ],
      "tokens": [
        "intoObjNamed",
        ": ",
        "$OBJ_NAME"
      ]
    }
  },
  {
    "desc": "IntoRec: Targets a record stored in a variable",
    "name": "into_rec",
    "uses": "flow",
    "with": {
      "params": {
        "$VAR": {
          "label": "var",
          "type": "variable_name"
        }
      },
      "roles": "SZK",
      "slots": [
        "fields"
      ],
      "tokens": [
        "intoRec",
        ": ",
        "$VAR"
      ]
    }
  },
  {
    "desc": "IntoRecList: Targets a list of records",
    "name": "into_rec_list",
    "uses": "flow",
    "with": {
      "params": {
        "$VAR": {
          "label": "var",
          "type": "variable_name"
        }
      },
      "roles": "SZK",
      "slots": [
        "list_target"
      ],
      "tokens": [
        "intoRecList",
        ": ",
        "$VAR"
      ]
    }
  },
  {
    "desc": "IntoTxtList: Targets a list of text",
    "name": "into_txt_list",
    "uses": "flow",
    "with": {
      "params": {
        "$VAR": {
          "label": "var",
          "type": "variable_name"
        }
      },
      "roles": "SZK",
      "slots": [
        "list_target"
      ],
      "tokens": [
        "intoTxtList",
        ": ",
        "$VAR"
      ]
    }
  },
  {
    "desc": "Is Empty: True if the text is empty.",
    "group": [
      "strings"
    ],
    "name": "is_empty",
    "uses": "flow",
    "with": {
      "params": {
        "$TEXT": {
          "label": "text",
          "type": "text_eval"
        }
      },
      "roles": "FZK",
      "slots": [
        "bool_eval"
      ],
      "tokens": [
        "isEmpty",
        ": ",
        "$TEXT"
      ]
    }
  },
  {
    "desc": "Is Kind Of: True if the object is compatible with the named kind.",
    "group": [
      "objects"
    ],
    "name": "is_kind_of",
    "uses": "flow",
    "with": {
      "params": {
        "$KIND": {
          "label": "is",
          "type": "text"
        },
        "$OBJECT": {
          "label": "object",
          "type": "object_eval"
        }
      },
      "roles": "FZKZSZK",
      "slots": [
        "bool_eval"
      ],
      "tokens": [
        "kindOf",
        ": ",
        "$OBJECT",
        ", ",
        "is",
        ": ",
        "$KIND"
      ]
    }
  },
  {
    "desc": "Join Strings: Returns multiple pieces of text as a single new piece of text.",
    "group": [
      "strings"
    ],
    "name": "join",
    "uses": "flow",
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
      "roles": "QZSZKZSZK",
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "join",
        " ",
        "sep",
        ": ",
        "$SEP",
        ", ",
        "parts",
        ": ",
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
    "spec": "kind of {object:object_eval}",
    "uses": "flow",
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
    "uses": "flow",
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
    "uses": "flow"
  },
  {
    "desc": "Value of List: Get a value from a list. The first element is is index 1.",
    "group": [
      "list"
    ],
    "name": "list_at",
    "spec": "list {list:assignment} at {index:number_eval}",
    "uses": "flow",
    "with": {
      "slots": [
        "number_eval",
        "text_eval",
        "record_eval"
      ]
    }
  },
  {
    "desc": "For each in list: Loops over the elements in the passed list, or runs the 'else' activity if empty.",
    "group": [
      "list"
    ],
    "name": "list_each",
    "uses": "flow",
    "with": {
      "params": {
        "$AS": {
          "label": "as",
          "type": "list_iterator"
        },
        "$DO": {
          "label": "do",
          "type": "activity"
        },
        "$ELSE": {
          "label": "else",
          "optional": true,
          "type": "list_empty_do"
        },
        "$LIST": {
          "label": "across",
          "type": "assignment"
        }
      },
      "roles": "CZSZKZKZSZKZKT",
      "slots": [
        "execute"
      ],
      "tokens": [
        "repeating",
        " ",
        "across",
        ": ",
        "$LIST",
        ", ",
        "$AS",
        ", ",
        "do",
        ": ",
        "$DO",
        ", ",
        "$ELSE",
        "."
      ]
    }
  },
  {
    "desc": "ElseIfEmpty: Runs an activity when a list is empty.",
    "group": [
      "list"
    ],
    "name": "list_empty_do",
    "uses": "flow",
    "with": {
      "params": {
        "$DO": {
          "label": "do",
          "type": "activity"
        }
      },
      "roles": "SZK",
      "tokens": [
        "elseIfEmptyDo",
        ": ",
        "$DO"
      ]
    }
  },
  {
    "desc": "Gather list: Transform the values from a list.\n\t\tThe named pattern gets called once for each value in the list.\n\t\tIt get called with two parameters: 'in' as each value from the list, \n\t\tand 'out' as the var passed to the gather.",
    "group": [
      "list"
    ],
    "name": "list_gather",
    "uses": "flow",
    "with": {
      "params": {
        "$FROM": {
          "label": "from",
          "type": "list_source"
        },
        "$USING": {
          "label": "using",
          "type": "text"
        },
        "$VAR": {
          "label": "var",
          "type": "variable_name"
        }
      },
      "roles": "CZKZKZSZKT",
      "tokens": [
        "gather",
        ": ",
        "$VAR",
        ", ",
        "$FROM",
        ", ",
        "using",
        ": ",
        "$USING",
        "."
      ]
    }
  },
  {
    "desc": "Length of List: Determines the number of values in a list.",
    "group": [
      "list"
    ],
    "name": "list_len",
    "spec": "length of {list:assignment}",
    "uses": "flow",
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
    "uses": "flow",
    "with": {
      "params": {
        "$FROM_LIST": {
          "label": "fromList",
          "type": "assignment"
        },
        "$TO_LIST": {
          "label": "toList",
          "type": "text"
        },
        "$USING_PATTERN": {
          "label": "usingPattern",
          "type": "text"
        }
      },
      "roles": "QZSZKZSZKZSZK",
      "slots": [
        "execute"
      ],
      "tokens": [
        "map",
        " ",
        "toList",
        ": ",
        "$TO_LIST",
        ", ",
        "fromList",
        ": ",
        "$FROM_LIST",
        ", ",
        "usingPattern",
        ": ",
        "$USING_PATTERN"
      ]
    }
  },
  {
    "desc": "Reduce List: Transform the values from one list by combining them into a single value.\n\t\tThe named pattern is called with two parameters: 'in' ( each element of the list ) and 'out' ( ex. a record ).",
    "group": [
      "list"
    ],
    "name": "list_reduce",
    "uses": "flow",
    "with": {
      "params": {
        "$FROM_LIST": {
          "label": "fromList",
          "type": "assignment"
        },
        "$INTO_VALUE": {
          "label": "intoValue",
          "type": "text"
        },
        "$USING_PATTERN": {
          "label": "usingPattern",
          "type": "text"
        }
      },
      "roles": "QZSZKZSZKZSZK",
      "slots": [
        "execute"
      ],
      "tokens": [
        "reduce",
        " ",
        "intoValue",
        ": ",
        "$INTO_VALUE",
        ", ",
        "fromList",
        ": ",
        "$FROM_LIST",
        ", ",
        "usingPattern",
        ": ",
        "$USING_PATTERN"
      ]
    }
  },
  {
    "desc": "Set Value of List: Overwrite an existing value in a list.",
    "group": [
      "list"
    ],
    "name": "list_set",
    "uses": "flow",
    "with": {
      "params": {
        "$FROM": {
          "label": "from",
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
      "roles": "QZSZKZSZKZSZK",
      "slots": [
        "execute"
      ],
      "tokens": [
        "set",
        " ",
        "list",
        ": ",
        "$LIST",
        ", ",
        "index",
        ": ",
        "$INDEX",
        ", ",
        "from",
        ": ",
        "$FROM"
      ]
    }
  },
  {
    "desc": "Slice of List: Create a new list from a section of another list.",
    "group": [
      "list"
    ],
    "name": "list_slice",
    "spec": "slice {list:assignment} {from entry%start?number} {ending before entry%end?number}",
    "uses": "flow",
    "with": {
      "slots": [
        "num_list_eval",
        "text_list_eval",
        "record_list_eval"
      ]
    }
  },
  {
    "desc": "Sort numbers: .",
    "group": [
      "list"
    ],
    "name": "list_sort_by_field",
    "uses": "flow",
    "with": {
      "params": {
        "$NAME": {
          "label": "name",
          "type": "text"
        }
      },
      "roles": "SZK",
      "tokens": [
        "byField",
        ": ",
        "$NAME"
      ]
    }
  },
  {
    "desc": "Sort numbers: .",
    "group": [
      "list"
    ],
    "name": "list_sort_numbers",
    "uses": "flow",
    "with": {
      "params": {
        "$BY_FIELD": {
          "label": "byField",
          "optional": true,
          "type": "list_sort_by_field"
        },
        "$ORDER": {
          "label": "order",
          "type": "list_order"
        },
        "$VAR": {
          "label": "numbers",
          "type": "variable_name"
        }
      },
      "roles": "CZSZKZKZKT",
      "tokens": [
        "sort",
        " ",
        "numbers",
        ": ",
        "$VAR",
        ", ",
        "$BY_FIELD",
        ", ",
        "$ORDER",
        "."
      ]
    }
  },
  {
    "desc": "Sort list: rearrange the elements in the named list by using the designated pattern to test pairs of elements.",
    "group": [
      "list"
    ],
    "name": "list_sort_text",
    "uses": "flow",
    "with": {
      "params": {
        "$BY_FIELD": {
          "label": "byField",
          "optional": true,
          "type": "list_sort_by_field"
        },
        "$CASE": {
          "label": "case",
          "type": "list_case"
        },
        "$ORDER": {
          "label": "order",
          "type": "list_order"
        },
        "$VAR": {
          "label": "text",
          "type": "variable_name"
        }
      },
      "roles": "CZSZKZKZKZKT",
      "tokens": [
        "sort",
        " ",
        "text",
        ": ",
        "$VAR",
        ", ",
        "$BY_FIELD",
        ", ",
        "$ORDER",
        ", ",
        "$CASE",
        "."
      ]
    }
  },
  {
    "desc": "Sort list: rearrange the elements in the named list by using the designated pattern to test pairs of elements.",
    "group": [
      "list"
    ],
    "name": "list_sort_using",
    "uses": "flow",
    "with": {
      "params": {
        "$USING": {
          "label": "using",
          "type": "text"
        },
        "$VAR": {
          "label": "records",
          "type": "variable_name"
        }
      },
      "roles": "CZSZKZSZKT",
      "tokens": [
        "sort",
        " ",
        "records",
        ": ",
        "$VAR",
        ", ",
        "using",
        ": ",
        "$USING",
        "."
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
    "uses": "flow",
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
    "uses": "flow",
    "with": {
      "params": {
        "$ARGUMENTS": {
          "label": "arguments",
          "optional": true,
          "type": "arguments"
        },
        "$NAME": {
          "label": "name",
          "type": "text"
        }
      },
      "roles": "QZSZKZSZK",
      "slots": [
        "record_eval"
      ],
      "tokens": [
        "make",
        " ",
        "name",
        ": ",
        "$NAME",
        ", ",
        "arguments",
        ": ",
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
    "uses": "flow",
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
    "uses": "flow",
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
    "uses": "flow",
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
    "uses": "flow",
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
    "uses": "flow",
    "with": {
      "slots": [
        "text_eval"
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
    "uses": "flow",
    "with": {
      "slots": [
        "bool_eval"
      ]
    }
  },
  {
    "desc": "Name Of: Full name of the object.",
    "group": [
      "objects"
    ],
    "name": "name_of",
    "spec": "name of {object:object_eval}",
    "uses": "flow",
    "with": {
      "slots": [
        "text_eval"
      ]
    }
  },
  {
    "desc": "Next: In a repeating loop, try the next iteration of the loop.",
    "group": [
      "flow"
    ],
    "name": "next",
    "uses": "flow",
    "with": {
      "params": {},
      "roles": "E",
      "slots": [
        "execute"
      ],
      "tokens": [
        "next"
      ]
    }
  },
  {
    "desc": "Is Not: Returns the opposite value.",
    "group": [
      "logic"
    ],
    "name": "not",
    "uses": "flow",
    "with": {
      "params": {
        "$TEST": {
          "label": "test",
          "type": "bool_eval"
        }
      },
      "roles": "FZK",
      "slots": [
        "bool_eval"
      ],
      "tokens": [
        "not",
        ": ",
        "$TEST"
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
    "uses": "flow",
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
    "uses": "flow",
    "with": {
      "params": {
        "$VALUES": {
          "label": "values",
          "repeats": true,
          "type": "number"
        }
      },
      "roles": "QZSZK",
      "slots": [
        "num_list_eval"
      ],
      "tokens": [
        "numbers",
        " ",
        "values",
        ": ",
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
    "uses": "flow",
    "with": {
      "slots": [
        "bool_eval"
      ]
    }
  },
  {
    "desc": "Object Name: Returns a noun's object id.",
    "group": [
      "objects"
    ],
    "name": "object_name",
    "spec": "object named {name:text_eval}",
    "uses": "flow",
    "with": {
      "slots": [
        "object_eval"
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
    "uses": "flow",
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
    "uses": "flow",
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
    "uses": "flow",
    "with": {
      "params": {
        "$NUM": {
          "label": "num",
          "type": "number_eval"
        }
      },
      "roles": "QZSZK",
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "printNumWord",
        " ",
        "num",
        ": ",
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
    "uses": "flow",
    "with": {
      "slots": [
        "number_eval"
      ]
    }
  },
  {
    "desc": "Put: put a value into the field of an record or object",
    "name": "put_at_field",
    "uses": "flow",
    "with": {
      "params": {
        "$AT_FIELD": {
          "label": "atField",
          "type": "text"
        },
        "$FROM": {
          "label": "from",
          "type": "assignment"
        },
        "$INTO": {
          "label": "into",
          "type": "fields"
        }
      },
      "roles": "CZKZKZSZKT",
      "slots": [
        "execute"
      ],
      "tokens": [
        "put",
        ": ",
        "$FROM",
        ", ",
        "$INTO",
        ", ",
        "atField",
        ": ",
        "$AT_FIELD",
        "."
      ]
    }
  },
  {
    "desc": "Put: add a value to a list",
    "name": "put_edge",
    "uses": "flow",
    "with": {
      "params": {
        "$AT_EDGE": {
          "label": "atEdge",
          "type": "list_edge"
        },
        "$FROM": {
          "label": "from",
          "type": "assignment"
        },
        "$INTO": {
          "label": "into",
          "type": "list_target"
        }
      },
      "roles": "CZKZKZKT",
      "slots": [
        "execute"
      ],
      "tokens": [
        "put",
        ": ",
        "$FROM",
        ", ",
        "$INTO",
        ", ",
        "$AT_EDGE",
        "."
      ]
    }
  },
  {
    "desc": "Put: replace one value in a list with another",
    "name": "put_index",
    "uses": "flow",
    "with": {
      "params": {
        "$AT_INDEX": {
          "label": "atIndex",
          "type": "number_eval"
        },
        "$FROM": {
          "label": "from",
          "type": "assignment"
        },
        "$INTO": {
          "label": "into",
          "type": "list_target"
        }
      },
      "roles": "CZKZKZSZKT",
      "slots": [
        "execute"
      ],
      "tokens": [
        "put",
        ": ",
        "$FROM",
        ", ",
        "$INTO",
        ", ",
        "atIndex",
        ": ",
        "$AT_INDEX",
        "."
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
    "uses": "flow",
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
    "name": "range",
    "uses": "flow",
    "with": {
      "params": {
        "$BY_STEP": {
          "label": "byStep",
          "optional": true,
          "type": "number_eval"
        },
        "$FROM": {
          "label": "from",
          "optional": true,
          "type": "number_eval"
        },
        "$TO": {
          "label": "to",
          "type": "number_eval"
        }
      },
      "roles": "FZKZSZKZSZK",
      "slots": [
        "num_list_eval"
      ],
      "tokens": [
        "range",
        ": ",
        "$TO",
        ", ",
        "from",
        ": ",
        "$FROM",
        ", ",
        "byStep",
        ": ",
        "$BY_STEP"
      ]
    }
  },
  {
    "desc": "ReciprocalOf: Returns the implied relative of a noun (ex. the source in a one-to-many relation.)",
    "group": [
      "relations"
    ],
    "name": "reciprocal_of",
    "uses": "flow",
    "with": {
      "params": {
        "$OBJ": {
          "label": "of",
          "type": "text_eval"
        },
        "$VIA": {
          "label": "via",
          "type": "relation_name"
        }
      },
      "roles": "FZKZSZK",
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "reciprocal",
        ": ",
        "$VIA",
        ", ",
        "of",
        ": ",
        "$OBJ"
      ]
    }
  },
  {
    "desc": "ReciprocalsOf: Returns the implied relative of a noun (ex. the sources of a many-to-many relation.)",
    "group": [
      "relations"
    ],
    "name": "reciprocals_of",
    "uses": "flow",
    "with": {
      "params": {
        "$OBJ": {
          "label": "of",
          "type": "text_eval"
        },
        "$VIA": {
          "label": "via",
          "type": "relation_name"
        }
      },
      "roles": "FZKZSZK",
      "slots": [
        "text_list_eval"
      ],
      "tokens": [
        "reciprocals",
        ": ",
        "$VIA",
        ", ",
        "of",
        ": ",
        "$OBJ"
      ]
    }
  },
  {
    "desc": "Relate: Relate two nouns.",
    "group": [
      "relations"
    ],
    "name": "relate",
    "uses": "flow",
    "with": {
      "params": {
        "$OBJ": {
          "label": "obj",
          "type": "text_eval"
        },
        "$TO_OBJ": {
          "label": "to",
          "type": "text_eval"
        },
        "$VIA": {
          "label": "via",
          "type": "relation_name"
        }
      },
      "roles": "CZKZSZKZSZKT",
      "slots": [
        "execute"
      ],
      "tokens": [
        "relate",
        ": ",
        "$OBJ",
        ", ",
        "to",
        ": ",
        "$TO_OBJ",
        ", ",
        "via",
        ": ",
        "$VIA",
        "."
      ]
    }
  },
  {
    "desc": "RelativeOf: Returns the relative of a noun (ex. the target of a one-to-one relation.)",
    "group": [
      "relations"
    ],
    "name": "relative_of",
    "uses": "flow",
    "with": {
      "params": {
        "$OBJ": {
          "label": "of",
          "type": "text_eval"
        },
        "$VIA": {
          "label": "via",
          "type": "relation_name"
        }
      },
      "roles": "FZKZSZK",
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "relative",
        ": ",
        "$VIA",
        ", ",
        "of",
        ": ",
        "$OBJ"
      ]
    }
  },
  {
    "desc": "RelativesOf: Returns the relatives of a noun as a list of names (ex. the targets of one-to-many relation).",
    "group": [
      "relations"
    ],
    "name": "relatives_of",
    "uses": "flow",
    "with": {
      "params": {
        "$OBJ": {
          "label": "of",
          "type": "text_eval"
        },
        "$VIA": {
          "label": "via",
          "type": "relation_name"
        }
      },
      "roles": "FZKZSZK",
      "slots": [
        "text_list_eval"
      ],
      "tokens": [
        "relatives",
        ": ",
        "$VIA",
        ", ",
        "of",
        ": ",
        "$OBJ"
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
    "uses": "flow",
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
    "uses": "flow",
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
    "uses": "flow",
    "with": {
      "params": {
        "$TEXT": {
          "label": "text",
          "type": "text_eval"
        }
      },
      "roles": "CZKT",
      "slots": [
        "execute"
      ],
      "tokens": [
        "say",
        ": ",
        "$TEXT",
        "."
      ]
    }
  },
  {
    "desc": "Shuffle Text: When called multiple times returns its inputs at random.",
    "group": [
      "format"
    ],
    "name": "shuffle_text",
    "uses": "flow",
    "with": {
      "params": {
        "$PARTS": {
          "label": "parts",
          "repeats": true,
          "type": "text_eval"
        }
      },
      "roles": "QZSZK",
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "shuffleText",
        " ",
        "parts",
        ": ",
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
    "uses": "flow",
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
    "uses": "flow",
    "with": {
      "params": {
        "$GO": {
          "label": "go",
          "optional": true,
          "type": "activity"
        }
      },
      "roles": "QZSZK",
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "slash",
        " ",
        "go",
        ": ",
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
    "uses": "flow",
    "with": {
      "params": {
        "$GO": {
          "label": "go",
          "optional": true,
          "type": "activity"
        }
      },
      "roles": "QZSZK",
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "span",
        " ",
        "go",
        ": ",
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
    "uses": "flow",
    "with": {
      "params": {
        "$PARTS": {
          "label": "parts",
          "repeats": true,
          "type": "text_eval"
        }
      },
      "roles": "QZSZK",
      "slots": [
        "text_eval"
      ],
      "tokens": [
        "stoppingText",
        " ",
        "parts",
        ": ",
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
    "uses": "flow",
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
    "uses": "flow",
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
    "spec": "text {values*text|comma-and}",
    "uses": "flow",
    "with": {
      "slots": [
        "text_list_eval"
      ]
    }
  },
  {
    "desc": "Not Equal To: Two values don't match exactly.",
    "group": [
      "comparison"
    ],
    "name": "unequal",
    "spec": "\u003c\u003e",
    "uses": "flow",
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
    "spec": "unpack {field:text} from {record:record_eval}",
    "uses": "flow",
    "with": {
      "slots": [
        "bool_eval",
        "number_eval",
        "text_eval",
        "record_eval",
        "num_list_eval",
        "text_list_eval",
        "record_list_eval",
        "assignment"
      ]
    }
  },
  {
    "desc": "While: Repeat a series of statements while a conditional is true.",
    "group": [
      "flow"
    ],
    "name": "while",
    "uses": "flow",
    "with": {
      "params": {
        "$DO": {
          "label": "do",
          "type": "activity"
        },
        "$TRUE": {
          "label": "while",
          "type": "bool_eval"
        }
      },
      "roles": "CZSZKZSZKT",
      "slots": [
        "execute"
      ],
      "tokens": [
        "repeating",
        " ",
        "while",
        ": ",
        "$TRUE",
        ", ",
        "do",
        ": ",
        "$DO",
        "."
      ]
    }
  }
];
const stub = [
  "text_value",
  "cycle_text",
  "shuffle_text",
  "stopping_text",
  "arguments",
  "argument",
  "debug_level",
  "render_template",
  "determine_act",
  "determine_num",
  "determine_text",
  "determine_bool",
  "determine_num_list",
  "determine_text_list",
  "list_case",
  "list_edge",
  "list_order"
];
