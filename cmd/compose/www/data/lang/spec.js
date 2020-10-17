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
    "desc": "Object Reference: Helper used when referring to objects.",
    "name": "object_ref",
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
    "with": {}
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
    "desc": "Assign Variable: Assigns one variable or object to another.\n\t\tUsed internally for pattern calls in templates. ex. { pattern: .something }.",
    "group": [
      "variables"
    ],
    "name": "assign_var",
    "uses": "run",
    "with": {
      "params": {
        "$NAME": {
          "label": "name",
          "type": "text_eval"
        }
      },
      "slots": [
        "assignment"
      ],
      "tokens": [
        "from var",
        "$NAME"
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
    "spec": "the {field:text_eval} of {object%obj:object_ref}",
    "uses": "run",
    "with": {
      "slots": [
        "bool_eval",
        "number_eval",
        "text_eval",
        "num_list_eval",
        "text_list_eval"
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
        "num_list_eval",
        "text_list_eval",
        "object_ref"
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
    "spec": "{object%obj:object_ref} is {trait:text_eval}",
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
          "type": "object_ref"
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
    "spec": "Is {object%obj:object_ref} a kind of {kind:singular_kind}",
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
    "spec": "kind of {object%obj:object_ref}",
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
    "desc": "Name Of: Full name of the object.",
    "group": [
      "objects"
    ],
    "name": "name_of",
    "spec": "name of {object%obj:object_ref}",
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
    "desc": "Length of Number List: Determines the number of elements in a list of numbers.",
    "group": [
      "format"
    ],
    "name": "number_list_count",
    "uses": "run",
    "with": {
      "params": {
        "$ELEMS": {
          "label": "elems",
          "type": "num_list_eval"
        }
      },
      "slots": [
        "number_eval"
      ],
      "tokens": [
        "len of numbers",
        "$ELEMS"
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
    "desc": "Set Boolean Field: Sets the named field to the passed boolean value.",
    "group": [
      "objects"
    ],
    "name": "set_field_bool",
    "uses": "run",
    "with": {
      "params": {
        "$FIELD": {
          "label": "field",
          "type": "text_eval"
        },
        "$OBJ": {
          "label": "obj",
          "type": "object_ref"
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
        "set field bool",
        "$OBJ",
        "$FIELD",
        "$VAL"
      ]
    }
  },
  {
    "desc": "Set Number Field: Sets the named field to the passed number.",
    "group": [
      "objects"
    ],
    "name": "set_field_num",
    "uses": "run",
    "with": {
      "params": {
        "$FIELD": {
          "label": "field",
          "type": "text_eval"
        },
        "$OBJ": {
          "label": "obj",
          "type": "object_ref"
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
        "set field num",
        "$OBJ",
        "$FIELD",
        "$VAL"
      ]
    }
  },
  {
    "desc": "Set Number List Field: Sets the named field to the passed number list.",
    "group": [
      "objects"
    ],
    "name": "set_field_num_list",
    "uses": "run",
    "with": {
      "params": {
        "$FIELD": {
          "label": "field",
          "type": "text_eval"
        },
        "$OBJ": {
          "label": "obj",
          "type": "object_ref"
        },
        "$VALS": {
          "label": "vals",
          "type": "num_list_eval"
        }
      },
      "slots": [
        "execute"
      ],
      "tokens": [
        "set field num list",
        "$OBJ",
        "$FIELD",
        "$VALS"
      ]
    }
  },
  {
    "desc": "Set Text Field: Sets the named field to the passed text value.",
    "group": [
      "objects"
    ],
    "name": "set_field_text",
    "uses": "run",
    "with": {
      "params": {
        "$FIELD": {
          "label": "field",
          "type": "text_eval"
        },
        "$OBJ": {
          "label": "obj",
          "type": "object_ref"
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
        "set field text",
        "$OBJ",
        "$FIELD",
        "$VAL"
      ]
    }
  },
  {
    "desc": "Set Text List Field: Sets the named field to the passed text list.",
    "group": [
      "objects"
    ],
    "name": "set_field_text_list",
    "uses": "run",
    "with": {
      "params": {
        "$FIELD": {
          "label": "field",
          "type": "text_eval"
        },
        "$OBJ": {
          "label": "obj",
          "type": "object_ref"
        },
        "$VALS": {
          "label": "vals",
          "type": "text_list_eval"
        }
      },
      "slots": [
        "execute"
      ],
      "tokens": [
        "set field text list",
        "$OBJ",
        "$FIELD",
        "$VALS"
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
    "desc": "Length of Text List: Determines the number of text elements in a list.",
    "group": [
      "format"
    ],
    "name": "text_list_count",
    "uses": "run",
    "with": {
      "params": {
        "$ELEMS": {
          "label": "elems",
          "type": "text_list_eval"
        }
      },
      "slots": [
        "number_eval"
      ],
      "tokens": [
        "len of texts",
        "$ELEMS"
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
