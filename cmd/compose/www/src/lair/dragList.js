let lastList= 0;

class DragList {
  constructor(items, inline=false) {
    this.name= `list-${++lastList}`;
    this.items= items;
    this.inline= !!inline;
  }
  dropFrom(drop, from) {
    throw new Error("not implemented");
  }
};

(function() {
  console.log("testing drag list");
  function test(og, src, dst, expect) {
    class TestList extends DragList {
      makeBlank() { return "_" }
    }
    const items= og.split('');
    const dl= new TestList(items);
    dl.move(src,dst);
    const res= items.join("");
    if (expect !== res) {
     console.log("Error, moving", og[src], "want:", expect, "have:", res);
    }
  };

  // a
  test("abc", 0,-1, "a_bc"); // trailing head
  test("abc", 0, 0, "abc");  // <no change>
  test("abc", 0, 1, "_abc"); // leading b
  test("abc", 0, 2, "bac");  // leading c
  test("abc", 0, 3, "bca");  // leading tail

  // b
  test("abc", 1,-1, "bac");  // trailing head
  test("abc", 1, 0, "ab_c"); // trailing a
  test("abc", 1, 1, "abc");  // <no change>
  test("abc", 1, 2, "a_bc"); // leading c
  test("abc", 1, 3, "acb");  // leading tail

  // c
  test("abc", 2,-1, "cab");  // trailing head
  test("abc", 2, 0, "acb");  // trailing a
  test("abc", 2, 1, "abc_"); // trailing b
  test("abc", 2, 2, "abc");  // <no change>
  test("abc", 2, 3, "ab_c"); // leading tail

})/*()*/;
