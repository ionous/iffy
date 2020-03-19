function getStory() {
  return {
    "id": "id-1709ef632af-3",
    "type": "test",
    "value": {
      "$TEST_NAME": {
        "id": "id-1709ef632af-0",
        "type": "text",
        "value": "hello, goodbye"
      },
      "$GO": [
        {
          "id": "id-1709ef632af-1",
          "type": "execute",
          "value": {
            "id": "id-1709ef632af-7",
            "type": "choose",
            "value": {
              "$FALSE": [
                {
                  "id": "id-1709ef632af-4",
                  "type": "execute",
                  "value": {
                    "id": "id-1709ef632af-15",
                    "type": "say",
                    "value": {
                      "$TEXT": {
                        "id": "id-1709ef632af-14",
                        "type": "text_eval",
                        "value": {
                          "id": "id-1709ef632af-20",
                          "type": "text_value",
                          "value": {
                            "$TEXT": {
                              "id": "id-1709ef632af-19",
                              "type": "lines",
                              "value": "goodbye"
                            }
                          }
                        }
                      }
                    }
                  }
                }
              ],
              "$IF": {
                "id": "id-1709ef632af-5",
                "type": "bool_eval",
                "value": {
                  "id": "id-1709ef632af-9",
                  "type": "bool_value",
                  "value": {
                    "$BOOL": {
                      "id": "id-1709ef632af-8",
                      "type": "bool",
                      "value": "$TRUE"
                    }
                  }
                }
              },
              "$TRUE": [
                {
                  "id": "id-1709ef632af-6",
                  "type": "execute",
                  "value": {
                    "id": "id-1709ef632af-11",
                    "type": "say",
                    "value": {
                      "$TEXT": {
                        "id": "id-1709ef632af-10",
                        "type": "text_eval",
                        "value": {
                          "id": "id-1709ef632af-13",
                          "type": "text_value",
                          "value": {
                            "$TEXT": {
                              "id": "id-1709ef632af-12",
                              "type": "lines",
                              "value": "hello"
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
      ],
      "$LINES": {
        "id": "id-1709ef632af-2",
        "type": "lines",
        "value": "hello"
      }
    }
  }
}
