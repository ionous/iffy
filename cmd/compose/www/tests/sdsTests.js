
// sdsTests.js
// self describing storage
(function() {
  // sayTest storydata.
  const compactData= [
    "::test@id-3",{
      "test_name::text@id-0": "hello, goodbye",
      "go::execute": [
        "id-1::choose@id-7",{
          "false::execute": [
            "id-4::say@id-15",{
              "text::text_eval@id-14": [
              "::text_value@id-20", {
                "text::lines@id-19": "goodbye"
              }]}],
          "if::bool_eval@id-5": [
            "::bool_value@id-9",{
              "bool@id-8": "$TRUE"
            }],
          "true::execute": [
            "id-6::say@id-11",{
              "text::text_eval@id-10": [
                "::text_value@id-13",{
                  "text::lines@id-12": "hello"
                }]}]}],
      "lines@id-2": "hello"
  }];
  // compact
  // -
  const compacted= compact(getStory());
  const got= JSON.stringify(compacted, 0, 2);
  const want= JSON.stringify(compactData, 0, 2);
  // pull some pairs of closing parens together, and opening brackets onto the same line
  const out= got.replace(/([}\]])\n\s+([}\]])/g, "$1$2").replace(/",\n\s+([{\[])/g, '", $1');
  if (got !== want) {
    console.log(out);
    throw new Error("mismatch");
  }

  // expand
  // -

}());
