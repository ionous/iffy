
// --------------------------------------------------------------------
// testing functions
function tagTest(str, match) {
  console.log(`testing ${str}`);
  const out= TagParser.parse(str);
  const mts= JSON.stringify(match.keys);
  const ots= JSON.stringify(out.keys);
  if (mts !== ots) {
    throw new Error("mismatch");
  }
  // test that our *match* specification make sense.
  const args= match.keys.filter(t=> t.startsWith("$")).sort();
  const sargs= JSON.stringify(args);
  const keys= Object.keys( match.args ).sort();
  const skeys= JSON.stringify(keys);
  if (sargs !== skeys) {
    throw new Error("mismatched spec");
  }
  // test each arg
  for (const arg of args) {
    const mp= match.args[arg];
    const mps= JSON.stringify(mp, Object.keys(mp).sort());
    //
    const op= match.args[arg];
    const ops= JSON.stringify(op, Object.keys(op).sort());
    if (mps !== ops) {
      throw new Error("mismatch");
    }
  }
}


function testSomeStuff() {
  tagTest("notags",{
    keys:["notags"],
    args:{}
  });
  tagTest("front {arg} back",{
    keys:["front ", "$ARG", " back"],
    args:{
      "$ARG": {
        label: "arg",
        arg: "arg",
        type: "arg",
      }
    }
  });
  tagTest("{arg} back",{
    keys:["$ARG", " back"],
    args:{
      "$ARG": {
        label: "arg",
        arg: "arg",
        type: "arg",
      }
    }
  });
  tagTest("front {arg}",{
    keys:["front ", "$ARG"],
    args:{
      "$ARG": {
        label: "arg",
        arg: "arg",
        type: "arg",
      }
    }
  });
  tagTest("{label%arg}",{
    keys:["$ARG"],
    args:{
      "$ARG": {
        label: "label",
        arg: "arg",
        type: "arg"
      }
    }
  });
  tagTest("{arg}",{
    keys:["$ARG"],
    args:{
      "$ARG": {
        label: "arg",
        arg: "arg",
        type: "arg",
      }
    }
  });
  tagTest("{arg|filter}",{
    keys:["$ARG"],
    args:{
      "$ARG": {
        label: "arg",
        arg: "arg",
        type: "arg",
        filters: ["filter"],
      }
    }
  });
  tagTest("{arg|filter=val}",{
    keys:["$ARG"],
    args:{
      "$ARG": {
        label: "arg",
        arg: "arg",
        type: "arg",
        filters: ["filter"],
        filterVals: {
          "filter":"val",
        },
      }
    }
  });
  tagTest("{arg|fe|fi|folter}",{
    keys:["$ARG"],
    args:{
      "$ARG": {
        label: "arg",
        arg: "arg",
        type: "arg",
        filters: ["fe", "fi", "folter"],
      }
    }
  });
  tagTest("{label%arg|filter}",{
    keys:["$ARG"],
    args:{
      "$ARG": {
        label: "label",
        arg: "arg",
        type: "arg",
        filters: ["filter"],
      }
    }
  });
  tagTest(String.raw`The {[description] is::%multiline}`,{
    keys:["The ","$MULTILINE"],
    args:{
      "$MULTILINE": {
        label: "description",
        // fix? change prefix/suffix into just "extra", always used as a suffix?
        // ( under the theory that the prefix can exist outside the tag )
        suffix: " is:",
        arg: "arg",
        type: "arg",
      }
    }
  });

  const reps={
    ":":{},
    "?":{optional: true},
    "*":{optional: true, repeats: true},
    "+":{repeats: true}
  };
  for ( const rep in reps )  {
    const ex= reps[rep];
    // {label%arg#type}
    tagTest("{label%arg"+rep+"type}", {
      keys:["$ARG"],
      args:{
        "$ARG": Object.assign({
          label: "label",
          arg: "arg",
          type: "type"
        }, ex)
      }
    });
    // {arg#type}
    tagTest("{arg"+rep+"type}",{
      keys:["$ARG"],
      args:{
        "$ARG": Object.assign({
          label: "arg",
          arg: "arg",
          type: "type"
        }, ex)
      }
    });
    // {arg}
    tagTest("{"+rep+"arg}",{
      keys:["$ARG"],
      args:{
        "$ARG": Object.assign({
          label: "arg",
          arg: "arg",
          type: "arg"
        }, ex)
      }
    });
  };
};
testSomeStuff();
