function getStory() {
  return {
    "id": "id-1722f98760b-2",
    "type": "story",
    "value": {
      "$PARAGRAPH": [
        {
          "id": "id-1722f98760b-1",
          "type": "paragraph",
          "value": {
            "$STORY_STATEMENT": [
              {
                "id": "id-1722f98760b-0",
                "type": "story_statement",
                "value": {
                  "id": "id-1722f98760b-5",
                  "type": "test_statement",
                  "value": {
                    "$NAME": {
                      "id": "id-1722f98760b-3",
                      "type": "text",
                      "value": "cycle"
                    },
                    "$TEST": {
                      "id": "id-1722f98760b-4",
                      "type": "testing",
                      "value": {
                        "id": "id-1722f98760b-8",
                        "type": "test_output",
                        "value": {
                          "$LINES": {
                            "id": "id-1722f98760b-6",
                            "type": "lines",
                            "value": "a\nb\nc"
                          },
                          "$GO": [
                            {
                              "id": "id-1722fcbabd5-0",
                              "type": "execute",
                              "value": {
                                "id": "id-1722fcbabd5-4",
                                "type": "for_each_num",
                                "value": {
                                  "$ELSE": [
                                    {
                                      "id": "id-1722fcbabd5-1",
                                      "type": "execute",
                                      "value": null
                                    }
                                  ],
                                  "$GO": [
                                    {
                                      "id": "id-1722fcbabd5-2",
                                      "type": "execute",
                                      "value": {
                                        "id": "id-1722fcbabd5-16",
                                        "type": "say_text",
                                        "value": {
                                          "$TEXT": {
                                            "id": "id-1722fcbabd5-15",
                                            "type": "text_eval",
                                            "value": {
                                              "id": "id-1722fcbabd5-18",
                                              "type": "cycle_text",
                                              "value": {
                                                "$PARTS": [
                                                  {
                                                    "id": "id-1722fcbabd5-17",
                                                    "type": "text_eval",
                                                    "value": {
                                                      "id": "id-1722fcbabd5-22",
                                                      "type": "text_value",
                                                      "value": {
                                                        "$TEXT": {
                                                          "id": "id-1722fcbabd5-21",
                                                          "type": "lines",
                                                          "value": "a"
                                                        }
                                                      }
                                                    }
                                                  },
                                                  {
                                                    "id": "id-1722fcbabd5-19",
                                                    "type": "text_eval",
                                                    "value": {
                                                      "id": "id-1722fcbabd5-24",
                                                      "type": "text_value",
                                                      "value": {
                                                        "$TEXT": {
                                                          "id": "id-1722fcbabd5-23",
                                                          "type": "lines",
                                                          "value": "b"
                                                        }
                                                      }
                                                    }
                                                  },
                                                  {
                                                    "id": "id-1722fcbabd5-20",
                                                    "type": "text_eval",
                                                    "value": {
                                                      "id": "id-1722fcbabd5-26",
                                                      "type": "text_value",
                                                      "value": {
                                                        "$TEXT": {
                                                          "id": "id-1722fcbabd5-25",
                                                          "type": "lines",
                                                          "value": "c"
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
                                  ],
                                  "$IN": {
                                    "id": "id-1722fcbabd5-3",
                                    "type": "num_list_eval",
                                    "value": {
                                      "id": "id-1722fcbabd5-8",
                                      "type": "range_over",
                                      "value": {
                                        "$START": {
                                          "id": "id-1722fcbabd5-5",
                                          "type": "number_eval",
                                          "value": {
                                            "id": "id-1722fcbabd5-10",
                                            "type": "num_value",
                                            "value": {
                                              "$NUM": {
                                                "id": "id-1722fcbabd5-9",
                                                "type": "number",
                                                "value": 0
                                              }
                                            }
                                          }
                                        },
                                        "$STEP": {
                                          "id": "id-1722fcbabd5-6",
                                          "type": "number_eval",
                                          "value": {
                                            "id": "id-1722fcbabd5-14",
                                            "type": "num_value",
                                            "value": {
                                              "$NUM": {
                                                "id": "id-1722fcbabd5-13",
                                                "type": "number",
                                                "value": 1
                                              }
                                            }
                                          }
                                        },
                                        "$STOP": {
                                          "id": "id-1722fcbabd5-7",
                                          "type": "number_eval",
                                          "value": {
                                            "id": "id-1722fcbabd5-12",
                                            "type": "num_value",
                                            "value": {
                                              "$NUM": {
                                                "id": "id-1722fcbabd5-11",
                                                "type": "number",
                                                "value": 3
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
            ]
          }
        }
      ]
    }
  }
}
