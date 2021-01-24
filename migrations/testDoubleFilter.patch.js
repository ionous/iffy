[
  {
    "op": "test",
    "path": "$..[?(@.value.type=='from_object')]",
    "reason": "replace from_object with from_text",
    "subpatches": [{
      "op": "replace",
      "path": "$.type",
      "value": "from_text"
    },{
      "op": "replace",
      "path": "$.value['$VAL'].type",
      "value": "text_eval"
    }]
  }
]
