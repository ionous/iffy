function getStory() {
  return {
  "id": "id-a-2",
  "type": "story",
  "value": {
    "$PARAGRAPH": [
      {
        "id": "id-a-1",
        "type": "paragraph",
        "value": {
          "$STORY_STATEMENT": [
            {
              "id": "id-a-0",
              "type": "story_statement",
              "value": {
                "id": "id-a-5",
                "type": "pattern_decl",
                "value": {
                  "$NAME": {
                    "id": "id-a-3",
                    "type": "pattern_name",
                    "value": "printName"
                  },
                  "$TYPE": {
                    "id": "id-a-4",
                    "type": "pattern_type",
                    "value": {
                      "$ACTIVITY": {
                        "id": "id-a-6",
                        "type": "patterned_activity",
                        "value": "$ACTIVITY"
                      }
                    }
                  },
                  "$OPTVARS": {
                    "id": "id-a-10",
                    "type": "pattern_variables_tail",
                    "value": {
                      "$VARIABLE_DECL": [
                        {
                          "id": "id-a-9",
                          "type": "variable_decl",
                          "value": {
                            "$TYPE": {
                              "id": "id-a-7",
                              "type": "variable_type",
                              "value": {
                                "$OBJECT": {
                                  "id": "id-a-13",
                                  "type": "object_type",
                                  "value": {
                                    "$AN": {
                                      "id": "id-a-11",
                                      "type": "an",
                                      "value": "$A"
                                    },
                                    "$KIND": {
                                      "id": "id-a-12",
                                      "type": "singular_kind",
                                      "value": "thing"
                                    }
                                  }
                                }
                              }
                            },
                            "$NAME": {
                              "id": "id-a-8",
                              "type": "variable_name",
                              "value": "target"
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
        "id": "id-a-15",
        "type": "paragraph",
        "value": {
          "$STORY_STATEMENT": [
            {
              "id": "id-a-14",
              "type": "story_statement",
              "value": {
                "id": "id-a-18",
                "type": "pattern_handler",
                "value": {
                  "$NAME": {
                    "id": "id-a-16",
                    "type": "pattern_name",
                    "value": "printName"
                  },
                  "$HOOK": {
                    "id": "id-a-17",
                    "type": "pattern_hook",
                    "value": {
                      "$ACTIVITY": {
                        "id": "id-a-20",
                        "type": "pattern_activity",
                        "value": {
                          "$GO": [

                            {
                              "id": "id-17311b2a4a8-0",
                              "type": "execute",
                              "value": {
                                "id": "id-17311b2a4a8-2",
                                "type": "say_text",
                                "value": {
                                  "$TEXT": {
                                    "id": "id-17311b2a4a8-1",
                                    "type": "text_eval",
                                    "value": {
                                      "id": "id-17311b2a4a8-4",
                                      "type": "singularize",
                                      "value": {
                                        "$TEXT": {
                                          "id": "id-17311b2a4a8-3",
                                          "type": "text_eval",
                                          "value": {
                                            "id": "id-17311b2a4a8-7",
                                            "type": "kind_of",
                                            "value": {
                                              "$OBJ": {
                                                "id": "id-17311b2a4a8-6",
                                                "type": "object_name",
                                                "value": {
                                                  "$NAME": {
                                                    "id": "id-17311b2a4a8-5",
                                                    "type": "text_eval",
                                                    "value": {
                                                      "id": "id-17311b2a4a8-9",
                                                      "type": "text_value",
                                                      "value": {
                                                        "$TEXT": {
                                                          "id": "id-17311b2a4a8-8",
                                                          "type": "text",
                                                          "value": "target"
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
                            }
                          ]
                        }
                      }
                    }
                  }
                }
              }
            },
            {
              "id": "id-a-21",
              "type": "story_statement",
              "value": {
                "id": "id-a-57",
                "type": "noun_statement",
                "value": {
                  "$LEDE": {
                    "id": "id-a-56",
                    "type": "lede",
                    "value": {
                      "$NOUN": [
                        {
                          "id": "id-a-54",
                          "type": "noun",
                          "value": {
                            "$COMMON_NOUN": {
                              "id": "id-a-60",
                              "type": "common_noun",
                              "value": {
                                "$DETERMINER": {
                                  "id": "id-a-58",
                                  "type": "determiner",
                                  "value": "$THE"
                                },
                                "$COMMON_NAME": {
                                  "id": "id-a-59",
                                  "type": "common_name",
                                  "value": "named object"
                                }
                              }
                            }
                          }
                        }
                      ],
                      "$NOUN_PHRASE": {
                        "id": "id-a-55",
                        "type": "noun_phrase",
                        "value": {
                          "$KIND_OF_NOUN": {
                            "id": "id-a-63",
                            "type": "kind_of_noun",
                            "value": {
                              "$ARE_AN": {
                                "id": "id-a-61",
                                "type": "are_an",
                                "value": "$ISA"
                              },
                              "$KIND": {
                                "id": "id-a-62",
                                "type": "singular_kind",
                                "value": "thing"
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
              "id": "id-172f216bee4-0",
              "type": "story_statement",
              "value": {
                "id": "id-172f216bee4-4",
                "type": "noun_statement",
                "value": {
                  "$LEDE": {
                    "id": "id-172f216bee4-3",
                    "type": "lede",
                    "value": {
                      "$NOUN": [
                        {
                          "id": "id-172f216bee4-1",
                          "type": "noun",
                          "value": {
                            "$COMMON_NOUN": {
                              "id": "id-172f216bee4-7",
                              "type": "common_noun",
                              "value": {
                                "$DETERMINER": {
                                  "id": "id-172f216bee4-5",
                                  "type": "determiner",
                                  "value": "$THE"
                                },
                                "$COMMON_NAME": {
                                  "id": "id-172f216bee4-6",
                                  "type": "common_name",
                                  "value": "unnamed object"
                                }
                              }
                            }
                          }
                        }
                      ],
                      "$NOUN_PHRASE": {
                        "id": "id-172f216bee4-2",
                        "type": "noun_phrase",
                        "value": {
                          "$KIND_OF_NOUN": {
                            "id": "id-172f216bee4-10",
                            "type": "kind_of_noun",
                            "value": {
                              "$ARE_AN": {
                                "id": "id-172f216bee4-8",
                                "type": "are_an",
                                "value": "$ISA"
                              },
                              "$KIND": {
                                "id": "id-172f216bee4-9",
                                "type": "singular_kind",
                                "value": "thing"
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
              "id": "id-172fbf04bb4-0",
              "type": "story_statement",
              "value": {
                "id": "id-172fbf04bb4-4",
                "type": "noun_statement",
                "value": {
                  "$LEDE": {
                    "id": "id-172fbf04bb4-3",
                    "type": "lede",
                    "value": {
                      "$NOUN": [
                        {
                          "id": "id-172fbf04bb4-1",
                          "type": "noun",
                          "value": {
                            "$COMMON_NOUN": {
                              "id": "id-172fbf04bb4-7",
                              "type": "common_noun",
                              "value": {
                                "$DETERMINER": {
                                  "id": "id-172fbf04bb4-5",
                                  "type": "determiner",
                                  "value": "$THE"
                                },
                                "$COMMON_NAME": {
                                  "id": "id-172fbf04bb4-6",
                                  "type": "common_name",
                                  "value": "unnamed object"
                                }
                              }
                            }
                          }
                        }
                      ],
                      "$NOUN_PHRASE": {
                        "id": "id-172fbf04bb4-2",
                        "type": "noun_phrase",
                        "value": {
                          "$NOUN_TRAITS": {
                            "id": "id-172fbf04bb4-10",
                            "type": "noun_traits",
                            "value": {
                              "$ARE_BEING": {
                                "id": "id-172fbf04bb4-8",
                                "type": "are_being",
                                "value": "$IS"
                              },
                              "$TRAIT": [
                                {
                                  "id": "id-172fbf04bb4-9",
                                  "type": "trait",
                                  "value": "privatelyNamed"
                                }
                              ]
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
      },
      {
        "id": "id-a-27",
        "type": "paragraph",
        "value": {
          "$STORY_STATEMENT": [
            {
              "id": "id-a-26",
              "type": "story_statement",
              "value": {
                "id": "id-a-30",
                "type": "test_statement",
                "value": {
                  "$NAME": {
                    "id": "id-a-28",
                    "type": "text",
                    "value": "testNameOfKind"
                  },
                  "$TEST": {
                    "id": "id-a-29",
                    "type": "testing",
                    "value": {
                      "id": "id-a-33",
                      "type": "test_output",
                      "value": {
                        "$LINES": {
                          "id": "id-a-31",
                          "type": "lines",
                          "value": "named object thing"
                        },
                        "$GO": [
                          {
                            "id": "id-a-64",
                            "type": "execute",
                            "value": {
                              "id": "id-a-66",
                              "type": "say_text",
                              "value": {
                                "$TEXT": {
                                  "id": "id-a-65",
                                  "type": "text_eval",
                                  "value": {
                                    "id": "id-a-68",
                                    "type": "render_template",
                                    "value": {
                                      "$TEMPLATE": {
                                        "id": "id-a-67",
                                        "type": "lines",
                                        "value": "{.named} {.unnamed}"
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
            },
            {
              "id": "id-b-0",
              "type": "story_statement",
              "value": {
                "id": "id-b-3",
                "type": "pattern_handler",
                "value": {
                  "$NAME": {
                    "id": "id-b-1",
                    "type": "pattern_name",
                    "value": "printName"
                  },
                  "$HOOK": {
                    "id": "id-b-2",
                    "type": "pattern_hook",
                    "value": {
                      "$ACTIVITY": {
                        "id": "id-b-5",
                        "type": "pattern_activity",
                        "value": {
                          "$GO": [
                            {
                              "id": "id-172fbf04bb4-11",
                              "type": "execute",
                              "value": {
                                "id": "id-172fbf04bb4-13",
                                "type": "say_text",
                                "value": {
                                  "$TEXT": {
                                    "id": "id-172fbf04bb4-12",
                                    "type": "text_eval",
                                    "value": {
                                      "id": "id-172fbf04bb4-16",
                                      "type": "get_field",
                                      "value": {
                                        "$OBJ": {
                                          "id": "id-172fbf04bb4-14",
                                          "type": "object_name",
                                          "value": {
                                            "$NAME": {
                                              "type": "text_eval",
                                              "value": {
                                                "type": "text_value",
                                                "value": {
                                                  "$TEXT": {
                                                    "type": "text",
                                                    "value": "target"
                                                  }
                                                }
                                              }
                                            }
                                          }
                                        },
                                        "$FIELD": {
                                          "id": "id-172fbf04bb4-15",
                                          "type": "text_eval",
                                          "value": {
                                            "id": "id-172fbf04bb4-20",
                                            "type": "text_value",
                                            "value": {
                                              "$TEXT": {
                                                "id": "id-172fbf04bb4-19",
                                                "type": "text",
                                                "value": "name"
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
                  },
                  "$FILTERS": {
                    "id": "id-172fbf04bb4-22",
                    "type": "pattern_filters",
                    "value": {
                      "$FILTER": [
                        {
                          "id": "id-172fbf04bb4-21",
                          "type": "bool_eval",
                          "value": {
                            "id": "id-172fbf04bb4-27",
                            "type": "is_true",
                            "value": {
                              "$TEST": {
                                "id": "id-172fbf04bb4-26",
                                "type": "bool_eval",
                                "value": {
                                  "id": "id-172fbf04bb4-30",
                                  "type": "get_field",
                                  "value": {
                                    "$OBJ": {
                                      "id": "id-172fbf04bb4-28",
                                      "type": "object_name",
                                      "value": {
                                        "$NAME": {
                                          "type": "text_eval",
                                          "value": {
                                            "type": "text_value",
                                            "value": {
                                              "$TEXT": {
                                                "type": "text",
                                                "value": "target"
                                              }
                                            }
                                          }
                                        }
                                      }
                                    },
                                    "$FIELD": {
                                      "id": "id-172fbf04bb4-29",
                                      "type": "text_eval",
                                      "value": {
                                        "id": "id-172fbf04bb4-34",
                                        "type": "text_value",
                                        "value": {
                                          "$TEXT": {
                                            "id": "id-172fbf04bb4-33",
                                            "type": "text",
                                            "value": "publiclyNamed"
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
            },
            {
              "id": "id-c-12",
              "type": "story_statement",
              "value": {
                "id": "id-c-15",
                "type": "kinds_of_kind",
                "value": {
                  "$PLURAL_KINDS": {
                    "id": "id-c-13",
                    "type": "plural_kinds",
                    "value": " things"
                  },
                  "$SINGULAR_KIND": {
                    "id": "id-c-14",
                    "type": "singular_kind",
                    "value": "kind"
                  }
                }
              }
            }
          ]
        }
      },
      {
        "id": "id-c-17",
        "type": "paragraph",
        "value": {
          "$STORY_STATEMENT": [
            {
              "id": "id-c-16",
              "type": "story_statement",
              "value": {
                "id": "id-c-21",
                "type": "kinds_possess_properties",
                "value": {
                  "$PLURAL_KINDS": {
                    "id": "id-c-18",
                    "type": "plural_kinds",
                    "value": " things"
                  },
                  "$PROPERTY_PHRASE": {
                    "id": "id-c-20",
                    "type": "property_phrase",
                    "value": {
                      "$PRIMITIVE_PHRASE": {
                        "id": "id-c-24",
                        "type": "primitive_phrase",
                        "value": {
                          "$PRIMITIVE_TYPE": {
                            "id": "id-c-22",
                            "type": "primitive_type",
                            "value": "$TEXT"
                          },
                          "$PROPERTY": {
                            "id": "id-c-23",
                            "type": "property",
                            "value": "printedName"
                          }
                        }
                      }
                    }
                  }
                }
              }
            },
            {
              "type": "story_statement",
              "value": {
                "type": "kinds_possess_properties",
                "value": {
                  "$PLURAL_KINDS": {
                    "type": "plural_kinds",
                    "value": " things"
                  },
                  "$PROPERTY_PHRASE": {
                    "type": "property_phrase",
                    "value": {
                      "$PRIMITIVE_PHRASE": {
                        "type": "primitive_phrase",
                        "value": {
                          "$PRIMITIVE_TYPE": {
                            "type": "primitive_type",
                            "value": "$TEXT"
                          },
                          "$PROPERTY": {
                            "type": "property",
                            "value": "pluralName"
                          }
                        }
                      }
                    }
                  }
                }
              }
            },
            {
              "type": "story_statement",
              "value": {
                "type": "kinds_possess_properties",
                "value": {
                  "$PLURAL_KINDS": {
                    "type": "plural_kinds",
                    "value": " things"
                  },
                  "$PROPERTY_PHRASE": {
                    "type": "property_phrase",
                    "value": {
                      "$PRIMITIVE_PHRASE": {
                        "type": "primitive_phrase",
                        "value": {
                          "$PRIMITIVE_TYPE": {
                            "type": "primitive_type",
                            "value": "$TEXT"
                          },
                          "$PROPERTY": {
                            "type": "property",
                            "value": "indefiniteArticle"
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
      },
      {
        "id": "id-172eed384e3-1",
        "type": "paragraph",
        "value": {
          "$STORY_STATEMENT": [
            {
              "id": "id-172fbe17c49-0",
              "type": "story_statement",
              "value": {
                "id": "id-172fbe17c49-3",
                "type": "pattern_handler",
                "value": {
                  "$NAME": {
                    "id": "id-172fbe17c49-1",
                    "type": "pattern_name",
                    "value": "printName"
                  },
                  "$HOOK": {
                    "id": "id-172fbe17c49-2",
                    "type": "pattern_hook",
                    "value": {
                      "$ACTIVITY": {
                        "id": "id-172fbe17c49-18",
                        "type": "pattern_activity",
                        "value": {
                          "$GO": [
                            {
                              "id": "id-172fbe17c49-17",
                              "type": "execute",
                              "value": {
                                "id": "id-172fbe17c49-20",
                                "type": "say_text",
                                "value": {
                                  "$TEXT": {
                                    "id": "id-172fbe17c49-19",
                                    "type": "text_eval",
                                    "value": {
                                      "id": "id-172fbe17c49-23",
                                      "type": "get_field",
                                      "value": {
                                        "$OBJ": {
                                          "id": "id-172fbe17c49-21",
                                          "type": "object_name",
                                          "value": {
                                            "$NAME": {
                                              "type": "text_eval",
                                              "value": {
                                                "type": "text_value",
                                                "value": {
                                                  "$TEXT": {
                                                    "type": "text",
                                                    "value": "target"
                                                  }
                                                }
                                              }
                                            }
                                          }
                                        },
                                        "$FIELD": {
                                          "id": "id-172fbe17c49-22",
                                          "type": "text_eval",
                                          "value": {
                                            "id": "id-172fbe17c49-27",
                                            "type": "text_value",
                                            "value": {
                                              "$TEXT": {
                                                "id": "id-172fbe17c49-26",
                                                "type": "text",
                                                "value": "printedName"
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
                  },
                  "$FILTERS": {
                    "id": "id-172fbe17c49-5",
                    "type": "pattern_filters",
                    "value": {
                      "$FILTER": [
                        {
                          "id": "id-172fbe17c49-4",
                          "type": "bool_eval",
                          "value": {
                            "id": "id-172fbe17c49-7",
                            "type": "is_not",
                            "value": {
                              "$TEST": {
                                "id": "id-172fbe17c49-6",
                                "type": "bool_eval",
                                "value": {
                                  "id": "id-172fbe17c49-9",
                                  "type": "is_empty",
                                  "value": {
                                    "$TEXT": {
                                      "id": "id-172fbe17c49-8",
                                      "type": "text_eval",
                                      "value": {
                                        "id": "id-172fbe17c49-12",
                                        "type": "get_field",
                                        "value": {
                                          "$OBJ": {
                                            "id": "id-172fbe17c49-10",
                                            "type": "object_name",
                                            "value": {
                                              "$NAME": {
                                                "type": "text_eval",
                                                "value": {
                                                  "type": "text_value",
                                                  "value": {
                                                    "$TEXT": {
                                                      "type": "text",
                                                      "value": "target"
                                                    }
                                                  }
                                                }
                                              }
                                            }
                                          },
                                          "$FIELD": {
                                            "id": "id-172fbe17c49-11",
                                            "type": "text_eval",
                                            "value": {
                                              "id": "id-172fbe17c49-16",
                                              "type": "text_value",
                                              "value": {
                                                "$TEXT": {
                                                  "id": "id-172fbe17c49-15",
                                                  "type": "text",
                                                  "value": "printedName"
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
      }
    ]
  }
}
}
