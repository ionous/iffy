function getStory() {
  return {
  "id": "id-1724419eaa5-2",
  "type": "story",
  "value": {
    "$PARAGRAPH": [
      {
        "id": "id-1724419eaa5-1",
        "type": "paragraph",
        "value": {
          "$STORY_STATEMENT": [
            {
              "id": "id-1724419eaa5-0",
              "type": "story_statement",
              "value": {
                "id": "id-1724419eaa5-5",
                "type": "pattern_decl",
                "value": {
                  "$NAME": {
                    "id": "id-1724419eaa5-3",
                    "type": "pattern_name",
                    "value": "print name"
                  },
                  "$TYPE": {
                    "id": "id-1724419eaa5-4",
                    "type": "pattern_type",
                    "value": {
                      "$ACTIVITY": {
                        "id": "id-1724419eaa5-6",
                        "type": "patterned_activity",
                        "value": "$ACTIVITY"
                      }
                    }
                  },
                  "$OPTVARS": {
                    "id": "id-1724419eaa5-10",
                    "type": "pattern_variables_tail",
                    "value": {
                      "$VARIABLE_DECL": [
                        {
                          "id": "id-1724419eaa5-9",
                          "type": "variable_decl",
                          "value": {
                            "$TYPE": {
                              "id": "id-1724419eaa5-7",
                              "type": "variable_type",
                              "value": {
                                "$OBJECT": {
                                  "id": "id-1724419eaa5-13",
                                  "type": "object_type",
                                  "value": {
                                    "$AN": {
                                      "id": "id-1724419eaa5-11",
                                      "type": "an",
                                      "value": "$A"
                                    },
                                    "$KINDS": {
                                      "id": "id-1724419eaa5-12",
                                      "type": "plural_kinds",
                                      "value": "things"
                                    }
                                  }
                                }
                              }
                            },
                            "$NAME": {
                              "id": "id-1724419eaa5-8",
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
        "id": "id-1724419eaa5-15",
        "type": "paragraph",
        "value": {
          "$STORY_STATEMENT": [
            {
              "id": "id-1724419eaa5-14",
              "type": "story_statement",
              "value": {
                "id": "id-1724419eaa5-18",
                "type": "pattern_handler",
                "value": {
                  "$NAME": {
                    "id": "id-1724419eaa5-16",
                    "type": "pattern_name",
                    "value": "print name"
                  },
                  "$HOOK": {
                    "id": "id-1724419eaa5-17",
                    "type": "pattern_hook",
                    "value": {
                      "$ACTIVITY": {
                        "id": "id-1724419eaa5-20",
                        "type": "pattern_activity",
                        "value": {
                          "$GO": [
                            {
                              "id": "id-1724419eaa5-19",
                              "type": "execute",
                              "value": {
                                "id": "id-1724419eaa5-23",
                                "type": "say_text",
                                "value": {
                                  "$TEXT": {
                                    "id": "id-1724419eaa5-22",
                                    "type": "text_eval",
                                    "value": {
                                      "id": "id-1724419eaa5-25",
                                      "type": "render_template",
                                      "value": {
                                        "$TEMPLATE": {
                                          "id": "id-1724419eaa5-24",
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
              "id": "id-1724419eaa5-21",
              "type": "story_statement",
              "value": {
                "id": "id-1724419eaa5-57",
                "type": "noun_statement",
                "value": {
                  "$LEDE": {
                    "id": "id-1724419eaa5-56",
                    "type": "lede",
                    "value": {
                      "$NOUN": [
                        {
                          "id": "id-1724419eaa5-54",
                          "type": "noun",
                          "value": {
                            "$COMMON_NOUN": {
                              "id": "id-1724419eaa5-60",
                              "type": "common_noun",
                              "value": {
                                "$DETERMINER": {
                                  "id": "id-1724419eaa5-58",
                                  "type": "determiner",
                                  "value": "$THE"
                                },
                                "$COMMON_NAME": {
                                  "id": "id-1724419eaa5-59",
                                  "type": "common_name",
                                  "value": "example"
                                }
                              }
                            }
                          }
                        }
                      ],
                      "$NOUN_PHRASE": {
                        "id": "id-1724419eaa5-55",
                        "type": "noun_phrase",
                        "value": {
                          "$KIND_OF_NOUN": {
                            "id": "id-1724419eaa5-63",
                            "type": "kind_of_noun",
                            "value": {
                              "$ARE_AN": {
                                "id": "id-1724419eaa5-61",
                                "type": "are_an",
                                "value": "$ISA"
                              },
                              "$KIND": {
                                "id": "id-1724419eaa5-62",
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
        "id": "id-1724419eaa5-27",
        "type": "paragraph",
        "value": {
          "$STORY_STATEMENT": [
            {
              "id": "id-1724419eaa5-26",
              "type": "story_statement",
              "value": {
                "id": "id-1724419eaa5-30",
                "type": "test_statement",
                "value": {
                  "$NAME": {
                    "id": "id-1724419eaa5-28",
                    "type": "text",
                    "value": "test name of kind"
                  },
                  "$TEST": {
                    "id": "id-1724419eaa5-29",
                    "type": "testing",
                    "value": {
                      "id": "id-1724419eaa5-33",
                      "type": "test_output",
                      "value": {
                        "$LINES": {
                          "id": "id-1724419eaa5-31",
                          "type": "lines",
                          "value": "thing"
                        },
                        "$GO": [
                          {
                            "id": "id-1724419eaa5-64",
                            "type": "execute",
                            "value": {
                              "id": "id-1724419eaa5-66",
                              "type": "say_text",
                              "value": {
                                "$TEXT": {
                                  "id": "id-1724419eaa5-65",
                                  "type": "text_eval",
                                  "value": {
                                    "id": "id-1724419eaa5-68",
                                    "type": "render_template",
                                    "value": {
                                      "$TEMPLATE": {
                                        "id": "id-1724419eaa5-67",
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
              "id": "id-172a53c6656-0",
              "type": "story_statement",
              "value": {
                "id": "id-172a53c6656-3",
                "type": "pattern_handler",
                "value": {
                  "$NAME": {
                    "id": "id-172a53c6656-1",
                    "type": "pattern_name",
                    "value": "print name"
                  },
                  "$HOOK": {
                    "id": "id-172a53c6656-2",
                    "type": "pattern_hook",
                    "value": {
                      "$ACTIVITY": {
                        "id": "id-172a53c6656-5",
                        "type": "pattern_activity",
                        "value": {
                          "$GO": [
                            {
                              "id": "id-172a53c6656-4",
                              "type": "execute",
                              "value": {
                                "id": "id-172a53c6656-7",
                                "type": "choose",
                                "value": {
                                  "$TRUE": [
                                    {
                                      "id": "id-172a6515f9d-0",
                                      "type": "execute",
                                      "value": {
                                        "id": "id-172a6515f9d-8",
                                        "type": "say_text",
                                        "value": {
                                          "$TEXT": {
                                            "id": "id-172a6515f9d-7",
                                            "type": "text_eval",
                                            "value": {
                                              "id": "id-172a6515f9d-10",
                                              "type": "render_template",
                                              "value": {
                                                "$TEMPLATE": {
                                                  "id": "id-172a6515f9d-9",
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
                                    "id": "id-172a53c6656-6",
                                    "type": "bool_eval",
                                    "value": {
                                      "id": "id-172a53c6656-11",
                                      "type": "is_true",
                                      "value": {
                                        "$TEST": {
                                          "id": "id-172a53c6656-10",
                                          "type": "bool_eval",
                                          "value": {
                                            "id": "id-172a53c6656-14",
                                            "type": "get_field",
                                            "value": {
                                              "$OBJ": {
                                                "id": "id-172a53c6656-12",
                                                "type": "text_eval",
                                                "value": {
                                                  "id": "id-172a53c6656-16",
                                                  "type": "text_value",
                                                  "value": {
                                                    "$TEXT": {
                                                      "id": "id-172a53c6656-15",
                                                      "type": "lines",
                                                      "value": "target"
                                                    }
                                                  }
                                                }
                                              },
                                              "$FIELD": {
                                                "id": "id-172a53c6656-13",
                                                "type": "text_eval",
                                                "value": {
                                                  "id": "id-172a53c6656-18",
                                                  "type": "text_value",
                                                  "value": {
                                                    "$TEXT": {
                                                      "id": "id-172a53c6656-17",
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
              "id": "id-172a6515f9d-12",
              "type": "story_statement",
              "value": {
                "id": "id-172a6515f9d-15",
                "type": "kinds_of_kind",
                "value": {
                  "$PLURAL_KINDS": {
                    "id": "id-172a6515f9d-13",
                    "type": "plural_kinds",
                    "value": "things"
                  },
                  "$SINGULAR_KIND": {
                    "id": "id-172a6515f9d-14",
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
        "id": "id-172a6515f9d-17",
        "type": "paragraph",
        "value": {
          "$STORY_STATEMENT": [
            {
              "id": "id-172a6515f9d-16",
              "type": "story_statement",
              "value": {
                "id": "id-172a6515f9d-21",
                "type": "kinds_possess_properties",
                "value": {
                  "$PLURAL_KINDS": {
                    "id": "id-172a6515f9d-18",
                    "type": "plural_kinds",
                    "value": "kinds"
                  },
                  "$PROPERTY_PHRASE": {
                    "id": "id-172a6515f9d-20",
                    "type": "property_phrase",
                    "value": {
                      "$PRIMITIVE_PHRASE": {
                        "id": "id-172a6515f9d-24",
                        "type": "primitive_phrase",
                        "value": {
                          "$PRIMITIVE_TYPE": {
                            "id": "id-172a6515f9d-22",
                            "type": "primitive_type",
                            "value": "$TEXT"
                          },
                          "$PROPERTY_NAME": {
                            "id": "id-172a6515f9d-23",
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
              "id": "id-172a6ab9b69-0",
              "type": "story_statement",
              "value": {
                "id": "id-172a6ab9b69-7",
                "type": "aspect_traits",
                "value": {
                  "$PLURAL_KINDS": {
                    "id": "id-172a6ab9b69-3",
                    "type": "plural_kinds",
                    "value": "Kinds"
                  },
                  "$TRAIT_PHRASE": {
                    "id": "id-172a6ab9b69-6",
                    "type": "attribute_phrase",
                    "value": {
                      "$ARE_EITHER": {
                        "id": "id-172a6ab9b69-4",
                        "type": "are_either",
                        "value": "$EITHER"
                      },
                      "$TRAIT": [
                        {
                          "id": "id-172a6ab9b69-5",
                          "type": "trait",
                          "value": "common named"
                        },
                        {
                          "id": "id-172a6ab9b69-8",
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
              "id": "id-172a6ab9b69-9",
              "type": "story_statement",
              "value": {
                "id": "id-172a6ab9b69-17",
                "type": "aspect_traits",
                "value": {
                  "$PLURAL_KINDS": {
                    "id": "id-172a6ab9b69-13",
                    "type": "plural_kinds",
                    "value": "Kinds"
                  },
                  "$TRAIT_PHRASE": {
                    "id": "id-172a6ab9b69-16",
                    "type": "attribute_phrase",
                    "value": {
                      "$ARE_EITHER": {
                        "id": "id-172a6ab9b69-14",
                        "type": "are_either",
                        "value": "$EITHER"
                      },
                      "$TRAIT": [
                        {
                          "id": "id-172a6ab9b69-15",
                          "type": "trait",
                          "value": "singular named"
                        },
                        {
                          "id": "id-172a6ab9b69-18",
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
