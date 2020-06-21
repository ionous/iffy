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
                      "type": "test_out",
                      "value": {
                        "$LINES": {
                          "id": "id-171c8ccc9f6-6",
                          "type": "lines",
                          "value": "6"
                        },
                        "$GO": [
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
                                            "$PARAMETERS": {
                                              "id": "id-171c8ccc9f6-18",
                                              "type": "parameters",
                                              "value": {
                                                "$PARAMS": [
                                                  {
                                                    "id": "id-171c8ccc9f6-17",
                                                    "type": "parameter",
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
                "id": "id-171ccf3429f-16",
                "type": "pattern_handler",
                "value": {
                  "$NAME": {
                    "id": "id-171ccf3429f-13",
                    "type": "pattern_name",
                    "value": "factorial"
                  },
                  "$HOOK": {
                    "id": "id-171ccf3429f-15",
                    "type": "pattern_hook",
                    "value": {
                      "$RESULT": {
                        "id": "id-171ccf3429f-20",
                        "type": "pattern_return",
                        "value": {
                          "$RESULT": {
                            "id": "id-171ccf3429f-19",
                            "type": "pattern_result",
                            "value": {
                              "$PRIMITIVE": {
                                "id": "id-171ccf3429f-21",
                                "type": "primitive_func",
                                "value": {
                                  "$NUMBER_EVAL": {
                                    "id": "id-171ccf3429f-22",
                                    "type": "number_eval",
                                    "value": {
                                      "id": "id-171ccf3429f-24",
                                      "type": "num_value",
                                      "value": {
                                        "$NUM": {
                                          "id": "id-171ccf3429f-23",
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
                  },
                  "$FILTERS": {
                    "id": "id-171cd1c092f-1",
                    "type": "pattern_filters",
                    "value": {
                      "$FILTER": [
                        {
                          "id": "id-171cd1c092f-0",
                          "type": "bool_eval",
                          "value": {
                            "id": "id-171cd1c092f-5",
                            "type": "compare_num",
                            "value": {
                              "$A": {
                                "id": "id-171cd1c092f-2",
                                "type": "number_eval",
                                "value": {
                                  "id": "id-171cd1c092f-7",
                                  "type": "get_var",
                                  "value": {
                                    "$NAME": {
                                      "id": "id-171cd1c092f-6",
                                      "type": "text",
                                      "value": "num"
                                    }
                                  }
                                }
                              },
                              "$IS": {
                                "id": "id-171cd1c092f-3",
                                "type": "compare_to",
                                "value": {
                                  "id": "id-171cd1c092f-8",
                                  "type": "equal",
                                  "value": {}
                                }
                              },
                              "$B": {
                                "id": "id-171cd1c092f-4",
                                "type": "number_eval",
                                "value": {
                                  "id": "id-171cd1c092f-10",
                                  "type": "num_value",
                                  "value": {
                                    "$NUM": {
                                      "id": "id-171cd1c092f-9",
                                      "type": "number",
                                      "value": 0
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
                "id": "id-171ccf3429f-38",
                "type": "pattern_handler",
                "value": {
                  "$NAME": {
                    "id": "id-171ccf3429f-35",
                    "type": "pattern_name",
                    "value": "factorial"
                  },
                  "$HOOK": {
                    "id": "id-171ccf3429f-37",
                    "type": "pattern_hook",
                    "value": {
                      "$RESULT": {
                        "id": "id-171ccf3429f-42",
                        "type": "pattern_return",
                        "value": {
                          "$RESULT": {
                            "id": "id-171ccf3429f-41",
                            "type": "pattern_result",
                            "value": {
                              "$PRIMITIVE": {
                                "id": "id-171ccf3429f-43",
                                "type": "primitive_func",
                                "value": {
                                  "$NUMBER_EVAL": {
                                    "id": "id-171ccf3429f-44",
                                    "type": "number_eval",
                                    "value": {
                                      "id": "id-171ccf3429f-47",
                                      "type": "product_of",
                                      "value": {
                                        "$A": {
                                          "id": "id-171ccf3429f-45",
                                          "type": "number_eval",
                                          "value": {
                                            "id": "id-171ccf3429f-49",
                                            "type": "get_var",
                                            "value": {
                                              "$NAME": {
                                                "id": "id-171ccf3429f-48",
                                                "type": "text",
                                                "value": "num"
                                              }
                                            }
                                          }
                                        },
                                        "$B": {
                                          "id": "id-171ccf3429f-46",
                                          "type": "number_eval",
                                          "value": {
                                            "id": "id-171ccf3429f-52",
                                            "type": "diff_of",
                                            "value": {
                                              "$A": {
                                                "id": "id-171ccf3429f-50",
                                                "type": "number_eval",
                                                "value": {
                                                  "id": "id-171ccf3429f-54",
                                                  "type": "get_var",
                                                  "value": {
                                                    "$NAME": {
                                                      "id": "id-171ccf3429f-53",
                                                      "type": "text",
                                                      "value": "num"
                                                    }
                                                  }
                                                }
                                              },
                                              "$B": {
                                                "id": "id-171ccf3429f-51",
                                                "type": "number_eval",
                                                "value": {
                                                  "id": "id-171ccf3429f-56",
                                                  "type": "num_value",
                                                  "value": {
                                                    "$NUM": {
                                                      "id": "id-171ccf3429f-55",
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
            }
          ]
        }
      }
    ]
  }
}
}
