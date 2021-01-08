[
  {
    "op": "test",
    "path": {
      "parent": "$..[?(@.type=='list_each')]"
    },
    "reason": "select the nodes containing type #choose.",
    "subpatches": [
      {
        "op": "remove",
        "path": {
          "parent": "$.value",
          "part": "$ELSE"
        },
        "reason": "remove #else block, always empty right now"
      },
      {
        "op": "move",
        "from": {
          "parent": "$.value",
          "part": "$GO"
        },
        "path": {
          "parent": "$.value",
          "part": "$DO"
        },
        "reason": "#go now #do"
      },
      {
        "op": "replace",
        "path": {
          "parent": "$.value",
          "part": "$AS"
        },
        "reason": "add #as seletion",
        "value": {
          "type": "list_iterator",
          "value": {
            "type": "as_???",
            "value": {
              "$VAR": {
                "type": "variable_name",
                "value": "???"
              }
            }
          }
        }
      },
      {
        "op": "move",
        "from": {
          "parent": "$.value['$WITH']",
          "part": "value"
        },
        "path": {
          "parent": "$.value['$AS'].value.value['$VAR']",
          "part": "value"
        },
        "reason": "move var name in"
      },
      {
        "op": "remove",
        "path": {
          "parent": "$.value",
          "part": "$WITH"
        },
        "reason": "remove #with"
      }
    ]
  }
]
