function getStory() {
  return {
  "id": "id-174cda99dde-0",
  "type": "story",
  "value": {
    "$PARAGRAPH": [
      {
        "id": "id-174cda99dde-5",
        "type": "paragraph",
        "value": {
          "$STORY_STATEMENT": [
            {
              "id": "id-174cda99dde-4",
              "type": "story_statement",
              "value": {
                "id": "id-174cda99dde-1",
                "type": "pattern_decl",
                "value": {
                  "$NAME": {
                    "id": "id-174cda99dde-2",
                    "type": "pattern_name",
                    "value": "starts with vowel"
                  },
                  "$TYPE": {
                    "id": "id-174cda99dde-3",
                    "type": "pattern_type",
                    "value": {
                      "$VALUE": {
                        "id": "id-174cda99dde-6",
                        "type": "variable_type",
                        "value": {
                          "$PRIMITIVE": {
                            "id": "id-174cda99dde-7",
                            "type": "primitive_type",
                            "value": "$BOOL"
                          }
                        }
                      }
                    }
                  },
                  "$OPTVARS": {
                    "id": "id-174cda99dde-8",
                    "type": "pattern_variables_tail",
                    "value": {
                      "$VARIABLE_DECL": [
                        {
                          "id": "id-174cda99dde-9",
                          "type": "variable_decl",
                          "value": {
                            "$TYPE": {
                              "id": "id-174cda99dde-10",
                              "type": "variable_type",
                              "value": {
                                "$PRIMITIVE": {
                                  "id": "id-174cda99dde-12",
                                  "type": "primitive_type",
                                  "value": "$TEXT"
                                }
                              }
                            },
                            "$NAME": {
                              "id": "id-174cda99dde-11",
                              "type": "variable_name",
                              "value": "text"
                            }
                          }
                        }
                      ]
                    }
                  }
                }
              }
            },
            {
              "id": "id-174cdb223ae-27",
              "type": "story_statement",
              "value": {
                "id": "id-174cdb223ae-25",
                "type": "comment",
                "value": {
                  "$LINES": {
                    "id": "id-174cdb223ae-26",
                    "type": "lines",
                    "value": "Determine if text starts with a vowel or vowel sound."
                  }
                }
              }
            }
          ]
        }
      },
      {
        "id": "id-174cda99dde-17",
        "type": "paragraph",
        "value": {
          "$STORY_STATEMENT": [
            {
              "id": "id-174cda99dde-16",
              "type": "story_statement",
              "value": {
                "id": "id-174cda99dde-13",
                "type": "pattern_actions",
                "value": {
                  "$NAME": {
                    "id": "id-174cda99dde-14",
                    "type": "pattern_name",
                    "value": "starts with vowel"
                  },
                  "$PATTERN_RULES": {
                    "id": "id-174cda99dde-15",
                    "type": "pattern_rules",
                    "value": {
                      "$PATTERN_RULE": [
                        {
                          "id": "id-174cda99dde-18",
                          "type": "pattern_rule",
                          "value": {
                            "$GUARD": {
                              "id": "id-174cda99dde-19",
                              "type": "bool_eval",
                              "value": {
                                "id": "id-174cdb223ae-0",
                                "type": "matches",
                                "value": {
                                  "$TEXT": {
                                    "id": "id-174cdb223ae-1",
                                    "type": "text_eval",
                                    "value": {
                                      "id": "id-174cdb223ae-3",
                                      "type": "get_var",
                                      "value": {
                                        "$NAME": {
                                          "id": "id-174cdb223ae-4",
                                          "type": "text_eval",
                                          "value": {
                                            "id": "id-174cdb223ae-5",
                                            "type": "text_value",
                                            "value": {
                                              "$TEXT": {
                                                "id": "id-174cdb223ae-6",
                                                "type": "text",
                                                "value": "text"
                                              }
                                            }
                                          }
                                        }
                                      }
                                    }
                                  },
                                  "$PATTERN": {
                                    "id": "id-174cdb223ae-2",
                                    "type": "text",
                                    "value": "^(?i:A|E|I|O|U)"
                                  }
                                }
                              }
                            },
                            "$HOOK": {
                              "id": "id-174cda99dde-20",
                              "type": "program_hook",
                              "value": {
                                "$RESULT": {
                                  "id": "id-174cdb223ae-7",
                                  "type": "program_return",
                                  "value": {
                                    "$RESULT": {
                                      "id": "id-174cdb223ae-8",
                                      "type": "program_result",
                                      "value": {
                                        "$PRIMITIVE": {
                                          "id": "id-174cdb223ae-9",
                                          "type": "primitive_func",
                                          "value": {
                                            "$BOOL_EVAL": {
                                              "id": "id-174cdb223ae-10",
                                              "type": "bool_eval",
                                              "value": {
                                                "id": "id-174cdb223ae-11",
                                                "type": "bool_value",
                                                "value": {
                                                  "$BOOL": {
                                                    "id": "id-174cdb223ae-12",
                                                    "type": "bool",
                                                    "value": "$TRUE"
                                                  }
                                                }
                                              }
                                            }
                                          }
                                        }
                                      }
                                    }
                                  }
                                }
                              }
                            }
                          }
                        },
                        {
                          "id": "id-174cdb223ae-13",
                          "type": "pattern_rule",
                          "value": {
                            "$GUARD": {
                              "id": "id-174cdb223ae-14",
                              "type": "bool_eval",
                              "value": {
                                "id": "id-174cdb223ae-16",
                                "type": "matches",
                                "value": {
                                  "$TEXT": {
                                    "id": "id-174cdb223ae-17",
                                    "type": "text_eval",
                                    "value": {
                                      "id": "id-174cdb223ae-31",
                                      "type": "get_var",
                                      "value": {
                                        "$NAME": {
                                          "id": "id-174cdb223ae-32",
                                          "type": "text_eval",
                                          "value": {
                                            "id": "id-174cdb223ae-33",
                                            "type": "text_value",
                                            "value": {
                                              "$TEXT": {
                                                "id": "id-174cdb223ae-34",
                                                "type": "text",
                                                "value": "text"
                                              }
                                            }
                                          }
                                        }
                                      }
                                    }
                                  },
                                  "$PATTERN": {
                                    "id": "id-174cdb223ae-18",
                                    "type": "text",
                                    "value": "^(?i:EU|EW|ONCE|ONE|OUI|UBI|UGAND|UKRAIN|UKULELE|ULYSS|UNA|UNESCO|UNI|UNUM|URA|URE|URI|URO|URU|USA|USE|USI|USU|UTA|UTE|UTI|UTO)"
                                  }
                                }
                              }
                            },
                            "$HOOK": {
                              "id": "id-174cdb223ae-15",
                              "type": "program_hook",
                              "value": {
                                "$RESULT": {
                                  "id": "id-174cdb223ae-19",
                                  "type": "program_return",
                                  "value": {
                                    "$RESULT": {
                                      "id": "id-174cdb223ae-20",
                                      "type": "program_result",
                                      "value": {
                                        "$PRIMITIVE": {
                                          "id": "id-174cdb223ae-21",
                                          "type": "primitive_func",
                                          "value": {
                                            "$BOOL_EVAL": {
                                              "id": "id-174cdb223ae-22",
                                              "type": "bool_eval",
                                              "value": {
                                                "id": "id-174cdb223ae-23",
                                                "type": "bool_value",
                                                "value": {
                                                  "$BOOL": {
                                                    "id": "id-174cdb223ae-24",
                                                    "type": "bool",
                                                    "value": "$FALSE"
                                                  }
                                                }
                                              }
                                            }
                                          }
                                        }
                                      }
                                    }
                                  }
                                }
                              }
                            }
                          }
                        },
                        {
                          "id": "id-174cdb223ae-28",
                          "type": "pattern_rule",
                          "value": {
                            "$GUARD": {
                              "id": "id-174cdb223ae-29",
                              "type": "bool_eval",
                              "value": {
                                "id": "id-174cdb223ae-35",
                                "type": "matches",
                                "value": {
                                  "$TEXT": {
                                    "id": "id-174cdb223ae-36",
                                    "type": "text_eval",
                                    "value": {
                                      "id": "id-174cdb223ae-38",
                                      "type": "get_var",
                                      "value": {
                                        "$NAME": {
                                          "id": "id-174cdb223ae-39",
                                          "type": "text_eval",
                                          "value": {
                                            "id": "id-174cdb223ae-40",
                                            "type": "text_value",
                                            "value": {
                                              "$TEXT": {
                                                "id": "id-174cdb223ae-41",
                                                "type": "text",
                                                "value": "text"
                                              }
                                            }
                                          }
                                        }
                                      }
                                    }
                                  },
                                  "$PATTERN": {
                                    "id": "id-174cdb223ae-37",
                                    "type": "text",
                                    "value": "^(?i:HEIR|HERB|HOMAGE|HONEST|HONOR|HONOUR|HORS|HOUR)"
                                  }
                                }
                              }
                            },
                            "$HOOK": {
                              "id": "id-174cdb223ae-30",
                              "type": "program_hook",
                              "value": {
                                "$RESULT": {
                                  "id": "id-174cdb223ae-42",
                                  "type": "program_return",
                                  "value": {
                                    "$RESULT": {
                                      "id": "id-174cdb223ae-43",
                                      "type": "program_result",
                                      "value": {
                                        "$PRIMITIVE": {
                                          "id": "id-174cdb223ae-44",
                                          "type": "primitive_func",
                                          "value": {
                                            "$BOOL_EVAL": {
                                              "id": "id-174cdb223ae-45",
                                              "type": "bool_eval",
                                              "value": {
                                                "id": "id-174cdb223ae-48",
                                                "type": "bool_value",
                                                "value": {
                                                  "$BOOL": {
                                                    "id": "id-174cdb223ae-49",
                                                    "type": "bool",
                                                    "value": "$TRUE"
                                                  }
                                                }
                                              }
                                            }
                                          }
                                        }
                                      }
                                    }
                                  }
                                }
                              }
                            }
                          }
                        }
                      ]
                    }
                  }
                }
              }
            }
          ]
        }
      },
      {
        "id": "id-174cdb223ae-59",
        "type": "paragraph",
        "value": {
          "$STORY_STATEMENT": [
            {
              "id": "id-174cdb223ae-58",
              "type": "story_statement",
              "value": {
                "id": "id-174cdb223ae-55",
                "type": "test_statement",
                "value": {
                  "$TEST_NAME": {
                    "id": "id-174cdb223ae-56",
                    "type": "test_name",
                    "value": "vowels"
                  },
                  "$TEST": {
                    "id": "id-174cdb223ae-57",
                    "type": "testing",
                    "value": {
                      "id": "id-174cdb223ae-60",
                      "type": "test_output",
                      "value": {
                        "$LINES": {
                          "id": "id-174cdb223ae-61",
                          "type": "lines",
                          "value": "ok"
                        }
                      }
                    }
                  }
                }
              }
            }
          ]
        }
      },
      {
        "id": "id-174cdb223ae-66",
        "type": "paragraph",
        "value": {
          "$STORY_STATEMENT": [
            {
              "id": "id-174cdb223ae-65",
              "type": "story_statement",
              "value": {
                "id": "id-174cdb223ae-62",
                "type": "test_rule",
                "value": {
                  "$TEST_NAME": {
                    "id": "id-174cdb223ae-63",
                    "type": "test_name",
                    "value": "vowels"
                  },
                  "$HOOK": {
                    "id": "id-174cdb223ae-64",
                    "type": "program_hook",
                    "value": {
                      "$ACTIVITY": {
                        "id": "id-174cdb223ae-67",
                        "type": "activity",
                        "value": {
                          "$EXE": [
                            {
                              "id": "id-174cdb223ae-68",
                              "type": "execute",
                              "value": {
                                "id": "id-174cdb223ae-80",
                                "type": "say_text",
                                "value": {
                                  "$TEXT": {
                                    "id": "id-174cdb223ae-81",
                                    "type": "text_eval",
                                    "value": {
                                      "id": "id-174cdb223ae-84",
                                      "type": "choose_text",
                                      "value": {
                                        "$IF": {
                                          "id": "id-174cdb223ae-85",
                                          "type": "bool_eval",
                                          "value": {
                                            "id": "id-174cdd39e3a-0",
                                            "type": "determine_bool",
                                            "value": {
                                              "$NAME": {
                                                "id": "id-174cdd39e3a-1",
                                                "type": "pattern_name",
                                                "value": "starts with vowel"
                                              },
                                              "$ARGUMENTS": {
                                                "id": "id-174cdd39e3a-2",
                                                "type": "arguments",
                                                "value": {
                                                  "$ARGS": [
                                                    {
                                                      "id": "id-174cdd39e3a-3",
                                                      "type": "argument",
                                                      "value": {
                                                        "$NAME": {
                                                          "id": "id-174cdd39e3a-4",
                                                          "type": "variable_name",
                                                          "value": "text"
                                                        },
                                                        "$FROM": {
                                                          "id": "id-174cdd39e3a-5",
                                                          "type": "assignment",
                                                          "value": {
                                                            "id": "id-174cdd39e3a-6",
                                                            "type": "assign_text",
                                                            "value": {
                                                              "$VAL": {
                                                                "id": "id-174cdd39e3a-7",
                                                                "type": "text_eval",
                                                                "value": {
                                                                  "id": "id-174cdd39e3a-8",
                                                                  "type": "text_value",
                                                                  "value": {
                                                                    "$TEXT": {
                                                                      "id": "id-174cdd39e3a-9",
                                                                      "type": "text",
                                                                      "value": "house"
                                                                    }
                                                                  }
                                                                }
                                                              }
                                                            }
                                                          }
                                                        }
                                                      }
                                                    }
                                                  ]
                                                }
                                              }
                                            }
                                          }
                                        },
                                        "$TRUE": {
                                          "id": "id-174cdb223ae-86",
                                          "type": "text_eval",
                                          "value": {
                                            "id": "id-174cdd39e3a-10",
                                            "type": "text_value",
                                            "value": {
                                              "$TEXT": {
                                                "id": "id-174cdd39e3a-11",
                                                "type": "text",
                                                "value": "error: house doesn't start with a vowel"
                                              }
                                            }
                                          }
                                        },
                                        "$FALSE": {
                                          "id": "id-174cdb223ae-87",
                                          "type": "text_eval",
                                          "value": {
                                            "id": "id-174cdd39e3a-12",
                                            "type": "text_value",
                                            "value": {
                                              "$TEXT": {
                                                "id": "id-174cdd39e3a-13",
                                                "type": "text",
                                                "value": "$EMPTY"
                                              }
                                            }
                                          }
                                        }
                                      }
                                    }
                                  }
                                }
                              }
                            },
                            {
                              "id": "id-174d0e35aff-0",
                              "type": "execute",
                              "value": {
                                "id": "id-174d0e35aff-1",
                                "type": "say_text",
                                "value": {
                                  "$TEXT": {
                                    "id": "id-174d0e35aff-2",
                                    "type": "text_eval",
                                    "value": {
                                      "id": "id-174d0e35aff-3",
                                      "type": "choose_text",
                                      "value": {
                                        "$IF": {
                                          "id": "id-174d0e35aff-4",
                                          "type": "bool_eval",
                                          "value": {
                                            "id": "id-174d0e35aff-5",
                                            "type": "determine_bool",
                                            "value": {
                                              "$NAME": {
                                                "id": "id-174d0e35aff-6",
                                                "type": "pattern_name",
                                                "value": "starts with vowel"
                                              },
                                              "$ARGUMENTS": {
                                                "id": "id-174d0e35aff-7",
                                                "type": "arguments",
                                                "value": {
                                                  "$ARGS": [
                                                    {
                                                      "id": "id-174d0e35aff-8",
                                                      "type": "argument",
                                                      "value": {
                                                        "$NAME": {
                                                          "id": "id-174d0e35aff-9",
                                                          "type": "variable_name",
                                                          "value": "text"
                                                        },
                                                        "$FROM": {
                                                          "id": "id-174d0e35aff-10",
                                                          "type": "assignment",
                                                          "value": {
                                                            "id": "id-174d0e35aff-11",
                                                            "type": "assign_text",
                                                            "value": {
                                                              "$VAL": {
                                                                "id": "id-174d0e35aff-12",
                                                                "type": "text_eval",
                                                                "value": {
                                                                  "id": "id-174d0e35aff-13",
                                                                  "type": "text_value",
                                                                  "value": {
                                                                    "$TEXT": {
                                                                      "id": "id-174d0e35aff-14",
                                                                      "type": "text",
                                                                      "value": "hour"
                                                                    }
                                                                  }
                                                                }
                                                              }
                                                            }
                                                          }
                                                        }
                                                      }
                                                    }
                                                  ]
                                                }
                                              }
                                            }
                                          }
                                        },
                                        "$TRUE": {
                                          "id": "id-174d0e35aff-15",
                                          "type": "text_eval",
                                          "value": {
                                            "id": "id-174d0e35aff-16",
                                            "type": "text_value",
                                            "value": {
                                              "$TEXT": {
                                                "id": "id-174d0e35aff-17",
                                                "type": "text",
                                                "value": "$EMPTY"
                                              }
                                            }
                                          }
                                        },
                                        "$FALSE": {
                                          "id": "id-174d0e35aff-18",
                                          "type": "text_eval",
                                          "value": {
                                            "id": "id-174d0e35aff-19",
                                            "type": "text_value",
                                            "value": {
                                              "$TEXT": {
                                                "id": "id-174d0e35aff-20",
                                                "type": "text",
                                                "value": "error: hour starts with a vowel sound"
                                              }
                                            }
                                          }
                                        }
                                      }
                                    }
                                  }
                                }
                              }
                            },
                            {
                              "id": "id-174d0e35aff-31",
                              "type": "execute",
                              "value": {
                                "id": "id-174d0e35aff-32",
                                "type": "say_text",
                                "value": {
                                  "$TEXT": {
                                    "id": "id-174d0e35aff-33",
                                    "type": "text_eval",
                                    "value": {
                                      "id": "id-174d0e35aff-34",
                                      "type": "choose_text",
                                      "value": {
                                        "$IF": {
                                          "id": "id-174d0e35aff-35",
                                          "type": "bool_eval",
                                          "value": {
                                            "id": "id-174d0e35aff-36",
                                            "type": "determine_bool",
                                            "value": {
                                              "$NAME": {
                                                "id": "id-174d0e35aff-37",
                                                "type": "pattern_name",
                                                "value": "starts with vowel"
                                              },
                                              "$ARGUMENTS": {
                                                "id": "id-174d0e35aff-38",
                                                "type": "arguments",
                                                "value": {
                                                  "$ARGS": [
                                                    {
                                                      "id": "id-174d0e35aff-39",
                                                      "type": "argument",
                                                      "value": {
                                                        "$NAME": {
                                                          "id": "id-174d0e35aff-40",
                                                          "type": "variable_name",
                                                          "value": "text"
                                                        },
                                                        "$FROM": {
                                                          "id": "id-174d0e35aff-41",
                                                          "type": "assignment",
                                                          "value": {
                                                            "id": "id-174d0e35aff-42",
                                                            "type": "assign_text",
                                                            "value": {
                                                              "$VAL": {
                                                                "id": "id-174d0e35aff-43",
                                                                "type": "text_eval",
                                                                "value": {
                                                                  "id": "id-174d0e35aff-44",
                                                                  "type": "text_value",
                                                                  "value": {
                                                                    "$TEXT": {
                                                                      "id": "id-174d0e35aff-45",
                                                                      "type": "text",
                                                                      "value": "beta"
                                                                    }
                                                                  }
                                                                }
                                                              }
                                                            }
                                                          }
                                                        }
                                                      }
                                                    }
                                                  ]
                                                }
                                              }
                                            }
                                          }
                                        },
                                        "$TRUE": {
                                          "id": "id-174d0e35aff-46",
                                          "type": "text_eval",
                                          "value": {
                                            "id": "id-174d0e35aff-47",
                                            "type": "text_value",
                                            "value": {
                                              "$TEXT": {
                                                "id": "id-174d0e35aff-48",
                                                "type": "text",
                                                "value": "error: beta doesn't start with a vowel"
                                              }
                                            }
                                          }
                                        },
                                        "$FALSE": {
                                          "id": "id-174d0e35aff-49",
                                          "type": "text_eval",
                                          "value": {
                                            "id": "id-174d0e35aff-50",
                                            "type": "text_value",
                                            "value": {
                                              "$TEXT": {
                                                "id": "id-174d0e35aff-51",
                                                "type": "text",
                                                "value": "$EMPTY"
                                              }
                                            }
                                          }
                                        }
                                      }
                                    }
                                  }
                                }
                              }
                            },
                            {
                              "id": "id-174d0e35aff-52",
                              "type": "execute",
                              "value": {
                                "id": "id-174d0e35aff-53",
                                "type": "say_text",
                                "value": {
                                  "$TEXT": {
                                    "id": "id-174d0e35aff-54",
                                    "type": "text_eval",
                                    "value": {
                                      "id": "id-174d0e35aff-55",
                                      "type": "choose_text",
                                      "value": {
                                        "$IF": {
                                          "id": "id-174d0e35aff-56",
                                          "type": "bool_eval",
                                          "value": {
                                            "id": "id-174d0e35aff-57",
                                            "type": "determine_bool",
                                            "value": {
                                              "$NAME": {
                                                "id": "id-174d0e35aff-58",
                                                "type": "pattern_name",
                                                "value": "starts with vowel"
                                              },
                                              "$ARGUMENTS": {
                                                "id": "id-174d0e35aff-59",
                                                "type": "arguments",
                                                "value": {
                                                  "$ARGS": [
                                                    {
                                                      "id": "id-174d0e35aff-60",
                                                      "type": "argument",
                                                      "value": {
                                                        "$NAME": {
                                                          "id": "id-174d0e35aff-61",
                                                          "type": "variable_name",
                                                          "value": "text"
                                                        },
                                                        "$FROM": {
                                                          "id": "id-174d0e35aff-62",
                                                          "type": "assignment",
                                                          "value": {
                                                            "id": "id-174d0e35aff-63",
                                                            "type": "assign_text",
                                                            "value": {
                                                              "$VAL": {
                                                                "id": "id-174d0e35aff-64",
                                                                "type": "text_eval",
                                                                "value": {
                                                                  "id": "id-174d0e35aff-65",
                                                                  "type": "text_value",
                                                                  "value": {
                                                                    "$TEXT": {
                                                                      "id": "id-174d0e35aff-66",
                                                                      "type": "text",
                                                                      "value": "alpha"
                                                                    }
                                                                  }
                                                                }
                                                              }
                                                            }
                                                          }
                                                        }
                                                      }
                                                    }
                                                  ]
                                                }
                                              }
                                            }
                                          }
                                        },
                                        "$TRUE": {
                                          "id": "id-174d0e35aff-67",
                                          "type": "text_eval",
                                          "value": {
                                            "id": "id-174d0e35aff-68",
                                            "type": "text_value",
                                            "value": {
                                              "$TEXT": {
                                                "id": "id-174d0e35aff-69",
                                                "type": "text",
                                                "value": "$EMPTY"
                                              }
                                            }
                                          }
                                        },
                                        "$FALSE": {
                                          "id": "id-174d0e35aff-70",
                                          "type": "text_eval",
                                          "value": {
                                            "id": "id-174d0e35aff-71",
                                            "type": "text_value",
                                            "value": {
                                              "$TEXT": {
                                                "id": "id-174d0e35aff-72",
                                                "type": "text",
                                                "value": "alpha starts with a vowel"
                                              }
                                            }
                                          }
                                        }
                                      }
                                    }
                                  }
                                }
                              }
                            },
                            {
                              "id": "id-174d0e35aff-26",
                              "type": "execute",
                              "value": {
                                "id": "id-174d0e35aff-27",
                                "type": "say_text",
                                "value": {
                                  "$TEXT": {
                                    "id": "id-174d0e35aff-28",
                                    "type": "text_eval",
                                    "value": {
                                      "id": "id-174d0e35aff-29",
                                      "type": "text_value",
                                      "value": {
                                        "$TEXT": {
                                          "id": "id-174d0e35aff-30",
                                          "type": "text",
                                          "value": "ok"
                                        }
                                      }
                                    }
                                  }
                                }
                              }
                            }
                          ]
                        }
                      }
                    }
                  }
                }
              }
            }
          ]
        }
      }
    ]
  }
}
}

