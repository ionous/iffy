function getStory() {
  return {
    "id": "id-171c8ccc9f6-2",
    "type": "story",
    "value": {
      "$PARAGRAPH": [
        {
          "id": "id-171c8ccc9f6-1",
          "type": "paragraph",
          "value": {
            "$STORY_STATEMENT": [
              {
                "id": "id-171c8ccc9f6-0",
                "type": "story_statement",
                "value": {
                  "id": "id-171c8ccc9f6-4",
                  "type": "test_statement",
                  "value": {
                    "$NAME": {
                      "id": "id-171c8ccc9f6-5",
                      "type": "text",
                      "value": "factorial"
                    },
                    "$TEST": {
                      "id": "id-171c8ccc9f6-3",
                      "type": "testing",
                      "value": {
                        "id": "id-171c8ccc9f6-8",
                        "type": "test_output",
                        "value": {
                          "$LINES": {
                            "id": "id-171c8ccc9f6-6",
                            "type": "lines",
                            "value": "6"
                          },
                          "$GO": {
                            "type": "activity",
                            "value": {
                              "$EXE": [
                                {
                                  "id": "id-171c8ccc9f6-7",
                                  "type": "execute",
                                  "value": {
                                    "id": "id-171c8ccc9f6-10",
                                    "type": "say_text",
                                    "value": {
                                      "$TEXT": {
                                        "id": "id-171c8ccc9f6-9",
                                        "type": "text_eval",
                                        "value": {
                                          "id": "id-171c8ccc9f6-12",
                                          "type": "print_num",
                                          "value": {
                                            "$NUM": {
                                              "id": "id-171c8ccc9f6-11",
                                              "type": "number_eval",
                                              "value": {
                                                "id": "id-171c8ccc9f6-14",
                                                "type": "determine_num",
                                                "value": {
                                                  "$PATTERN": {
                                                    "id": "id-171c8ccc9f6-13",
                                                    "type": "pattern_name",
                                                    "value": "factorial"
                                                  },
                                                  "$ARGUMENTS": {
                                                    "id": "id-171c8ccc9f6-18",
                                                    "type": "arguments",
                                                    "value": {
                                                      "$ARGS": [
                                                        {
                                                          "id": "id-171c8ccc9f6-17",
                                                          "type": "argument",
                                                          "value": {
                                                            "$FROM": {
                                                              "id": "id-171c8ccc9f6-15",
                                                              "type": "assignment",
                                                              "value": {
                                                                "id": "id-171c8ccc9f6-20",
                                                                "type": "assign_num",
                                                                "value": {
                                                                  "$VAL": {
                                                                    "id": "id-171c8ccc9f6-19",
                                                                    "type": "number_eval",
                                                                    "value": {
                                                                      "id": "id-171c8ccc9f6-22",
                                                                      "type": "num_value",
                                                                      "value": {
                                                                        "$NUM": {
                                                                          "id": "id-171c8ccc9f6-21",
                                                                          "type": "number",
                                                                          "value": 3
                                                                        }
                                                                      }
                                                                    }
                                                                  }
                                                                }
                                                              }
                                                            },
                                                            "$NAME": {
                                                              "id": "id-171c8ccc9f6-16",
                                                              "type": "variable_name",
                                                              "value": "num"
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
              },
              {
                "id": "id-171ccf3429f-0",
                "type": "story_statement",
                "value": {
                  "id": "id-171ccf3429f-3",
                  "type": "pattern_decl",
                  "value": {
                    "$NAME": {
                      "id": "id-171ccf3429f-1",
                      "type": "pattern_name",
                      "value": "factorial"
                    },
                    "$TYPE": {
                      "id": "id-171ccf3429f-2",
                      "type": "pattern_type",
                      "value": {
                        "$VALUE": {
                          "id": "id-171ccf3429f-4",
                          "type": "variable_type",
                          "value": {
                            "$PRIMITIVE": {
                              "id": "id-171ccf3429f-5",
                              "type": "primitive_type",
                              "value": "$NUMBER"
                            }
                          }
                        }
                      }
                    },
                    "$OPTVARS": {
                      "id": "id-171ccf3429f-9",
                      "type": "pattern_variables_tail",
                      "value": {
                        "$VARIABLE_DECL": [
                          {
                            "id": "id-171ccf3429f-8",
                            "type": "variable_decl",
                            "value": {
                              "$TYPE": {
                                "id": "id-171ccf3429f-6",
                                "type": "variable_type",
                                "value": {
                                  "$PRIMITIVE": {
                                    "id": "id-171ccf3429f-10",
                                    "type": "primitive_type",
                                    "value": "$NUMBER"
                                  }
                                }
                              },
                              "$NAME": {
                                "id": "id-171ccf3429f-7",
                                "type": "variable_name",
                                "value": "num"
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
          "id": "id-171ccf3429f-12",
          "type": "paragraph",
          "value": {
            "$STORY_STATEMENT": [
              {
                "id": "id-171ccf3429f-11",
                "type": "story_statement",
                "value": {
                  "type": "pattern_actions",
                  "value": {
                    "$NAME": {
                      "type": "pattern_name",
                      "value": "factorial"
                    },
                    "$PATTERN_RULES": {
                      "type": "pattern_rules",
                      "value": {
                        "$PATTERN_RULE": [
                          {
                            "type": "pattern_rule",
                            "value": {
                              "$GUARD": {
                                "type": "bool_eval",
                                "value": {
                                  "type": "compare_num",
                                  "value": {
                                    "$A": {
                                      "type": "number_eval",
                                      "value": {
                                        "type": "get_var",
                                        "value": {
                                          "$NAME": {
                                            "type": "text",
                                            "value": "num"
                                          }
                                        }
                                      }
                                    },
                                    "$B": {
                                      "type": "number_eval",
                                      "value": {
                                        "type": "num_value",
                                        "value": {
                                          "$NUM": {
                                            "type": "number",
                                            "value": 0
                                          }
                                        }
                                      }
                                    },
                                    "$IS": {
                                      "type": "comparator",
                                      "value": {
                                        "type": "equal",
                                        "value": {}
                                      }
                                    }
                                  }
                                }
                              },
                              "$HOOK": {
                                "type": "program_hook",
                                "value": {
                                  "$RESULT": {
                                    "type": "program_return",
                                    "value": {
                                      "$RESULT": {
                                        "type": "program_result",
                                        "value": {
                                          "$PRIMITIVE": {
                                            "type": "primitive_func",
                                            "value": {
                                              "$NUMBER_EVAL": {
                                                "type": "number_eval",
                                                "value": {
                                                  "type": "num_value",
                                                  "value": {
                                                    "$NUM": {
                                                      "type": "number",
                                                      "value": 1
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
              },
              {
                "id": "id-171ccf3429f-34",
                "type": "story_statement",


                "value": {
                  "type": "pattern_actions",
                  "value": {
                    "$NAME": {
                      "type": "pattern_name",
                      "value": "factorial"
                    },
                    "$PATTERN_RULES": {
                      "type": "pattern_rules",
                      "value": {
                        "$PATTERN_RULE": [
                          {
                            "type": "pattern_rule",
                            "value": {
                              "$GUARD": {
                                "type": "bool_eval",
                                "value": {
                                  "type": "always",
                                  "value": {}
                                }
                              },
                              "$HOOK": {
                                "type": "program_hook",
                                "value": {
                                  "$RESULT": {
                                    "type": "program_return",
                                    "value": {
                                      "$RESULT": {
                                        "type": "program_result",
                                        "value": {
                                          "$PRIMITIVE": {
                                            "type": "primitive_func",
                                            "value": {
                                              "$NUMBER_EVAL": {
                                                "type": "number_eval",
                                                "value": {
                                                  "type": "product_of",
                                                  "value": {
                                                    "$A": {
                                                      "type": "number_eval",
                                                      "value": {
                                                        "type": "get_var",
                                                        "value": {
                                                          "$NAME": {
                                                            "type": "text",
                                                            "value": "num"
                                                          }
                                                        }
                                                      }
                                                    },
                                                    "$B": {
                                                      "type": "number_eval",
                                                      "value": {
                                                        "type": "diff_of",
                                                        "value": {
                                                          "$A": {
                                                            "type": "number_eval",
                                                            "value": {
                                                              "type": "get_var",
                                                              "value": {
                                                                "$NAME": {
                                                                  "type": "text",
                                                                  "value": "num"
                                                                }
                                                              }
                                                            }
                                                          },
                                                          "$B": {
                                                            "type": "number_eval",
                                                            "value": {
                                                              "type": "num_value",
                                                              "value": {
                                                                "$NUM": {
                                                                  "type": "number",
                                                                  "value": 1
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
