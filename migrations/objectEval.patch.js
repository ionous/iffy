[
  {
    "op": "test",
    "path": "$..[?(@.type=='get_field')]",
    "reason": "replace get_field with get_at_field",
    "subpatches": [
      {
        "op": "add",
        "path": "$.value['$FROM']",
        "value": {
          "type": "from_source_fields",
          "value": {
            "type": "from_var",
            "value": {
              "$VAR": {
                "type": "variable_name",
                "value": null
              }
            }
          }
        }
      },
      {
        "op": "copy",
        "from": "$.value['$OBJECT'].value.value['$NAME'].value",
        "path": "$.value['$FROM'].value.value['$VAR'].value",
        "reason": "notes 'patches' applies to the current doc, not the elements selected by the test."
      },
      {
        "op": "remove",
        "path": "$.value['$OBJECT']",
        "reason": "notes 'patches' applies to the current doc, not the elements selected by the test."
      },
      {
        "op": "replace",
        "path": "$.type",
        "value": "get_at_field"
      }
    ]
  },
  {
    "op": "test",
    "path": "$..[?(@.value.type=='object_name')]",
    "reason": "chop out object_name and replace it with its own text_eval",
    "subpatches": [
      {
        "op": "move",
        "from": "$.value.value['$NAME'].value",
        "path": "$.value"
      },
      {
        "op": "replace",
        "path": "$.type",
        "value": "text_eval"
      }
    ]
  },
  {
    "op": "test",
    "path": "$..[?(@.type=='from_object')]",
    "reason": "replace from_object with from_text",
    "subpatches": [
      {
        "op": "replace",
        "path": "$.type",
        "value": "from_text"
      }
    ]
  },
  {
    "op": "test",
    "path": "$..[?(@.type=='object_eval')]",
    "reason": "replace remaining object evals with text eval",
    "subpatches": [
      {
        "op": "replace",
        "path": "$.type",
        "value": "text_eval"
      }
    ]
  },
  {
    "op": "test",
    "path": "$..[?(@.type=='unpack')]",
    "reason": "replace unpacks containing unpacks with get_at_field.from_rec",
    "subpatches": [
      {
        "op": "test",
        "path": "$..[?(@.type=='unpack')]",
        "patches": [
          {
            "op": "add",
            "path": "$.value['$FROM']",
            "reason": "need the new source fields sub structure.",
            "value": {
              "type": "from_source_fields",
              "value": {
                "type": "from_rec",
                "value": {
                  "$REC": null
                }
              }
            }
          },
          {
            "op": "copy",
            "from": "$.value['$RECORD']",
            "path": "$.value['$FROM'].value.value['$REC']",
            "reason": "notes 'patches' applies to the current doc, not the elements selected by the test."
          },
          {
            "op": "remove",
            "path": "$.value['$RECORD']",
            "reason": "removes the old record evaluation."
          },
          {
            "op": "replace",
            "path": "$.type",
            "reason": "finally, rename #unpack to #get_at_field",
            "value": "get_at_field"
          }
        ]
      }
    ]
  },
  {
    "op": "test",
    "path": "$..[?(@.type=='unpack')]",
    "reason": "replace regular unpack nodes with get_at_field.from_var",
    "subpatches": [
      {
        "op": "add",
        "path": "$.value['$FROM']",
        "reason": "need the new source fields sub structure.",
        "value": {
          "type": "from_source_fields",
          "value": {
            "type": "from_var",
            "value": {
              "$VAR": {
                "type": "variable_name",
                "value": ""
              }
            }
          }
        }
      },
      {
        "op": "copy",
        "from": "$.value['$RECORD'].value.value['$NAME'].value",
        "path": "$.value['$FROM'].value.value['$VAR'].value",
        "reason": "notes 'patches' applies to the current doc, not the elements selected by the test."
      },
      {
        "op": "remove",
        "path": "$.value['$RECORD']",
        "reason": "notes 'patches' applies to the current doc, not the elements selected by the test."
      },
      {
        "op": "replace",
        "path": "$.type",
        "reason": "finally, rename #unpack to #get_at_field",
        "value": "get_at_field"
      }
    ]
  }
]
