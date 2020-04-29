function getStory() {
return {
  "id": "id-1709ef632af-3",
  "type": "test",
  "value": {
    "$TEST_NAME": {
      "id": "id-1709ef632af-0",
      "type": "text",
      "value": "factorials"
    },
    "$GO": [
      {
        "id": "id-171c4a25ccd-0",
        "type": "execute",
        "value": {
          "id": "id-171c4a25ccd-2",
          "type": "say_text",
          "value": {
            "$TEXT": {
              "id": "id-171c4a25ccd-1",
              "type": "text_eval",
              "value": {
                "id": "id-171c4a25ccd-4",
                "type": "print_num",
                "value": {
                  "$NUM": {
                    "id": "id-171c4a25ccd-3",
                    "type": "number_eval",
                    "value": {
                      "id": "id-171c4a25ccd-9",
                      "type": "determine_num",
                      "value": {
                        "$NAME": {
                          "id": "id-171c4a25ccd-5",
                          "type": "text",
                          "value": "factorial"
                        },
                        "$PARAMETERS": [
                          {
                            "id": "id-171c4a25ccd-8",
                            "type": "parameter",
                            "value": {
                              "$FROM": {
                                "id": "id-171c4a25ccd-6",
                                "type": "assignment",
                                "value": {
                                  "id": "id-171c4a25ccd-11",
                                  "type": "assign_num",
                                  "value": {
                                    "$VAL": {
                                      "id": "id-171c4a25ccd-10",
                                      "type": "number_eval",
                                      "value": {
                                        "id": "id-171c4a25ccd-13",
                                        "type": "num_value",
                                        "value": {
                                          "$NUM": {
                                            "id": "id-171c4a25ccd-12",
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
                                "id": "id-171c4a25ccd-7",
                                "type": "text",
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
    ],
    "$LINES": {
      "id": "id-1709ef632af-2",
      "type": "lines",
      "value": "6"
    }
  }
}
}
