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
                    "value": "print name"
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
                                    "$KINDS": {
                                      "id": "id-a-12",
                                      "type": "plural_kinds",
                                      "value": "things"
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
                    "value": "print name"
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
                              "id": "id-a-19",
                              "type": "execute",
                              "value": {
                                "id": "id-a-23",
                                "type": "say_text",
                                "value": {
                                  "$TEXT": {
                                    "id": "id-a-22",
                                    "type": "text_eval",
                                    "value": {
                                      "id": "id-a-25",
                                      "type": "render_template",
                                      "value": {
                                        "$TEMPLATE": {
                                          "id": "id-a-24",
                                          "type": "lines",
                                          "value": "{nameOfKind: .target}"
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
                                  "value": "example"
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
                    "value": "test name of kind"
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
                          "value": "thing"
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
                                        "value": "{.example}"
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
                    "value": "print name"
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
                              "id": "id-b-4",
                              "type": "execute",
                              "value": {
                                "id": "id-b-7",
                                "type": "choose",
                                "value": {
                                  "$TRUE": [
                                    {
                                      "id": "id-c-0",
                                      "type": "execute",
                                      "value": {
                                        "id": "id-c-8",
                                        "type": "say_text",
                                        "value": {
                                          "$TEXT": {
                                            "id": "id-c-7",
                                            "type": "text_eval",
                                            "value": {
                                              "id": "id-c-10",
                                              "type": "render_template",
                                              "value": {
                                                "$TEMPLATE": {
                                                  "id": "id-c-9",
                                                  "type": "lines",
                                                  "value": "{.target.printedName}"
                                                }
                                              }
                                            }
                                          }
                                        }
                                      }
                                    }
                                  ],
                                  "$FALSE": [],
                                  "$IF": {
                                    "id": "id-b-6",
                                    "type": "bool_eval",
                                    "value": {
                                      "id": "id-b-11",
                                      "type": "is_true",
                                      "value": {
                                        "$TEST": {
                                          "id": "id-b-10",
                                          "type": "bool_eval",
                                          "value": {
                                            "id": "id-b-14",
                                            "type": "get_field",
                                            "value": {
                                              "$OBJ": {
                                                "id": "id-b-12",
                                                "type": "text_eval",
                                                "value": {
                                                  "id": "id-b-16",
                                                  "type": "text_value",
                                                  "value": {
                                                    "$TEXT": {
                                                      "id": "id-b-15",
                                                      "type": "lines",
                                                      "value": "target"
                                                    }
                                                  }
                                                }
                                              },
                                              "$FIELD": {
                                                "id": "id-b-13",
                                                "type": "text_eval",
                                                "value": {
                                                  "id": "id-b-18",
                                                  "type": "text_value",
                                                  "value": {
                                                    "$TEXT": {
                                                      "id": "id-b-17",
                                                      "type": "lines",
                                                      "value": "publicly-named"
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
              "id": "id-c-12",
              "type": "story_statement",
              "value": {
                "id": "id-c-15",
                "type": "kinds_of_kind",
                "value": {
                  "$PLURAL_KINDS": {
                    "id": "id-c-13",
                    "type": "plural_kinds",
                    "value": "things"
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
                    "value": "kinds"
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
                          "$PROPERTY_NAME": {
                            "id": "id-c-23",
                            "type": "property_name",
                            "value": "printed name"
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
                    "value": "kinds"
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
                          "$PROPERTY_NAME": {
                            "type": "property_name",
                            "value": "printed plural name"
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
                    "value": "kinds"
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
                          "$PROPERTY_NAME": {
                            "type": "property_name",
                            "value": "indefinite article"
                          }
                        }
                      }
                    }
                  }
                }
              }
            },
            {
              "id": "id-d-0",
              "type": "story_statement",
              "value": {
                "id": "id-d-7",
                "type": "aspect_traits",
                "value": {
                  "$PLURAL_KINDS": {
                    "id": "id-d-3",
                    "type": "plural_kinds",
                    "value": "Kinds"
                  },
                  "$TRAIT_PHRASE": {
                    "id": "id-d-6",
                    "type": "attribute_phrase",
                    "value": {
                      "$ARE_EITHER": {
                        "id": "id-d-4",
                        "type": "are_either",
                        "value": "$EITHER"
                      },
                      "$TRAIT": [
                        {
                          "id": "id-d-5",
                          "type": "trait",
                          "value": "common named"
                        },
                        {
                          "id": "id-d-8",
                          "type": "trait",
                          "value": "proper named"
                        }
                      ]
                    }
                  }
                }
              }
            },
            {
              "id": "id-d-9",
              "type": "story_statement",
              "value": {
                "id": "id-d-17",
                "type": "aspect_traits",
                "value": {
                  "$PLURAL_KINDS": {
                    "id": "id-d-13",
                    "type": "plural_kinds",
                    "value": "Kinds"
                  },
                  "$TRAIT_PHRASE": {
                    "id": "id-d-16",
                    "type": "attribute_phrase",
                    "value": {
                      "$ARE_EITHER": {
                        "id": "id-d-14",
                        "type": "are_either",
                        "value": "$EITHER"
                      },
                      "$TRAIT": [
                        {
                          "id": "id-d-15",
                          "type": "trait",
                          "value": "singular named"
                        },
                        {
                          "id": "id-d-18",
                          "type": "trait",
                          "value": "plural named"
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
