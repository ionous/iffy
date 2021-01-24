[
  {
    "op": "test",
    "path": {
      "parent": "$..[?(@.type=='choose')]"
    },
    "reason": "select the nodes containing type #choose.",
    "subpatches": [
      {
        "op": "test",
        "patches": [
          {
            "op": "replace",
            "path": "$.value['$ELSE']",
            "reason": "#else needs a different structure than #false due to brancher.",
            "value": {
              "type": "brancher",
              "value": {
                "type": "choose_nothing_else",
                "value": {
                  "$DO": {}
                }
              }
            }
          },
          {
            "from": {
              "parent": "$.value",
              "field": "$FALSE"
            },
            "op": "move",
            "path": {
              "parent": "$.value['$ELSE'].value.value",
              "field": "$DO"
            },
            "reason": "notes 'patches' applies to the current doc, not the elements selected by the test."
          }
        ],
        "path": {
          "parent": "$.value['$FALSE'].value['$EXE'].*"
        },
        "reason": "change #false to #else *if* false isnt empty."
      },
      {
        "op": "remove",
        "path": {
          "parent": "$.value",
          "field": "$FALSE"
        },
        "reason": "remove #false block in case it wasnt moved in the test."
      },
      {
        "op": "replace",
        "path": {
          "parent": "$.value",
          "field": "$DO"
        },
        "reason": "#do is required; #true was optional, so first create a blank.",
        "value": {
          "type": "activity",
          "value": {
            "$EXE": []
          }
        }
      },
      {
        "from": {
          "parent": "$.value['$TRUE'].value",
          "field": "$EXE"
        },
        "op": "copy",
        "path": {
          "parent": "$.value['$DO'].value",
          "field": "$EXE"
        },
        "reason": "now... copy the #true actions (if they exist)."
      },
      {
        "op": "remove",
        "path": {
          "parent": "$.value",
          "field": "$TRUE"
        },
        "reason": "remove #true ( if it existed. )"
      },
      {
        "op": "replace",
        "path": {
          "parent": "$",
          "field": "type"
        },
        "reason": "finally, rename #choose to #choose_action",
        "value": "get_at_field"
      }
    ]
  }
]
