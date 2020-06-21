function compact(src) {
  const makeTag= function(field, obj) {
    const ant= obj.type ? `::${obj.type}` : "";
    const id= obj.id? `@${obj.id}`:"";
    return `${field}${ant}${id}`;
  };
  // return an array of one tag and one object
  const makeRun= function(src) {
    const tag= makeTag("", src);
    const val= makeObject(src.value);
    return [tag, val];
  };
  const makeArray= function(src) {
    const arr= [];
    for (var el of src) {
      // slot-like
      var val;
      if ('type' in el) {
        val= makeRun(el.value);
      } else {
        val= makeValue(el.value);
      }
      const tag= makeTag("", el);
      arr.push( tag, val );
    }
    return arr;
  };
  const makeObject= function(src) {
    var ret;
    // slot-like
    if ('type' in src) {
       ret= makeRun(src);
    } else {
      const out= {};
      for (const lhs in src) {
        if (lhs.startsWith("$")) {
          const n= lhs.replace("$", "").toLowerCase();
          const rhs= src[lhs];
          const at= makeTag(n, rhs);
          out[at]= makeValue(rhs.value || rhs);
        }
      }
      ret= out;
    }
    return ret;
  };
  const makeValue= function(src) {
    var ret;
    switch (typeof src) {
      case 'undefined':
       ret= null;
      break;
      case 'string':
      case 'number':
      case 'boolean':
        ret= src;
      break;
      default:
        if (!src) {
          ret= null;
        } else if (Array.isArray(src)) {
          ret= makeArray(src);
        } else {
          ret= makeObject(src);
        }
    }; // switch
    return ret;
  }
  return makeRun(src);
};

// sdsTests.js
// self describing storage
(function() {
  var referenceData= {
    "id": "id-1709ef632af-3",
    "type": "test",
    "value": {
      "$TEST_NAME": {
        "id": "id-1709ef632af-0",
        "type": "text",
        "value": "hello, goodbye"
      },
      "$GO": [{
          "id": "id-1709ef632af-1",
          "type": "execute",
          "value": {
            "id": "id-1709ef632af-7",
            "type": "choose",
            "value": {
              "$FALSE": [{
                  "id": "id-1709ef632af-4",
                  "type": "execute",
                  "value": {
                    "id": "id-1709ef632af-15",
                    "type": "say",
                    "value": {
                      "$TEXT": {
                        "id": "id-1709ef632af-14",
                        "type": "text_eval",
                          "id": "id-1709ef632af-20",
                          "type": "text_value",
                          "value": {
                            "$TEXT": {
                              "id": "id-1709ef632af-19",
                              "type": "lines",
                              "value": "goodbye"
              }}}}}}],
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
              }}}},
              "$TRUE": [{
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
      }}}}}}}]}}}],
      "$LINES": {
        "id": "id-1709ef632af-2",
        "type": "lines",
        "value": "hello"
  }}};
  const compactData= [
    "::test@id-1709ef632af-3",{
      "test_name::text@id-1709ef632af-0": "hello, goodbye",
      "go": [
        "::execute@id-1709ef632af-1",[
          "::choose@id-1709ef632af-7",{
            "false": [
              "::execute@id-1709ef632af-4",[
                "::say@id-1709ef632af-15",{
                  "text::text_value@id-1709ef632af-20": {
                    "text::lines@id-1709ef632af-19": "goodbye"
                  }}]],
            "if::bool_eval@id-1709ef632af-5": [
              "::bool_value@id-1709ef632af-9",{
                "bool::bool@id-1709ef632af-8": "$TRUE"
              }],
            "true": [
              "::execute@id-1709ef632af-6",[
                "::say@id-1709ef632af-11",{
                  "text::text_eval@id-1709ef632af-10": [
                    "::text_value@id-1709ef632af-13",{
                      "text::lines@id-1709ef632af-12": "hello"
                    }]}]]}]],
      "lines::lines@id-1709ef632af-2": "hello"
  }];
  // compact
  // -
  const compacted= compact(referenceData);
  const got= JSON.stringify(compacted, 0, 2);
  const want= JSON.stringify(compactData, 0, 2);
  if (got !== want) {
    console.log(got);
    throw new Error("mismatch");
  }

  // expand
  // -

}());
